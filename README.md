# api-basic-go [![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev) [![systemd](https://img.shields.io/badge/systemd-deploy-139e4b?logo=systemd)](https://systemd.io) [![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

A simple Go HTTP API with zero external dependencies (`net/http` only) — built as a live demo for the workshop: **"Self-host your App without Docker: Why Not?"**

This repo demonstrates how to build, run, and self-host a Go application **without Docker**, using **systemd** as the service manager on a Linux host.

> Workshop slides: [Self-host your App without Docker: Why Not?](https://sintiasnn.github.io/assets/slides/Self-host%20your%20App%20without%20Docker_%20Why%20Not_.pdf)

---

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
- [Deployment](#deployment)
- [Docker](#docker)
- [CORS Configuration](#cors-configuration)
- [Author](#author)
- [References](#references)

---

## Overview

**What this repo demonstrates:**

- Building a zero-dependency Go HTTP API from scratch
- Running the app directly with `go run` (no build step)
- Deploying as a **systemd service** on a Linux host — **no Docker required**
- Alternative deployment via Docker / Docker Compose
- CORS configuration for frontend integration

**The thesis:** You don't always need Docker. For a simple Go binary, a systemd service is lighter, simpler, and more direct.

---

## Project Structure

```
api-basic-go/
├── main.go                         # Entry point — all logic in one file
├── go.mod                          # Go module (no external deps)
├── Makefile                        # Build, run, docker helpers
├── Dockerfile                      # Multi-stage distroless image
├── compose.yaml                    # Docker Compose service definition
├── .dockerignore
│
├── deploy/
│   └── systemd/
│       ├── api-basic-go.service    # systemd unit file
│       └── api-basic-go.env.example
│
├── docs/
│   ├── SETUP.md                    # Setup guide overview
│   ├── LINUX_SYSTEMD_SETUP.md      # Ubuntu native setup
│   └── WSL_SYSTEMD_SETUP.md        # Windows + WSL setup
│
├── REQUIREMENTS.md                 # Minimum requirements
└── README.md
```

---

## API Endpoints

### System

| Method | Path | Description | Example |
|--------|------|-------------|---------|
| `GET` | `/` | Welcome message | `curl http://localhost:8080/` |
| `GET` | `/health` | Health check | `curl http://localhost:8080/health` |
| `GET` | `/hello?name=Nama` | Greet by name | `curl 'http://localhost:8080/hello?name=Andi'` |

### Todo (in-memory, non-persistent)

| Method | Path | Description | Example |
|--------|------|-------------|---------|
| `GET` | `/todos` | List all todos | `curl http://localhost:8080/todos` |
| `POST` | `/todos` | Create a todo | `curl -X POST http://localhost:8080/todos -H 'Content-Type: application/json' -d '{"title":"Belajar Go","done":false}'` |
| `GET` | `/todos/{id}` | Get todo by ID | `curl http://localhost:8080/todos/1` |
| `PATCH` | `/todos/{id}` | Update todo fields | `curl -X PATCH http://localhost:8080/todos/1 -H 'Content-Type: application/json' -d '{"done":true}'` |
| `DELETE` | `/todos/{id}` | Delete a todo | `curl -X DELETE http://localhost:8080/todos/1 -i` |

> Data is stored in-memory — it will be lost on server restart.

---

## Getting Started

### Prerequisites

See [REQUIREMENTS.md](REQUIREMENTS.md) for full details. Minimum:

- Go 1.24+
- Port `8080` available (configurable via `PORT` env)

Full setup guides:
- Linux native: [`docs/LINUX_SYSTEMD_SETUP.md`](docs/LINUX_SYSTEMD_SETUP.md)
- Windows + WSL: [`docs/WSL_SYSTEMD_SETUP.md`](docs/WSL_SYSTEMD_SETUP.md)
- Overview: [`docs/SETUP.md`](docs/SETUP.md)

### Run directly

```bash
go run .
# or
make run
```

Server starts at `http://localhost:8080`.

---

## Deployment

### systemd (no Docker)

The primary deployment method — run the Go binary as a systemd service.

**1) Build and install the binary:**

```bash
make build
sudo cp bin/api-basic-go /usr/local/bin/api-basic-go
```

**2) (Optional) Create a dedicated user:**

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin api-basic-go || true
```

**3) Prepare environment file (optional):**

```bash
sudo cp deploy/systemd/api-basic-go.env.example /etc/api-basic-go.env
sudoedit /etc/api-basic-go.env
```

**4) Install and enable the service:**

```bash
sudo cp deploy/systemd/api-basic-go.service /etc/systemd/system/api-basic-go.service
sudo systemctl daemon-reload
sudo systemctl enable --now api-basic-go
```

**5) Check status and logs:**

```bash
systemctl status api-basic-go
journalctl -u api-basic-go -f
```

> If you created a dedicated user, uncomment the `User=`/`Group=` lines in the unit file, then `sudo systemctl daemon-reload && sudo systemctl restart api-basic-go`.

---

## Docker

Alternative deployment via container.

### Build & run

```bash
docker build -t api-basic-go .
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go
```

### Docker Compose

```bash
docker compose up --build -d
# or: make compose-up
```

---

## CORS Configuration

Configure allowed origins via the `CORS_ALLOWED_ORIGINS` environment variable:

```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://example.com go run .
```

- Comma-separated list of origins
- Use `*` to allow all origins (default if unset)

---

## Author

**Ni Putu Sintia Wati**

- GitHub: [@sintiasnn](https://github.com/sintiasnn)

---

## References

- [Go net/http Documentation](https://pkg.go.dev/net/http)
- [systemd Documentation](https://systemd.io)
- [Docker Documentation](https://docs.docker.com)
