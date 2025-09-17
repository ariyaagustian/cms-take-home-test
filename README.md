# CMS Fullstack (Go + Gin + Postgres + MinIO)

> Skeleton proyek untuk memenuhi _take-home_ CMS. Jalankan `docker compose` untuk deps, lalu `go run` untuk server.

## Quickstart
1. Salin `.env.example` menjadi `.env`, sesuaikan variabel.
2. Jalankan deps: `make up` atau `make dev` (dev juga menjalankan server).
3. Jalankan migrasi: `make migrate-up`
4. Seed data demo: `make seed`
5. Coba endpoint:
   - Health: `GET http://localhost:8080/healthz`
   - Login (seed): `POST http://localhost:8080/api/auth/login`
   - CT CRUD: `POST http://localhost:8080/api/content-types`
   - Entry CRUD: `POST http://localhost:8080/api/entries`

## Testing
- Unit & integration: `make test`
- E2E (sederhana): `make e2e`

## Deploy
- Build image: `make build`
- Health check: `/healthz`

## Catatan
- Struktur repo dan requirement mengikuti dokumen tes. Lihat `/docs` untuk arsitektur dan API docs ringkas.
