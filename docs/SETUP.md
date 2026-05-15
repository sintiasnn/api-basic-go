# Environment Setup Guide

This document summarizes the setup options to meet the project prerequisites and links to step-by-step guides.

## Setup Options

- Linux (Ubuntu): follow `docs/LINUX_SYSTEMD_SETUP.md`.
  - Suitable for Linux servers/VMs or desktop Linux.
- Windows + WSL (Ubuntu): follow `docs/WSL_SYSTEMD_SETUP.md`.
  - Suitable for Windows users running WSL with systemd.

Both target Ubuntu with systemd enabled, Docker Engine + Compose v2, Git, and Go 1.24+.

## Quick Verification

Run the following commands to validate your environment:

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## Next Steps

- Read `README.md` for instructions on running the app (Go/Docker/Compose).
- See `REQUIREMENTS.md` for minimum requirements details.
