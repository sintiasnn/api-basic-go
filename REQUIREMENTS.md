# Minimum Requirements

Lingkungan target: Ubuntu (Linux native) atau Ubuntu pada Windows dengan WSL (dengan systemd).

Wajib:
- OS: Ubuntu (native) atau Windows + WSL (Ubuntu, systemd aktif).
- Go: 1.21+ (sesuai `go.mod`).
- Docker: Engine 20.10+ dan Docker Compose v2.
- Git: 2.30+ (atau versi stabil setara).
- systemd: akses `systemctl` untuk menjalankan service (lihat `deploy/systemd/`).

Lainnya:
- Port: `8080` tersedia (dapat diubah via env `PORT`).
- Opsional: GNU Make (untuk `Makefile`), `curl` (uji endpoint).

Catatan:
- Proyek tidak memiliki dependency eksternal (menggunakan `net/http`).
- Variabel lingkungan utama: `PORT`, `CORS_ALLOWED_ORIGINS`.

Panduan setup langkah demi langkah:
- Linux native: `docs/LINUX_SYSTEMD_SETUP.md`
- Windows + WSL: `docs/WSL_SYSTEMD_SETUP.md`
- Ringkasan & tautan: `docs/SETUP.md`
