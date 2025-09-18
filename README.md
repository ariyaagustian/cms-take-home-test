# CMS Fullstack (Go + Gin + Postgres + Redis + MinIO)

> Proyek CMS fullstack berbasis Go. Gunakan `make` untuk menjalankan server & dependensi secara efisien.

---

## 🚀 Quickstart

1. Salin file `.env.example` menjadi `.env`, lalu sesuaikan nilai variabel jika perlu.
2. Jalankan dependensi dan server lokal:
   ```bash
   make dev
   ```
   Ini akan:
   - Menjalankan service `db`, `redis`, `minio`
   - Menunggu 5 detik agar semua service siap
   - Menjalankan backend menggunakan `go run`

3. Jalankan migrasi database:
   ```bash
   make migrate-up
   ```

4. Isi data awal (seed):
   ```bash
   make seed
   ```

5. Coba akses beberapa endpoint:

| Endpoint                         | Deskripsi                                |
|----------------------------------|-------------------------------------------|
| `GET /healthz`                  | Health check                              |
| `POST /api/auth/login`         | Login menggunakan data hasil seed         |
| `POST /api/content-types`      | Create Content Type                        |
| `POST /api/entries/:slug`      | Create Entry untuk konten tertentu         |
| `POST /api/media`              | Upload media ke MinIO                      |
| `GET /api/media/:id/preview`   | Dapatkan signed URL sementara (preview)   |

---

## 🧪 Testing

- Unit & integration test:
  ```bash
  make test
  ```

- E2E test:
  ```bash
  make e2e
  ```

---

## 🛠️ Perintah Makefile

| Perintah                | Fungsi                                       |
|-------------------------|----------------------------------------------|
| `make dev`              | Jalankan semua service & Go app (dev mode)   |
| `make build-backend`    | Build image backend menggunakan Docker       |
| `make run-backend-image`| Jalankan backend dari image Docker           |
| `make up`               | Jalankan semua container                     |
| `make down`             | Hentikan semua container                     |
| `make migrate-up`       | Jalankan migrasi database                    |
| `make migrate-down`     | Rollback satu migrasi                        |
| `make seed`             | Seed data awal                               |
| `make fmt`              | Format seluruh kode di folder `server/`      |

---

## 📦 Dependencies via Docker Compose

- **PostgreSQL** – database utama
- **Redis** – caching (opsional)
- **MinIO** – penyimpanan file/media
- **Imagor** – image resizer server

---

## 🗂️ Struktur Proyek

```
.
├── admin-ui/           # Kode Frontend (Vite + Tailwindcss)
├── server/             # Kode backend (Gin)
│   └── cmd/api         # Entry point aplikasi
│   └── internal/       # Model, repo, service, handler
├── migrations/         # File migrasi Postgres
├── tests/e2e/          # E2E test sederhana
├── .env.example        # Contoh konfigurasi ENV
├── docker-compose.yml  # Definisi semua service
├── Makefile            # Task otomatis
```

---

## 🧠 Tips

- Jika ingin menjalankan backend di dalam container:
  ```bash
  make build-backend
  make run-backend-image
  ```

- Jika ingin menjalankan langsung tanpa container:
  ```bash
  make dev
  ```

- Gagal konek DB? Pastikan `.env` sesuai:
  ```
  DB_HOST=db
  DB_PORT=5432
  ...
  ```

---

## ✅ Catatan

- Token JWT & password hashing sudah diterapkan.
- Audit log disimpan untuk aksi `entry`.
- Endpoint public tersedia di:
  - `GET /api/public/:slug`
  - `GET /api/public/:slug/:id`

---