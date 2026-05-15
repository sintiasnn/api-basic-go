# Setup WSL (systemd), Docker, Git, and Go 1.24+

Step-by-step guide for setting up the environment in WSL (Windows) with Ubuntu to meet the project prerequisites: WSL with systemd enabled, Docker + Compose v2, Git, and Go 1.24+.

Summary:
- Target: Ubuntu on WSL.
- WSL: requires systemd support (Windows 11 / recent Windows 10 with WSL update).
- End state: `systemctl` running, Docker Engine active, Git installed, Go 1.24+ available.

## 0) Prepare/Update WSL (Windows)

Run in PowerShell (Administrator):

```
wsl --update
wsl --list --online
wsl --install -d Ubuntu
```

If you already have a WSL distro, just run `wsl --update` to ensure you have the latest version with systemd support.

## 1) Enable systemd in WSL

Run inside your WSL distro terminal (Ubuntu):

```
sudo nano /etc/wsl.conf
```

Add/modify to contain:

```
[boot]
systemd=true
```

Save, then from PowerShell (Windows):

```
wsl --shutdown
```

Reopen your WSL distro and verify:

```
systemctl is-system-running
# or
systemctl status
```

If it's not active, make sure WSL has been updated and the `wsl --shutdown` step was performed.

## 2) Install Git (Ubuntu)

```
sudo apt update
sudo apt install -y git
```

Verify:

```
git --version
```

## 3) Install Docker Engine + Compose v2 (inside WSL Ubuntu)

Official recommendation for Ubuntu:

```
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Add Docker key and repo (Ubuntu)
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Enable and start Docker via systemd
sudo systemctl enable --now docker

# (Optional) run without sudo
sudo usermod -aG docker $USER
newgrp docker
```

Verify:

```
docker version
docker compose version
docker run --rm hello-world
```

Note:
- Alternative: use Docker Desktop for Windows with WSL integration. However, for the systemd + service scenario, installing Docker Engine inside the WSL distro is often more consistent.

## 4) Install Go 1.24+

Option A (recommended — official tarball):

```
GO_VER=1.24.0  # replace with latest 1.24.x release
curl -LO https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz

echo 'export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH"' | tee -a ~/.profile ~/.bashrc >/dev/null
source ~/.profile || true

go version
```

Option B (apt) — version may lag behind 1.24 (not recommended if 1.24+ is required):

```
sudo apt update
sudo apt install -y golang
```

## 5) Full verification

Make sure all components are ready:

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## 6) Run this project

From the repo root:

```
# Run directly
go run .

# Or via Make
make run

# Or with Docker
docker build -t api-basic-go .
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go

# With Docker Compose
docker compose up --build -d
```

## Quick Troubleshooting

- systemd not active in WSL:
  - Make sure `/etc/wsl.conf` contains `[boot]\nsystemd=true`, run `wsl --shutdown`, then reopen the distro.
  - Run `wsl --update` in PowerShell to update WSL.
- Docker fails to start:
  - Check `systemctl status docker` and `journalctl -u docker -e`.
  - Make sure your user is in the `docker` group (`id`); if not: `sudo usermod -aG docker $USER` then `newgrp docker`.
- Go not detected:
  - Make sure PATH includes `/usr/local/go/bin` and `$HOME/go/bin`. Reload your shell or run `source ~/.profile`.
