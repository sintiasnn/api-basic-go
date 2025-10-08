# Setup Ubuntu (systemd), Docker, Git, dan Go 1.21+

Panduan langkah demi langkah untuk Ubuntu (Linux native) agar memenuhi prasyarat proyek: systemd aktif, Docker + Compose v2, Git, dan Go 1.21+.

Ringkas:
- Distro: Ubuntu (direkomendasikan).
- Hasil akhir: `systemctl` berjalan, Docker Engine aktif, Git terpasang, Go 1.21+ tersedia.

## 1) Verifikasi systemd

Pastikan systemd aktif:

```
systemctl is-system-running
# atau
systemctl status
```

Jika perintah tidak ada, Anda mungkin menggunakan distro tanpa systemd. Gunakan distro dengan systemd (Ubuntu/Debian, dll.).

## 2) Pasang Git

```
sudo apt update
sudo apt install -y git

git --version
```

## 3) Pasang Docker Engine + Compose v2

Ikuti repositori resmi Docker untuk Ubuntu:

```
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Tambah key & repo Docker (Ubuntu)
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Aktifkan Docker via systemd
sudo systemctl enable --now docker

# (Opsional) gunakan docker tanpa sudo
sudo usermod -aG docker $USER
newgrp docker

docker version
docker compose version
```

## 4) Pasang Go 1.21+

Opsi A (disarankan — tarball resmi):

```
GO_VER=1.21.13  # ganti ke rilis 1.21.x terbaru
curl -LO https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz

echo 'export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH"' | tee -a ~/.profile ~/.bashrc >/dev/null
source ~/.profile || true

go version
```

Opsi B (apt) — sering lebih lama dari 1.21 (tidak direkomendasikan bila butuh 1.21+):

```
sudo apt update
sudo apt install -y golang
```

## 5) Verifikasi keseluruhan

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## 6) Jalankan proyek

Di root repo:

```
# Jalankan langsung
go run .

# Via Make
make run

# Docker build/run
docker build -t api-basic-go .
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go

# Docker Compose
docker compose up --build -d
```

## Troubleshooting singkat

- Docker gagal start:
  - Cek `systemctl status docker` dan `journalctl -u docker -e`.
  - Pastikan user masuk grup `docker` (`id`), jika belum: `sudo usermod -aG docker $USER` lalu `newgrp docker`.
- Go tidak terdeteksi:
  - Pastikan PATH menyertakan `/usr/local/go/bin` dan `$HOME/go/bin`. Muat ulang shell atau `source ~/.profile`.
