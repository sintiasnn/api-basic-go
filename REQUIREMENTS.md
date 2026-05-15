# Minimum Requirements

Target environment: Ubuntu (Linux native) or Ubuntu on Windows with WSL (with systemd).

Required:
- OS: Ubuntu (native) or Windows + WSL (Ubuntu, systemd enabled).
- Go: 1.24+ (as specified in `go.mod`).
- Docker: Engine 20.10+ and Docker Compose v2.
- Git: 2.30+ (or equivalent stable version).
- systemd: `systemctl` access to run the service (see `deploy/systemd/`).

Others:
- Port: `8080` available (can be changed via `PORT` env).
- Optional: GNU Make (for `Makefile`), `curl` (to test endpoints).

Notes:
- This project has no external dependencies (uses `net/http`).
- Main environment variables: `PORT`, `CORS_ALLOWED_ORIGINS`.

Step-by-step setup guides:
- Linux native: `docs/LINUX_SYSTEMD_SETUP.md`
- Windows + WSL: `docs/WSL_SYSTEMD_SETUP.md`
- Overview & links: `docs/SETUP.md`
