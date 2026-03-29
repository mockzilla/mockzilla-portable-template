# Connexions Portable Template

A template for creating self-contained mock API servers using [Connexions](https://github.com/mockzilla/connexions).

Place your OpenAPI specs in the `openapi/` directory, build a single binary, and run it anywhere - no external files needed.

## Quick start

1. Click **Use this template** on GitHub to create your own repository
2. Replace `openapi/petstore.yml` with your own OpenAPI spec(s)
3. Build and run:

```sh
go build -o mock-server .
./mock-server
```

The server starts on port 2200 by default.

## Flags

```
--port PORT     Override server port (default: 2200)
--config PATH   Unified config YAML (app settings + per-service config)
--context PATH  Per-service context YAML for value replacements
```

## Adding more specs

Drop any `.yml`, `.yaml`, or `.json` OpenAPI spec files into the `openapi/` directory. Each spec becomes a separate service, accessible at `/{service-name}/`.

## Configuration

To embed a config file, add `app.yml` to the project root and update the embed directive in `main.go`:

```go
//go:embed openapi app.yml
var content embed.FS
```

See the [Connexions docs](https://github.com/mockzilla/connexions) for config options.
