# BÁO CÁO DỰ ÁN — NOTES APP

---

## 1. GIỚI THIỆU

Notes App là ứng dụng full-stack ghi chú cho phép tạo, xem, sửa, xoá, tìm kiếm ghi chú và lọc theo thẻ (tag). Dự án được xây dựng với mục tiêu thực hành quy trình DevOps bao gồm:

- Thiết lập CI/CD tự động với GitHub Actions
- Container hoá database với Docker
- Triển khai lên VPS production
- Thiết lập monitoring và cảnh báo với Prometheus + Grafana + Alertmanager

---

## 2. KIẾN TRÚC HỆ THỐNG

```
┌─────────────┐       ┌─────────────┐       ┌──────────────┐
│  Frontend   │ ───→  │   Backend   │ ───→  │  PostgreSQL  │
│  React 19   │ :3000 │   Go + Gin  │ :8080 │      17      │
│  Vite 8     │       │   pgx/v5    │       │  (Docker)    │
│  TS 6       │       │             │       │              │
└─────────────┘       └──────┬──────┘       └──────────────┘
                             │
                    ┌────────┴────────┐
                    │   /metrics      │
                    │  (Prometheus)   │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │         Monitoring          │
              │  (VPS — native systemd)     │
              │                             │
              │  ┌──────────┐  ┌─────────┐  │
              │  │Prometheus│  │Grafana  │  │
              │  └─────┬────┘  └─────────┘  │
              │        │                     │
              │  ┌─────┴────┐               │
              │  │Alertma-  │               │
              │  │nager     │               │
              │  └──────────┘               │
              └─────────────────────────────┘
```

### Luồng dữ liệu

1. Người dùng truy cập Frontend (port 3000)
2. Frontend gọi API tới Backend (port 8080) qua `fetch()`
3. Backend xử lý request, CRUD với PostgreSQL
4. Backend expose metrics tại `/metrics` cho Prometheus scrape
5. Prometheus đánh giá alert rules → gửi Alertmanager → thông báo qua Telegram
6. Grafana trực quan hoá metrics từ Prometheus

---

## 3. TECH STACK

| Layer | Công nghệ | Phiên bản |
|---|---|---|
| **Frontend** | React + TypeScript + Vite | React 19, TS 6, Vite 8 |
| **Backend** | Go + Gin + pgx | Go 1.26.4, Gin v1.12, pgx v5 |
| **Database** | PostgreSQL (Docker) | 17 Alpine |
| **Monitoring** | Prometheus + Grafana + Alertmanager | (native systemd) |
| **Node exporter** | node_exporter | (native systemd) |
| **CI/CD** | GitHub Actions | — |
| **Deploy** | VPS Ubuntu + systemd + SCP | — |
| **Domain** | giaquyen.click | — |

### Frontend dependencies chính

| Package | Công dụng |
|---|---|
| `react`, `react-dom` | UI framework |
| `vite` + `@vitejs/plugin-react` | Bundler (rolldown) + dev server |
| `typescript` ~6.0.2 | Type checking |
| `eslint` + `typescript-eslint` | Linting (flat config, không type-aware) |

### Backend dependencies chính

| Package | Công dụng |
|---|---|
| `github.com/gin-gonic/gin` | HTTP framework + routing |
| `github.com/jackc/pgx/v5` | PostgreSQL driver + connection pool |
| `github.com/prometheus/client_golang` | Custom Prometheus metrics |

---

## 4. TÍNH NĂNG CHÍNH

### 4.1. CRUD Notes

| Thao tác | Method | Endpoint | Mô tả |
|---|---|---|---|
| Danh sách | GET | `/api/notes?search=&tag=` | Lọc theo nội dung (ILIKE) + tag (ANY) |
| Chi tiết | GET | `/api/notes/:id` | Lấy note theo UUID |
| Tạo mới | POST | `/api/notes` | Body: `{title, content, tags}` — title bắt buộc |
| Cập nhật | PUT | `/api/notes/:id` | Body: `{title, content, tags}` |
| Xoá | DELETE | `/api/notes/:id` | Xoá note, trả 404 nếu không tồn tại |

### 4.2. Tìm kiếm và lọc

- **Tìm kiếm:** sử dụng `ILIKE` với ký tự đại diện `%` — tìm trong cả `title` và `content`
- **Lọc theo tag:** sử dụng `ANY(tags)` — tìm notes có chứa tag cụ thể
- Kết hợp cả search và tag bằng `WHERE` + `AND`

### 4.3. Prometheus Metrics

Endpoint `/metrics` expose 2 metric custom:

| Metric | Loại | Labels |
|---|---|---|
| `notes_api_requests_total` | Counter | `method`, `path`, `status` |
| `notes_api_request_duration_seconds` | Histogram | `method`, `path` |

### 4.4. Seed dữ liệu mẫu

File `server/seed.sql` chứa 20 ghi chú mẫu với đa dạng chủ đề và tags, có thời gian tạo cách nhau để dễ kiểm tra tính năng.

---

## 5. CƠ SỞ DỮ LIỆU

### 5.1. Thông tin kết nối

- **Host:** `localhost:5432`
- **User:** `notes`
- **Password:** `notespass`
- **Database:** `notesdb`
- **Connection string:** `postgres://notes:notespass@localhost:5432/notesdb?sslmode=disable`
- **Cấu hình pool:** MaxConns = 10, MinConns = 2

### 5.2. Schema

```sql
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 5.3. Đặc điểm

- **Primary key:** UUID v4 (`gen_random_uuid()`) — phù hợp distributed systems
- **Tags:** kiểu `TEXT[]` — dùng `ANY()` để truy vấn
- **Migration:** tự động chạy khi server Go khởi động (embedded SQL)
- **`updated_at`:** cập nhật trong code Go (`UPDATE ... SET updated_at = NOW()`), không dùng DB trigger
- **Persistent volume:** Docker volume `pgdata` giữ dữ liệu qua container restart

---

## 6. CI/CD PIPELINE

### 6.1. CI (Continuous Integration)

File: `.github/workflows/ci.yml`

**Kích hoạt:** Mọi push hoặc pull request vào nhánh `main`

```
Trigger: push / pull_request → main
         │
    ┌────┴────┐
    │ go test │  go test ./... -v
    └────┬────┘
         │
    ┌────┴────┐
    │ go build│  go build -o notes-server ./cmd/server
    └────┬────┘
         │
    ┌────┴────────┐
    │ npm ci      │
    │ npm run build│  Node 22, cache npm
    └─────────────┘
```

**Lưu ý:** CI dùng Go 1.23 dù go.mod yêu cầu Go 1.26.4 — cần cập nhật nếu bump version.

### 6.2. CD (Continuous Deployment)

File: `.github/workflows/cd.yml`

**Kích hoạt:** Chỉ push vào nhánh `main` (không chạy trên PR)

```
Trigger: push → main
         │
    ┌────┴────┐
    │ go test │
    └────┬────┘
         │
    ┌────┴──────┐
    │ go build  │  → notes-server binary
    └────┬──────┘
         │
    ┌────┴────────────┐
    │ npm ci + build  │  → frontend/dist/
    └────┬────────────┘
         │
    ┌────┴──────┐
    │ SCP to VPS│  → /opt/notes-app/ (port 24700, user root)
    └────┬──────┘
         │
    ┌────┴───────────┐
    │ systemctl      │  restart notes-api service
    │ restart        │
    └────┬───────────┘
         │
    ┌────┴───────────────┐
    │ Health check       │  curl https://notes-app.giaquyen.click/api/health
    └────┬───────────────┘
         │
    ┌────┴───────────┐
    │ Telegram notify │  ✅ thành công / ❌ thất bại
    └────────────────┘
```

**Chi tiết deploy:**

- **SSH:** port 24700, user root
- **Service:** `notes-api` (systemd)
- **Health check:** `https://notes-app.giaquyen.click/api/health` — mong đợi HTTP 200
- **Secrets:** lưu trong GitHub Secrets (SSH_HOST, SSH_PORT, SSH_USER, SSH_PASSWORD, TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID)

---

## 7. MONITORING

### 7.1. Kiến trúc monitoring (VPS production)

```
┌──────────────────────────────────────┐
│              VPS                      │
│                                      │
│  ┌──────────┐    ┌───────────────┐   │
│  │   App     │    │ node_exporter │   │
│  │ :8080     │    │ :9100         │   │
│  │ /metrics  │    │ CPU/Disk/Net  │   │
│  └─────┬─────┘    └──────┬────────┘   │
│        │                 │            │
│  ┌─────┴─────────────────┴──────┐    │
│  │        Prometheus             │    │
│  │  scrape targets + alerting    │    │
│  └──────────────┬────────────────┘    │
│                 │                     │
│          ┌──────┴──────┐             │
│          │ Alertmanager│             │
│          │  Telegram   │             │
│          └─────────────┘             │
│                                      │
│  ┌──────────────┐                    │
│  │   Grafana    │                    │
│  │  Dashboard   │                    │
│  └──────────────┘                    │
└──────────────────────────────────────┘
```

### 7.2. Prometheus

- **Cấu hình:** `/etc/prometheus/prometheus.yml`
- **Alert rules:** `/etc/prometheus/alert.rules.yml`
- **File config trong repo:** `monitoring/prometheus/`

### 7.3. Alertmanager

- **Cấu hình:** `/etc/alertmanager/alertmanager.yml`
- **File config trong repo:** `monitoring/alertmanager/`
- **Kênh thông báo:** Telegram

### 7.4. Alerts

| Alert | Mô tả | Metric |
|---|---|---|
| **APIDOwn** | API health check fail | HTTP status != 200 |
| **HighCPU** | CPU usage vượt ngưỡng | node_exporter CPU metric |
| **DiskSpaceLow** | Disk usage > 90% | node_exporter disk metric |

### 7.5. Grafana

- Dashboard trực quan hoá metrics từ Prometheus
- Chạy native service: `grafana-server.service`

---

## 8. CẤU TRÚC THƯ MỤC

```
notes-app/
│
├── docker-compose.yml          # PostgreSQL 17 container
├── AGENTS.md                   # Hướng dẫn cho OpenCode AI
├── REPORT.md                   # Báo cáo dự án (file này)
├── .gitignore
├── README.md
│
├── .github/
│   └── workflows/
│       ├── ci.yml              # CI pipeline
│       └── cd.yml              # CD pipeline
│
├── frontend/                   # React + Vite
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tsconfig.app.json
│   ├── tsconfig.node.json
│   ├── eslint.config.js
│   └── src/
│       ├── main.tsx            # Entrypoint
│       ├── App.tsx             # Root component
│       ├── types/note.ts       # TypeScript types
│       ├── services/api.ts     # API client (fetch)
│       └── components/         # 6 components với CSS Modules
│           ├── NoteCard/
│           ├── NoteForm/
│           ├── NoteList/
│           ├── SearchBar/
│           ├── ConfirmDialog/
│           └── LoadingSpinner/
│
├── server/                     # Go + Gin
│   ├── go.mod / go.sum
│   ├── seed.sql                # 20 notes mẫu
│   ├── cmd/server/main.go      # Entrypoint, graceful shutdown
│   ├── migrations/
│   │   └── 001_create_notes.sql
│   └── internal/
│       ├── models/note.go      # Note struct + input types
│       ├── database/database.go # pgx pool + CRUD queries
│       ├── handlers/
│       │   ├── health.go       # GET /api/health
│       │   ├── health_test.go   # Unit test health
│       │   ├── notes.go        # CRUD handlers
│       │   └── handlers_test.go # Unit test Create validation
│       └── router/router.go    # Gin engine + CORS + Prometheus metrics
│
└── monitoring/                 # Monitoring config (production)
    ├── prometheus/
    │   ├── prometheus.yml
    │   └── alert.rules.yml
    └── alertmanager/
        └── alertmanager.yml
```

---

## 9. HƯỚNG DẪN CHẠY LOCAL

### Yêu cầu

- Docker & Docker Compose
- Go 1.26+
- Node.js (npm)

### Các bước

```bash
# 1. Khởi động PostgreSQL
docker compose up -d

# 2. Seed dữ liệu mẫu (20 notes)
# Windows PowerShell:
Get-Content server/seed.sql | docker compose exec -T postgres psql -U notes -d notesdb
# Linux/macOS:
cat server/seed.sql | docker compose exec -T postgres psql -U notes -d notesdb

# 3. Chạy API server (port 8080)
cd server && go run ./cmd/server

# 4. Chạy Frontend (port 3000) — terminal khác
cd frontend && npm run dev
```

### Commands hữu ích

| Thư mục | Lệnh | Mô tả |
|---|---|---|
| `frontend/` | `npm run dev` | Dev server HMR :3000 |
| `frontend/` | `npm run build` | `tsc -b` + `vite build` |
| `frontend/` | `npm run lint` | ESLint (flat config) |
| `server/` | `go run ./cmd/server` | Chạy API server |
| `server/` | `go build ./...` | Build |
| `server/` | `go vet ./...` | Static analysis |
| `server/` | `go test ./... -v` | Chạy unit tests |

---

## 10. DEPLOY PRODUCTION

### Thông tin VPS

| Thông số | Giá trị |
|---|---|
| Host | giaquyen.click |
| SSH port | 24700 |
| User | root |
| App directory | `/opt/notes-app/` |
| Systemd service | `notes-api` |
| Frontend | served qua reverse proxy (Nginx) |
| API domain | `https://notes-app.giaquyen.click` |

### Systemd services

| Service | Mô tả |
|---|---|
| `notes-api` | Go server (port 8080) |
| `prometheus` | Prometheus metrics collector |
| `alertmanager` | Alertmanager (Telegram notify) |
| `grafana-server` | Grafana dashboard |
| `node_exporter` | Node exporter (system metrics) |

### GitHub Secrets

| Secret | Mô tả |
|---|---|
| `SSH_HOST` | VPS IP address |
| `SSH_PORT` | SSH port (24700) |
| `SSH_USER` | SSH username (root) |
| `SSH_PASSWORD` | SSH password |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token |
| `TELEGRAM_CHAT_ID` | Telegram chat ID nhận thông báo |

---

