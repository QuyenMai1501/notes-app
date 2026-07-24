# Notes App

Full-stack ghi chú với React + Go + PostgreSQL. Cho phép tạo, sửa, xoá, tìm kiếm ghi chú và lọc theo tag.

## Tech stack

| Layer | Công nghệ |
|---|---|
| Frontend | React 19, TypeScript 6, Vite 8, CSS Modules |
| Backend | Go 1.26, Gin, pgx/v5 |
| Database | PostgreSQL 17 (Docker) |

## Yêu cầu

- Docker & Docker Compose
- Go 1.26+
- Node.js (npm)

## Quick start

```bash
# 1. Khởi động PostgreSQL
docker compose up -d

# 2. Seed dữ liệu mẫu (20 notes)
Get-Content server/seed.sql | docker compose exec -T postgres psql -U notes -d notesdb

# 3. Chạy API server (port 8080)
cd server && go run ./cmd/server

# 4. Chạy frontend (port 3000) — terminal khác
cd frontend && npm run dev
```

## Cấu trúc thư mục

```
notes-app/
├── docker-compose.yml     # PostgreSQL 17 container
├── AGENTS.md              # Hướng dẫn cho OpenCode AI
├── frontend/              # React + Vite (port 3000)
│   ├── src/
│   │   ├── components/    # 6 components với CSS Modules
│   │   ├── services/api.ts
│   │   ├── types/note.ts
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   └── vite.config.ts
└── server/                # Go + Gin (port 8080)
    ├── cmd/server/main.go
    ├── internal/
    │   ├── models/
    │   ├── database/
    │   ├── handlers/
    │   └── router/
    ├── migrations/
    └── seed.sql
```

## API Endpoints

| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/health` | Kiểm tra trạng thái |
| GET | `/api/notes` | Danh sách notes (`?search=&tag=`) |
| GET | `/api/notes/:id` | Chi tiết note |
| POST | `/api/notes` | Tạo note (`title`, `content`, `tags`) |
| PUT | `/api/notes/:id` | Cập nhật note |
| DELETE | `/api/notes/:id` | Xoá note |

## Commands

### Frontend

| Lệnh | Mô tả |
|---|---|
| `npm run dev` | Dev server HMR :3000 |
| `npm run build` | Typecheck + bundle |
| `npm run lint` | ESLint kiểm tra code |

### Server

| Lệnh | Mô tả |
|---|---|
| `go run ./cmd/server` | Chạy API server |
| `go build ./...` | Build |
| `go vet ./...` | Kiểm tra code |
| `go mod tidy` | Cập nhật dependencies |

## Database

- PostgreSQL 17 chạy trong Docker container
- Volume `pgdata` giữ dữ liệu persistent
- Kết nối mặc định: `postgres://notes:notespass@localhost:5432/notesdb`
- Migration tự động chạy khi server Go khởi động
- `updated_at` được cập nhật trong code Go (không phải DB trigger)

## Seed dữ liệu

```bash
Get-Content server/seed.sql | docker compose exec -T postgres psql -U notes -d notesdb
```
