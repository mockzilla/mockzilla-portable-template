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
