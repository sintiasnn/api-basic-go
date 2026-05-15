# Setup Ubuntu (systemd), Docker, Git, and Go 1.24+

Step-by-step guide for Ubuntu (Linux native) to meet the project prerequisites: systemd enabled, Docker + Compose v2, Git, and Go 1.24+.

Summary:
- Distro: Ubuntu (recommended).
- End state: `systemctl` running, Docker Engine active, Git installed, Go 1.24+ available.

## 1) Verify systemd

Make sure systemd is running:

```
systemctl is-system-running
# or
systemctl status
```

If the command is not found, you may be using a distro without systemd. Use a distro with systemd (Ubuntu/Debian, etc.).

## 2) Install Git

```
sudo apt update
sudo apt install -y git

git --version
```

## 3) Install Docker Engine + Compose v2

Follow the official Docker repository for Ubuntu:

```
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Add Docker key & repo (Ubuntu)
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Enable Docker via systemd
sudo systemctl enable --now docker

# (Optional) run docker without sudo
sudo usermod -aG docker $USER
newgrp docker

docker version
docker compose version
```

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

Option B (apt) — often lags behind 1.24 (not recommended if 1.24+ is required):

```
sudo apt update
sudo apt install -y golang
```

## 5) Full verification

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## 6) Run the project

From the repo root:

```
# Run directly
go run .

# Via Make
make run

# Docker build/run
docker build -t api-basic-go .
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go

# Docker Compose
docker compose up --build -d
```

## Quick Troubleshooting

- Docker fails to start:
  - Check `systemctl status docker` and `journalctl -u docker -e`.
  - Make sure your user is in the `docker` group (`id`); if not: `sudo usermod -aG docker $USER` then `newgrp docker`.
- Go not detected:
  - Make sure PATH includes `/usr/local/go/bin` and `$HOME/go/bin`. Reload your shell or run `source ~/.profile`.
