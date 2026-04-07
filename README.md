# Connexions Portable Template

A template for creating self-contained mock API servers using [Connexions](https://github.com/mockzilla/connexions).

Place your OpenAPI specs in `openapi/` and static responses in `static/`, build a single binary, and run it anywhere - no external files needed.

## Quick start

1. Click [**Use this template**](https://github.com/mockzilla/connexions-portable-template/generate) to create your own repository
2. Add your OpenAPI specs to `openapi/` and/or static responses to `static/`
3. Push to main - binaries for Linux, macOS, and Windows are built automatically and published to **Releases**

## OpenAPI specs

Drop `.yml`, `.yaml`, or `.json` OpenAPI spec files into the `openapi/` directory.
Each spec becomes a separate service, accessible at `/{service-name}/`.

## Static responses

Create directories under `static/{service-name}/{method}/` with JSON response files:

```
static/
  hello-world/
    get/
      index.json          -> GET /
    post/
      index.json          -> POST /
```

See the [Connexions docs](https://github.com/mockzilla/connexions) for more options.

## Release

Every push to main/master builds binaries for Linux, macOS, and Windows (amd64 and arm64) and publishes them to the **Releases** page. Download the binary for your platform from the `latest` release and run it locally:

```bash
./your-repo-name
```

## Mockzilla workflow

The included GitHub Actions workflow (`.github/workflows/mockzilla.yml`) publishes your specs to [Mockzilla](https://mockzilla.org) automatically:

- **Push to main/master** — publishes the latest specs to your main simulation
- **Pull request with `Mockzilla` label** — deploys a preview simulation for the PR (torn down when the PR is closed)

The `Mockzilla` label is created automatically on first push via the setup workflow.

Your simulation will be available at:
- `https://api.mockzilla.org/gh/{org}/{repo}/` — main branch
- `https://api.mockzilla.org/gh/{org}/{repo}/pr-{n}/` — per pull request

### Action inputs

You can customize the action in `.github/workflows/mockzilla.yml`:

```yaml
- uses: mockzilla/actions/portable@v1
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    region: us-east-1        # optional — preferred AWS region, used as a hint on first deploy only
    memory: 128              # optional — memory in MB (e.g. 128, 256, 512)
    timeout: 10000           # optional — request timeout in milliseconds
    timeout-minutes: 5       # optional — max minutes to wait for simulation to become active (default: 5)
```

| Input | Required | Description |
|---|---|---|
| `token` | yes | `GITHUB_TOKEN` — used to verify repo identity |
| `region` | no | Preferred AWS region (e.g. `us-east-1`, `ap-southeast-1`). Used as a hint on first deploy — if at capacity, the nearest available region is used. Has no effect after the simulation is deployed. |
| `memory` | no | Memory in megabytes (e.g. `128`, `256`, `512`). |
| `timeout` | no | Request timeout for the simulation in milliseconds. |
| `timeout-minutes` | no | Max minutes the action polls for the simulation to become active. Defaults to `5`. |
