# Setup WSL (systemd), Docker, Git, dan Go 1.24+

 Panduan langkah demi langkah menyiapkan lingkungan di WSL (Windows) untuk Ubuntu agar memenuhi prasyarat proyek: WSL dengan systemd aktif, Docker + Compose v2, Git, dan Go 1.24+.

Ringkas:
- Target: Ubuntu di WSL.
- WSL: butuh systemd aktif (Windows 11/10 terbaru dengan WSL update).
- Hasil akhir: `systemctl` berjalan, Docker Engine aktif, Git terpasang, Go 1.24+ tersedia.

## 0) Siapkan/Update WSL (Windows)

Jalankan di PowerShell (Administrator):

```
wsl --update
wsl --list --online
wsl --install -d Ubuntu
```

Jika sudah punya distro WSL, cukup `wsl --update` untuk memastikan versi terbaru yang mendukung systemd.

## 1) Aktifkan systemd di WSL

Jalankan di terminal distro WSL (Ubuntu):

```
sudo nano /etc/wsl.conf
```

Isi/ubah agar mengandung:

```
[boot]
systemd=true
```

Simpan, lalu dari PowerShell (Windows):

```
wsl --shutdown
```

Buka kembali distro WSL dan verifikasi:

```
systemctl is-system-running
# atau
systemctl status
```

Jika belum active, pastikan WSL sudah di-update dan langkah `wsl --shutdown` dilakukan.

## 2) Pasang Git (Ubuntu)

```
sudo apt update
sudo apt install -y git
```

Verifikasi:

```
git --version
```

## 3) Pasang Docker Engine + Compose v2 (di dalam WSL Ubuntu)

Rekomendasi resmi untuk Ubuntu:

```
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Tambah key dan repo Docker (Ubuntu)
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Aktifkan dan mulai Docker via systemd
sudo systemctl enable --now docker

# (Opsional) jalankan tanpa sudo
sudo usermod -aG docker $USER
newgrp docker
```

Verifikasi:

```
docker version
docker compose version
docker run --rm hello-world
```

Catatan:
- Alternatif: gunakan Docker Desktop for Windows dengan integrasi WSL. Namun untuk skenario systemd + service, memasang Docker Engine di dalam distro WSL sering lebih konsisten.

## 4) Pasang Go 1.24+

Opsi A (disarankan, menggunakan tarball resmi):

```
GO_VER=1.24.0  # ganti ke rilis 1.24.x terbaru
curl -LO https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz

echo 'export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH"' | tee -a ~/.profile ~/.bashrc >/dev/null
source ~/.profile || true

go version
```

Opsi B (apt) — versi bisa lebih lama dari 1.24 (tidak direkomendasikan bila butuh 1.24+):

```
sudo apt update
sudo apt install -y golang
```

## 5) Verifikasi keseluruhan

Pastikan semua komponen siap:

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## 6) Jalankan proyek ini

Di root repo:

```
# Jalankan langsung
go run .

# Atau via Make
make run

# Atau dengan Docker
docker build -t api-basic-go .
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go

# Dengan Docker Compose
docker compose up --build -d
```

## Troubleshooting ringkas

- systemd tidak aktif di WSL:
  - Pastikan `/etc/wsl.conf` berisi `[boot]\nsystemd=true`, lakukan `wsl --shutdown`, lalu buka ulang distro.
  - Jalankan `wsl --update` di PowerShell untuk memperbarui WSL.
- Docker gagal start:
  - Cek `systemctl status docker` dan `journalctl -u docker -e`.
  - Pastikan user masuk grup `docker` (`id`), jika belum: `sudo usermod -aG docker $USER` lalu `newgrp docker`.
- Go tidak terdeteksi:
  - Pastikan PATH menyertakan `/usr/local/go/bin` dan `$HOME/go/bin`. Muat ulang shell atau `source ~/.profile`.
