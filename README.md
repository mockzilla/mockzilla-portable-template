# Mockzilla Portable Template

An open-source API mock server template built with [Mockzilla](https://mockzilla.org/).

Source engine: [github.com/mockzilla/mockzilla](https://github.com/mockzilla/mockzilla)
Website: [mockzilla.org](https://mockzilla.org/)

Drop your OpenAPI specs into `services/<name>/`, push to main, and get a hosted simulation URL and a portable `.mockz` binary you can run anywhere. No code, no configuration, no separate infrastructure.

## Quick start

1. Click [**Use this template**](https://github.com/mockzilla/mockzilla-portable-template/generate) to create your own repository.
2. Add your services under `services/<name>/` (see layout below).
3. Push to main. The included GitHub Actions:
   - Publish your specs to a hosted simulation at `https://<label>.api.mockz.io`, where the label is your repo name.
   - Pack a `.mockz` archive and attach it to the latest GitHub release for offline use.

This repo uses the [Mockzilla engine](https://github.com/mockzilla/mockzilla) to serve realistic responses from OpenAPI specs.

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

Run it locally with the Mockzilla CLI:

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
- **Run by hand**: takes the mocks down and frees your simulation slot

The `Mockzilla` label is created automatically on first push via the setup workflow.

To take the mocks down, run the workflow manually with **delete** ticked, under
Actions → Mockzilla → Run workflow, or:

```bash
gh workflow run mockzilla.yml -f delete=true
```

The action removes a repository's mocks only on a run with `delete: true`, and
push and pull request triggers cannot pass an input, which is what the
`workflow_dispatch` trigger is there for. Handy on the free plan, where one
repository occupies the single slot.

Your simulation will be available at:
- `https://{label}.api.mockz.io`: main branch
- `https://{label}-pr{n}.api.mockz.io`: per pull request
- `https://{label}-{branch}.api.mockz.io`: any other branch you add to the workflow's `push` trigger

The label is the repo name as a host name, up to 55 characters, with `-2`, `-3` if
it is taken. It never ends in `-pr` and digits, which is kept for pull requests. The
`host` input asks for another before the first deploy. A branch has its name in its
host, lowercased, with everything but letters and digits removed: `feature/new-api`
gives `{label}-featurenewapi`. Branches deploy on plans with PR environments.

### Action inputs

You can customize the action in `.github/workflows/mockzilla.yml`:

```yaml
- uses: mockzilla/actions@v1
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    region: us-east-1        # optional. Preferred AWS region, used as a hint on first deploy only.
    environment: '{"ENV":"production","DEBUG":"true"}'  # optional
    host: mockzilla.net      # optional. A domain, or a label on one (petstore.mockzilla.net).
    services-dir: services   # optional. Directory with per-service folders (default: 'services').
    timeout-minutes: 5       # optional. Max minutes to wait for simulation to become active (default: 5).
    delete: false            # optional. Remove this repository from Mockzilla (default: false).
```

| Input | Required | Description |
|---|---|---|
| `token` | yes | `GITHUB_TOKEN`, used to verify repo identity. |
| `region` | no | Preferred AWS region (e.g. `us-east-1`, `ap-southeast-1`). Used as a hint on first deploy. If at capacity, the nearest available region is used. Has no effect after the simulation is deployed. |
| `environment` | no | JSON object of environment variables to set in the simulation (e.g. `'{"ENV":"production"}'`). |
| `host` | no | The domain the simulation answers on: `mockz.io`, `mockz.net`, `mockz.org`, `mockzilla.org`, `mockzilla.de` or `mockzilla.net`. Put a label in front to ask for it: `petstore.mockz.io` answers at `https://petstore.api.mockz.io`. Fixed at the first deploy. Defaults to the org setting or `mockz.io`, with the repo name as the label. |
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
      - uses: mockzilla/actions@v1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          delete: true
```

Trigger it manually from the **Actions** tab when you're ready.

### Check your simulation URL

The label is picked on the server, so the host can't be worked out locally. The
action prints the URL, sets it as its `url` output and posts it on the pull
request. To read it from the current pull request's comment:

```bash
gh pr view --json comments -q '.comments[].body' | grep -o 'simulation live at [^ ]*' | tail -1 | cut -d' ' -f4
```

## Disclaimer

This project is not affiliated with or endorsed by any of the API providers whose specifications may be used as examples.
