# Mockzilla Portable Template

A template for declaring mock API services and publishing them to [Mockzilla](https://github.com/mockzilla/mockzilla).

Drop each service into `services/<name>/`, push, and get both a hosted simulation and a downloadable `.mockz` package you can run locally.

## Quick start

1. Click [**Use this template**](https://github.com/mockzilla/mockzilla-portable-template/generate) to create your own repository.
2. Add your services under `services/<name>/` (see layout below).
3. Push to main. The included GitHub Actions:
   - Publish your specs to a hosted simulation at `https://api.mockz.io/gh/<org>/<repo>/`.
   - Pack a `.mockz` archive and attach it to the latest GitHub release for offline use.

## Layout

Each service lives in its own folder under `services/`. The folder name **is** the service identity: what you name it is what gets mounted at.

```
services/
  petstore/
    openapi.yml          # OpenAPI spec (any *.{yml,yaml,json} name works)
    config.yml           # optional: latency, errors, mount, upstream, cache
    context.yml          # optional: replacement values for mock data
  hello-world/
    v1/
      get/index.json     # → GET /hello-world/v1
      post/index.json    # → POST /hello-world/v1
app.yml                  # optional: global settings (port, history, etc.)
```

A service folder is in **static mode** as soon as it contains any
`<path>/index.<ext>` file. No `static/` wrapper. The verb defaults to
`GET`; wrap in a `<method>/` dir (`get`, `post`, etc.) as the
immediate parent of `index.<ext>` to override. The path before that
can be empty, one segment, or many.

If a folder has both a spec file and static endpoints, the two are
merged: spec endpoints register first, static files override matching
`(path, method)` pairs and add any new ones, and the spec file itself
is served at `GET /<service>/<filename>` as a literal asset for docs.

Routes the example serves:

```
GET  /petstore/pets             ← from services/petstore/openapi.yml
POST /petstore/pets
GET  /petstore/pets/{petId}
GET  /hello-world/v1            ← from services/hello-world/v1/get/index.json
POST /hello-world/v1
```

## Per-service config (`services/<name>/config.yml`)

```yaml
latency: 100ms              # constant latency
# OR percentile latencies
latencies:
  p50: 50ms
  p95: 200ms
errors:                     # percentile error injection
  p5: 500                   # 5% of requests → 500
mount: pets/v2              # override URL prefix (default: <folder-name>)
upstream:                   # forward to a real backend
  url: https://petstore3.swagger.io/api/v3
  timeout: 10s
cache:
  requests: true            # cache GET responses
```

## Per-service context (`services/<name>/context.yml`)

Flat replacement values. Keys are used by the mock generator to fill matching fields in responses. **No service-name wrapper.**

```yaml
name: ["Fluffy", "Spot", "Rover"]
tag: ["cat", "dog", "bird"]
```

## Static responses

Drop response files under `services/<name>/<path>/index.<ext>`. The
verb defaults to `GET`; wrap in a `<method>/` dir (`get`, `post`,
`put`, `patch`, `delete`, `head`, `options`, `trace`) as the
immediate parent of `index.<ext>` to override.

```
services/api/
  index.json                    → GET /api
  users/index.json              → GET /api/users      (implicit GET)
  users/post/index.json         → POST /api/users     (explicit method)
  users/{id}/index.json         → GET /api/users/{id}
  users/{id}/delete/index.json  → DELETE /api/users/{id}
```

Extension determines content-type:

| Extension | Content-Type |
|---|---|
| `.json` | `application/json` |
| `.html` / `.htm` | `text/html` |
| `.xml` | `application/xml` |
| `.yaml` / `.yml` | `application/yaml` |
| `.txt` | `text/plain` |

If a service folder has only static endpoints (no spec), the OpenAPI
document is synthesized from the file tree at startup.

## Release

Every push to main/master packs your service tree into a `.mockz`
archive (a gzipped tarball with a `.mockzilla.json` manifest declaring
the contained services) and publishes it to the `latest` release.

Run it locally with the mockzilla CLI:

```bash
brew tap mockzilla/tap
brew install mockzilla

curl -L "https://github.com/<org>/<repo>/releases/latest/download/<repo>.mockz" -o mocks.mockz
mockzilla mocks.mockz
```

To inspect what's inside before running:

```bash
mockzilla info mocks.mockz
```

## Mockzilla workflow

The included GitHub Actions workflow (`.github/workflows/mockzilla.yml`) publishes your specs to [Mockzilla](https://mockzilla.org) automatically:

- **Push to main/master**: publishes the latest specs to your main simulation
- **Pull request with `Mockzilla` label**: deploys a preview simulation for the PR (torn down when the PR is closed)

The `Mockzilla` label is created automatically on first push via the setup workflow.

Your simulation will be available at:
- `https://api.mockzilla.org/gh/{org}/{repo}/`: main branch
- `https://api.mockzilla.org/gh/{org}/{repo}/pr-{n}/`: per pull request

### Action inputs

You can customize the action in `.github/workflows/mockzilla.yml`:

```yaml
- uses: mockzilla/actions/portable@v1
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    region: us-east-1        # optional. Preferred AWS region, used as a hint on first deploy only.
    memory-size: 256         # optional. Memory in MB (default: 128).
    timeout: 60              # optional. Request timeout in seconds.
    environment: '{"ENV":"production","DEBUG":"true"}'  # optional
    host: api.mockzilla.net  # optional. API host for the simulation URL.
    services-dir: services   # optional. Directory with per-service folders (default: 'services').
    timeout-minutes: 5       # optional. Max minutes to wait for simulation to become active (default: 5).
    delete: false            # optional. Remove this repository from Mockzilla (default: false).
```

| Input | Required | Description |
|---|---|---|
| `token` | yes | `GITHUB_TOKEN`, used to verify repo identity. |
| `region` | no | Preferred AWS region (e.g. `us-east-1`, `ap-southeast-1`). Used as a hint on first deploy. If at capacity, the nearest available region is used. Has no effect after the simulation is deployed. |
| `memory-size` | no | Memory in megabytes (e.g. `128`, `256`, `512`). Defaults to `128`. |
| `timeout` | no | Request timeout for the simulation in seconds (e.g. `30`, `60`). |
| `environment` | no | JSON object of environment variables to set in the simulation (e.g. `'{"ENV":"production"}'`). |
| `host` | no | API host for the simulation URL (`api.mockzilla.org`, `api.mockzilla.de`, or `api.mockzilla.net`). Defaults to org setting or `api.mockzilla.org`. |
| `services-dir` | no | Directory containing per-service folders. Defaults to `services`. |
| `timeout-minutes` | no | Max minutes the action polls for the simulation to become active. Defaults to `5`. |
| `delete` | no | Remove this repository from Mockzilla. When set to `true`, the action skips publishing and deletes all mock APIs for this repo. Useful on the free plan to free up your slot before connecting a different repository. Defaults to `false`. |

### Removing this repository from Mockzilla

On the free plan you can only have one repository connected to Mockzilla at a time. To switch to a different repo, run the action with `delete: true` on the old one first:

```yaml
name: mockzilla-remove

on:
  workflow_dispatch:

jobs:
  remove:
    runs-on: ubuntu-latest
    steps:
      - uses: mockzilla/actions/portable@v1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          delete: true
```

Trigger it manually from the **Actions** tab when you're ready.

### Check your simulation URL

After pushing, your URL is deterministic:

```bash
echo "https://api.mockz.io/gh/$(gh repo view --json nameWithOwner -q .nameWithOwner)/$(git branch --show-current)/"
```
