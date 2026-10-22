# Rendering

MIT-licensed starter for a new HTML rendering service. Fork it, then replace the module path `rendering` in `go.mod` and the imports before you publish your own service.

## License

[MIT](LICENSE). Copyright (c) 2026 Rendering template contributors.

## Local setup

- Linux or macOS, with Docker, Compose file format v3.9, and `build-essential`
- `cp .env.dev.example .env.dev` and `cp .env.prod.example .env.prod`
- `make init` builds the dev container. Remote debug is on port 2346.
- Page: http://localhost:8081/
- Health: http://localhost:8081/healthz

`DEBUG=true` serves HTTP on `HTTP_PORT`. Otherwise the server requires `SSL_CERT_PATH` and `SSL_KEY_PATH` and serves TLS on `HTTP_PORT`. Do not commit real `.env` files or certificates.

## Routes

- `GET /` and `GET /index.html` render the placeholder index.
- `GET /healthz` returns `{"status":"ok"}`.
- Unknown paths render the 404 template. A known path with the wrong method renders 405.
- Status templates: 400, 401, 403, 404, 405, 500, 502, 503.

## Add a page

1. Add a `{{define "name"}}` template under `templates/`.
2. Return that name and view data from `internal/app`.
3. Register the route in `internal/adapters/http`.

## Tests and CI

- `make test` runs `go test -race ./...`. It does not need Docker, `.env`, or TLS certificates.
- `make lint` runs golangci-lint v2.15.0, the same version as CI.
- GitHub Actions runs test (`-race`), lint, govulncheck, and `go build ./cmd/api` on pull requests and pushes to `main` or `master`. CI uses no secrets and does not start the server.
