# api-basic-go

API sederhana dengan Go (tanpa dependency eksternal), menggunakan `net/http`.

## Menjalankan

- Prasyarat: Go 1.21+
- Jalankan server:

```
go run .
```

- Port default: `8080` (bisa diubah via env `PORT`).

## Endpoint

- `GET /`: Welcome message.
  - Contoh: `curl http://localhost:8080/`

- `GET /health`: Cek status.
  - Contoh: `curl http://localhost:8080/health`

- `GET /hello?name=Nama`: Sapa dengan nama.
  - Contoh: `curl 'http://localhost:8080/hello?name=Andi'`

- Todo (in-memory, non-persistent):
  - `GET /todos`: List semua todo.
    - `curl http://localhost:8080/todos`
  - `POST /todos`: Buat todo baru.
    - `curl -X POST http://localhost:8080/todos -H 'Content-Type: application/json' -d '{"title":"Belajar Go","done":false}'`
  - `GET /todos/{id}`: Ambil todo by id.
    - `curl http://localhost:8080/todos/1`
  - `PATCH /todos/{id}`: Update sebagian field (`title`, `done`).
    - `curl -X PATCH http://localhost:8080/todos/1 -H 'Content-Type: application/json' -d '{"done":true}'`
  - `DELETE /todos/{id}`: Hapus todo by id.
    - `curl -X DELETE http://localhost:8080/todos/1 -i`

Catatan: Data disimpan in-memory, akan hilang saat server restart.

## CORS

- Env `CORS_ALLOWED_ORIGINS`: daftar origin yang diizinkan (pisah koma), atau `*` untuk semua origin.
- Contoh: `CORS_ALLOWED_ORIGINS=http://localhost:3000,https://example.com go run .`

## Docker

- Build image:

```
docker build -t api-basic-go .
```

- Run container (port 8080):

```
docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* api-basic-go
```

## Deploy

### systemd

1) Build dan install binary ke `/usr/local/bin` (di host):

```
make build
sudo cp bin/api-basic-go /usr/local/bin/api-basic-go
```

2) (Opsional) Buat user khusus untuk menjalankan service:

```
sudo useradd --system --no-create-home --shell /usr/sbin/nologin api-basic-go || true
```

3) Siapkan environment file (opsional):

```
sudo cp deploy/systemd/api-basic-go.env.example /etc/api-basic-go.env
sudoedit /etc/api-basic-go.env
```

4) Install unit file dan aktifkan service:

```
sudo cp deploy/systemd/api-basic-go.service /etc/systemd/system/api-basic-go.service
sudo systemctl daemon-reload
sudo systemctl enable --now api-basic-go
```

5) Cek status dan logs:

```
systemctl status api-basic-go
journalctl -u api-basic-go -f
```

Catatan: Jika Anda membuat user `api-basic-go`, buka file unit dan uncomment baris `User`/`Group` lalu reload + restart.

### Docker Compose

```
docker compose up --build -d
# or: make compose-up
```
