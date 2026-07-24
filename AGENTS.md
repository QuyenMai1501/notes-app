# Notes App — AGENTS.md

## Stack
- **frontend/** — React 19 + TypeScript 6 + Vite 8, CSS Modules, port 3000
- **server/** — Go 1.26 + Gin + pgx/v5, port 8080
- **DB** — PostgreSQL 17 in Docker (docker-compose.yml)

## Quick start
```bash
docker compose up -d                        # start PostgreSQL
cd server && go run ./cmd/server            # start API server
cd frontend && npm run dev                  # start Vite dev server
```

## Commands
| frontend/ | server/ |
|---|---|
| `npm run dev` — Vite HMR :3000 | `go run ./cmd/server` — Gin :8080 |
| `npm run build` — `tsc -b && vite build` | `go build ./...` |
| `npm run lint` — `eslint .` (flat config) | `go vet ./...` |

No test runner configured in either package.

## Server layout
```
server/
├── cmd/server/main.go            # entrypoint, graceful shutdown
├── internal/
│   ├── models/note.go            # Note struct, Create/Update input types
│   ├── database/database.go      # pgx pool, migrations, CRUD queries
│   ├── handlers/health.go        # GET /api/health
│   ├── handlers/notes.go         # CRUD handlers (List, Get, Create, Update, Delete)
│   └── router/router.go          # Gin engine setup, CORS (allow :3000), routes
└── migrations/001_create_notes.sql  # schema reference
```

## API
| Method | Path | Notes |
|---|---|---|
| GET | /api/health | `{"status":"ok"}` |
| GET | /api/notes | `?search=&tag=` filters |
| GET | /api/notes/:id | |
| POST | /api/notes | body: `{title, content, tags}` — title required |
| PUT | /api/notes/:id | body: `{title, content, tags}` |
| DELETE | /api/notes/:id | |

## Database
- `notes` table: id (UUID PK), title, content, tags (TEXT[]), created_at, updated_at
- Migration runs automatically on server start (embedded SQL in `database.RunMigrations`)
- Connection via `DATABASE_URL` env var; defaults to `postgres://notes:notespass@localhost:5432/notesdb?sslmode=disable`

## Frontend quirks
- `typescript-eslint` **without** type-aware rules (no `project` in parserOptions)
- `verbatimModuleSyntax: true` — use `import type` for type-only imports
- `erasableSyntaxOnly: true` — no enums, no namespaces, no parameter properties
- CSS Modules — each component has its own `.module.css` file
- State via `useState` + `useReducer` (no external state library)
- React Compiler not enabled
