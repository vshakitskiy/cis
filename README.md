# CIS — Corporate Information Systems

Project management app: projects, tasks (kanban board), time tracking, comments, file attachments.

**Backend:** Go, Chi, pgx, JWT auth
**Frontend:** React, TypeScript, Vite, Tailwind, shadcn/ui
**Database:** PostgreSQL 16

## Setup

You need Docker, Go 1.25+, [Goose](https://github.com/pressly/goose) and [Bun](https://bun.sh).

```bash
# start postgres
docker compose up -d

# run migrations
cd backend
make migrate

# set env (or use direnv — see backend/.envrc)
export DATABASE_URL=postgres://admeanie:assword@localhost:5432/dev?sslmode=disable
export JWT_SECRET=dev-secret-change-in-production

# start backend
go run ./cmd/server

# start frontend (in another terminal)
cd frontend
bun install
bun run dev
```

App runs at http://localhost:3000, API at http://localhost:8080.

## Environment Variables

| Variable       | Required | Default | Description               |
| -------------- | -------- | ------- | ------------------------- |
| `DATABASE_URL` | yes      | —       | Postgres connection string |
| `JWT_SECRET`   | yes      | —       | Token signing key          |
| `SERVER_PORT`  | no       | 8080    | Backend port               |

## API

All routes under `/api`. Auth routes are public, everything else needs `Authorization: Bearer <token>`.

```
POST   /api/auth/register
POST   /api/auth/login

GET    /api/projects
POST   /api/projects
GET    /api/projects/:id
PUT    /api/projects/:id
DELETE /api/projects/:id

GET    /api/projects/:pid/tasks
POST   /api/projects/:pid/tasks
GET    /api/projects/:pid/tasks/:id
PUT    /api/projects/:pid/tasks/:id
DELETE /api/projects/:pid/tasks/:id

GET    /api/projects/:pid/tasks/:tid/time-entries
POST   /api/projects/:pid/tasks/:tid/time-entries
DELETE /api/projects/:pid/tasks/:tid/time-entries/:id

GET    /api/projects/:pid/tasks/:tid/comments
POST   /api/projects/:pid/tasks/:tid/comments
PUT    /api/projects/:pid/tasks/:tid/comments/:id
DELETE /api/projects/:pid/tasks/:tid/comments/:id

GET    /api/projects/:pid/tasks/:tid/attachments
POST   /api/projects/:pid/tasks/:tid/attachments
GET    /api/projects/:pid/tasks/:tid/attachments/:id/download
DELETE /api/projects/:pid/tasks/:tid/attachments/:id

GET    /api/projects/:pid/report
```

## Database

Six tables: `users`, `projects`, `tasks`, `time_entries`, `comments`, `attachments`. See `ER.jpg` for the full diagram. Migrations are in `backend/migrations/`.

Roles: `admin`, `manager`, `employee`.
Task statuses: `todo`, `in_progress`, `done`.
Priorities: `low`, `medium`, `high`, `critical`.

## Project Layout

```
backend/
  cmd/server/       — entry point
  internal/
    handler/        — HTTP handlers
    service/        — business logic
    repository/     — SQL queries
    model/          — structs
    middleware/      — JWT auth
    config/         — env config
    db/             — connection pool
  migrations/       — goose SQL files
  uploads/          — attachment storage

frontend/src/
  api/              — fetch wrapper
  components/       — UI (shadcn/ui, forms, layout, tasks, etc.)
  context/          — auth state
  pages/            — routes
  lib/              — utils, constants
  types/            — TS types
```

## Production

```bash
# backend
cd backend && go build -o server ./cmd/server

# frontend
cd frontend && bun run build
# serve dist/ with nginx/caddy, proxy /api to the Go server
```

Change `JWT_SECRET`, use real DB credentials, update CORS origins in `cmd/server/main.go`.
