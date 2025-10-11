# Panduan Setup Lingkungan

Dokumen ini merangkum opsi setup untuk memenuhi prasyarat proyek dan menautkan ke panduan langkah demi langkah.

## Opsi Setup

- Linux (Ubuntu): ikuti `docs/LINUX_SYSTEMD_SETUP.md`.
  - Cocok untuk server/VM Linux atau desktop Linux.
- Windows + WSL (Ubuntu): ikuti `docs/WSL_SYSTEMD_SETUP.md`.
  - Cocok untuk pengguna Windows yang menggunakan WSL dengan systemd.

Keduanya menargetkan Ubuntu dengan systemd aktif, Docker Engine + Compose v2, Git, dan Go 1.24+.

## Verifikasi Cepat

Jalankan perintah berikut untuk memvalidasi lingkungan:

```
systemctl is-system-running
systemctl status docker

git --version
docker version
docker compose version
go version
```

## Langkah Selanjutnya

- Baca `README.md` untuk cara menjalankan aplikasi (Go/Docker/Compose).
- Lihat `REQUIREMENTS.md` untuk rincian prasyarat minimum.
