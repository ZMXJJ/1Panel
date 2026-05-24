# AGENTS.md

## Cursor Cloud specific instructions

### Architecture Overview

1Panel is a 3-component system:
- **Agent** (Go, `/workspace/agent/`): system operations executor, listens on Unix socket `/etc/1panel/agent.sock`
- **Core** (Go, `/workspace/core/`): API server on port 9999, proxies to Agent via Unix socket
- **Frontend** (Vue 3 + Vite, `/workspace/frontend/`): dev server on port 4004, proxies `/api/v2` → `localhost:9999`

### Prerequisites

- Go 1.25.10+ installed at `/usr/local/go/bin/go`
- Node.js v22+ (pre-installed)
- Required system directories: `/opt/1panel/{conf,db,log,secret}`, `/etc/1panel/`
- Config file at `/opt/1panel/conf/app.yaml` (copy from `core/cmd/server/conf/app.yaml`)
- Params file at `/usr/local/bin/1pctl` with `BASE_DIR=/opt`, `ORIGINAL_PORT=9999`, `ORIGINAL_VERSION=v2.0.0`

### Running Services (dev mode)

Start order matters — agent first, then core:

```bash
# 1. Build
export PATH=/usr/local/go/bin:$PATH
cd /workspace/agent && CGO_ENABLED=0 go build -trimpath -o /workspace/build/1panel-agent ./cmd/server/main.go
cd /workspace/core && CGO_ENABLED=0 go build -trimpath -o /workspace/build/1panel-core ./cmd/server/main.go

# 2. Start agent (creates /etc/1panel/agent.sock)
sudo /workspace/build/1panel-agent

# 3. Start core (listens on :9999)
sudo /workspace/build/1panel-core

# 4. Start frontend dev server
cd /workspace/frontend && npm run dev
```

### Important Gotchas

- The core binary embeds web assets via `//go:embed` in `core/cmd/server/web/web.go`. For dev builds, create placeholder files: `mkdir -p core/cmd/server/web/assets && touch core/cmd/server/web/assets/.gitkeep` and create a minimal `core/cmd/server/web/index.html`.
- Both services must run as root (sudo) because the agent creates restricted Unix socket with 0600 permissions.
- Default login: `admin` / `admin123` (configured in `/opt/1panel/conf/app.yaml`).
- The frontend type-check (`npm run type-check`) has pre-existing TS errors — this is a known state of the repo.
- The `Docker Compose command not found` error in agent logs is expected when Docker is not installed — it's non-fatal.

### Lint / Test / Build Commands

| Component | Lint | Test | Build |
|-----------|------|------|-------|
| Frontend | `cd frontend && npx eslint --ext .js,.ts,.vue ./src` | `npm run type-check` (has known errors) | `npm run build:pro` |
| Core | `cd core && go vet ./...` | `cd core && go test ./...` | `cd core && CGO_ENABLED=0 go build -trimpath -o ../build/1panel-core ./cmd/server/main.go` |
| Agent | `cd agent && go vet ./...` | `cd agent && go test ./...` | `cd agent && CGO_ENABLED=0 go build -trimpath -o ../build/1panel-agent ./cmd/server/main.go` |
