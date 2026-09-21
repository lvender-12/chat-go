# chat-go

Backend chat API berbasis Go (Fiber) dengan autentikasi JWT (cookie), fitur teman, dan realtime chat via WebSocket + RabbitMQ.

## Fitur

- Auth: register, login, logout (JWT di cookie `AuthToken`)
- Users: ambil/edit profil, upload avatar
- Friends: kirim / terima / tolak friend request, daftar teman
- Chat: riwayat chat via WebSocket, kirim/edit/hapus pesan (broadcast realtime via RabbitMQ)
- Response API seragam (`status`, `message`, `data`)
- Swagger UI di `/docs`
- Migrasi database otomatis (opsional via config)

## Tech stack

| Layer | Teknologi |
|-------|-----------|
| HTTP / WS | Fiber v3, Fiber WebSocket |
| Auth | JWT (`golang-jwt`) + cookie |
| Password | Argon2 |
| Database | MySQL |
| Message broker | RabbitMQ (`amqp091-go`) |
| Migrasi | golang-migrate |
| Docs | swag / Swagger UI |

## Struktur proyek

```text
chat-go/
├── .github/workflows/   # CI (GitHub Actions)
├── cmd/server/          # entrypoint aplikasi
├── config/              # config.json (lokal) & config.example.json
├── docs/                # swagger generated
├── internal/
│   ├── app/             # State, Response, ErrorHandler
│   ├── auth/            # register / login / logout
│   ├── users/           # profile, edit profile, upload avatar
│   ├── friends/         # friend requests & list
│   ├── chat/            # websocket hub + messaging
│   ├── rabbit/          # koneksi RabbitMQ + chat bus (publish/consume)
│   ├── middleware/      # auth & conversation access
│   ├── config/          # load config
│   ├── database/        # DB + migrate
│   ├── router/          # route setup
│   └── utils/           # jwt, password, validation
├── migrations/          # SQL migrations
├── public/              # static/public assets
└── test/                # unit tests
```

## Prasyarat

- Go 1.27+
- MySQL
- RabbitMQ

## Setup

1. Clone repo dan masuk ke folder proyek.

2. Salin config contoh:

```bash
cp config/config.example.json config/config.json
```

3. Sesuaikan `config/config.json` (DB, JWT secret, RabbitMQ, migration path, dll).

Contoh struktur config:

```json
{
  "app": {
    "name": "chat-app",
    "host": "0.0.0.0",
    "port": 3000,
    "log_level": "debug"
  },
  "database": {
    "driver": "mysql",
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "password",
    "name": "chatgo",
    "charset": "utf8mb4",
    "parse_time": true
  },
  "jwt": {
    "secret": "change-me",
    "exp": 24
  },
  "migration": {
    "enabled": true,
    "path": "./migrations"
  },
  "storage": {
    "path": "./public"
  },
  "cors": {
    "allow_origins": [
      "http://localhost:5500",
      "http://127.0.0.1:5500"
    ],
    "allow_methods": [
      "GET",
      "POST",
      "PUT",
      "PATCH",
      "DELETE",
      "OPTIONS"
    ],
    "allow_headers": [
      "Content-Type",
      "Accept"
    ],
    "allow_credentials": true
  },
  "rabbit": {
    "host": "127.0.0.1",
    "name": "guest",
    "password": "guest",
    "port": "5672"
  }
}
```

Tambah origin frontend baru cukup di `cors.allow_origins`.

> `config/config.json` di-ignore git. Jangan commit secret production.

4. Buat database MySQL sesuai `database.name`.

5. Pastikan RabbitMQ berjalan dan kredensial di `rabbit` valid.

6. Install dependency & jalankan server:

```bash
go mod tidy
go run ./cmd/server
```

Server default listen di `:3000`.

Swagger UI: [http://localhost:3000/docs](http://localhost:3000/docs)

## Realtime chat (Hub + RabbitMQ)

- **Hub** (`internal/chat`): track koneksi WebSocket lokal per conversation.
- **RabbitMQ bus** (`internal/rabbit`): fanout exchange `chat.fanout` untuk sebar event antar instance.
- Alur: `Send` / `Edit` / `Delete` → simpan DB → `bus.Publish` → consumer → `hub.Broadcast` → WebSocket client.

## Format response

Semua response HTTP (success & error) memakai struct yang sama:

```json
{
  "status": 200,
  "message": "success",
  "data": {}
}
```

- Ada payload → diisi di `data`
- Tidak ada payload → `data` bernilai `null`
- Error dari `fiber.NewError` / error lain diformat lewat `app.ErrorHandler`

## API overview

Base path: `/api/v1`

### Auth

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| POST | `/auth/register` | - | Register user |
| POST | `/auth/login` | - | Login, set cookie `AuthToken` |
| POST | `/auth/logout` | - | Clear cookie |

### Users

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/users/profile` | Cookie | Profil user login |
| PATCH | `/users/profile` | Cookie | Edit profil (field opsional: `username`, `email`, `display_name`) |
| POST | `/users/profile/avatar` | Cookie | Upload avatar (`multipart/form-data`, field `avatar`, max 5 MB) |

### Friends

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| POST | `/friend/add-friend` | Cookie | Kirim friend request (atau auto-accept jika ada request balik) |
| POST | `/friend/accept-friend/:id` | Cookie | Terima request (by request ID); buat conversation |
| POST | `/friend/reject-friend/:id` | Cookie | Tolak request |
| GET | `/friend/friends` | Cookie | Daftar teman |
| GET | `/friend/friends-request` | Cookie | Daftar friend request |

### Chat

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/chat/:id` | Cookie + WS | WebSocket + riwayat chat |
| POST | `/chat/:id` | Cookie | Kirim pesan |
| PUT | `/chat/:id` | Cookie | Edit pesan |
| DELETE | `/chat/:id` | Cookie | Soft delete pesan |

`:id` = conversation ID. User harus anggota conversation tersebut.

## Testing

```bash
go test ./... -count=1 -race
```

Cakupan unit test saat ini:

- `app.JSON` / `ErrorHandler`
- email validation
- password hash & verify
- JWT generate/parse
- ambil user ID dari cookie token
- load config
- middleware `CheckAuth`

## CI

GitHub Actions (`.github/workflows/ci.yml`) jalan otomatis di setiap push/PR ke `main`:

1. Setup Go (versi dari `go.mod`)
2. `go mod download` + `go mod verify`
3. `go vet ./...`
4. `go test ./... -race`
5. `go build ./cmd/server`

## Development notes

- Regenerasi Swagger (jika annotation berubah):

```bash
swag init -g cmd/server/main.go -o docs
```

- Log level mengikuti `app.log_level` (`debug` / lainnya).
- Profile upload path publik di-ignore: `/public/profile`.
