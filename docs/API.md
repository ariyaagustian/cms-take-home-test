# 📖 CMS API Documentation

Berikut adalah daftar endpoint API dari CMS Fullstack.

---

## 🩺 Health Check
### `GET /healthz`
- **Deskripsi**: Mengecek apakah server berjalan.
- **Response**: 
  ```json
  { "status": "ok" }
  ```

---

## 🔐 Authentication
### `POST /api/auth/login`
- **Deskripsi**: Login user.
- **Body**:
  ```json
  {
    "email": "admin@cms.local",
    "password": "yourpassword"
  }
  ```
- **Response**:
  ```json
  {
    "token": "jwt-token",
    "token_type": "Bearer",
    "user": {
      "id": "uuid",
      "email": "admin@cms.local",
      "name": "Admin",
      "role": "Admin"
    }
  }
  ```

### `POST /api/auth/register`
- **Deskripsi**: Registrasi user baru.
- **Body**:
  ```json
  {
    "email": "user@cms.local",
    "username": "User",
    "password": "password123"
  }
  ```

---

## 🗂️ Content Types
### `POST /api/content-types`
- Buat content type baru.

### `GET /api/content-types`
- Ambil semua content types.

### `GET /api/content-types/:id`
- Detail satu content type.

### `PUT /api/content-types/:id`
- Update content type.

### `DELETE /api/content-types/:id`
- Hapus content type.

### `POST /api/content-types/:id/fields`
- Tambah field ke content type.
- **Body**:
  ```json
  {
    "name": "content",
    "kind": "wysiwyg",
    "options": {}
  }
  ```

---

## ✍️ Entries
### `POST /api/entries/:slug`
- Buat entry baru untuk content type `:slug`.

### `GET /api/entries/:slug`
- List semua entry by content type.

### `GET /api/entries/:slug/:id`
- Detail entry.

### `PUT /api/entries/:slug/:id`
- Update entry.

### `DELETE /api/entries/:slug/:id`
- Hapus entry.

### `POST /api/entries/:slug/:id/publish`
- Publish entry langsung.

### `POST /api/entries/:slug/:id/rollback/:version`
- Rollback entry ke versi tertentu.

---

## 🖼️ Media
### `POST /api/media`
- Upload file ke MinIO.

### `GET /api/media`
- List semua media.

### `GET /api/media/preview/:id`
- Ambil signed URL untuk preview.

### `DELETE /api/media/:id`
- Hapus media.

---

## 👤 Admin Roles & Users
### Roles
- `GET /api/admin/roles` → List roles
- `POST /api/admin/roles` → Create role

### Users
- `GET /api/admin/users` → List users
- `GET /api/admin/users/:id/roles` → Ambil role user
- `POST /api/admin/users/:id/roles` → Set role untuk user

---

## 🌐 Public API
### `GET /api/public/:slug`
- List published entries berdasarkan content type `:slug`.

### `GET /api/public/:slug/:id`
- Detail entry published berdasarkan `id`.

---

## 📝 Catatan
- Semua endpoint **private** (auth/admin) butuh **JWT Bearer Token** di header:
  ```
  Authorization: Bearer <token>
  ```
- Public API (`/api/public/...`) tidak butuh authentication.
