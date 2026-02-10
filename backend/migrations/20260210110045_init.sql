-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('admin', 'manager', 'employee');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
        CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_priority') THEN
        CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'critical');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  role user_role NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS projects (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  owner_id BIGINT NOT NULL REFERENCES users (id),
  deadline DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tasks (
  id BIGSERIAL PRIMARY KEY,
  project_id BIGINT NOT NULL REFERENCES projects (id),
  assignee_id BIGINT REFERENCES users (id),
  title TEXT NOT NULL,
  description TEXT,
  status task_status NOT NULL DEFAULT 'todo',
  priority task_priority NOT NULL DEFAULT 'medium',
  deadline DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS time_entries (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES tasks (id),
  user_id BIGINT NOT NULL REFERENCES users (id),
  minutes INT NOT NULL CHECK (minutes > 0),
  description TEXT,
  date DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS comments (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES tasks (id),
  user_id BIGINT NOT NULL REFERENCES users (id),
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS attachments (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES tasks (id),
  uploaded_by BIGINT NOT NULL REFERENCES users (id),
  filename TEXT NOT NULL,
  filepath TEXT NOT NULL UNIQUE,
  size BIGINT NOT NULL CHECK (size >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_projects_owner_id ON projects (owner_id);

CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks (project_id);

CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks (assignee_id);

CREATE INDEX IF NOT EXISTS idx_time_entries_task_id ON time_entries (task_id);

CREATE INDEX IF NOT EXISTS idx_time_entries_user_date ON time_entries (user_id, date);

CREATE INDEX IF NOT EXISTS idx_comments_task_id ON comments (task_id);

CREATE INDEX IF NOT EXISTS idx_attachments_task_id ON attachments (task_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attachments;

DROP TABLE IF EXISTS comments;

DROP TABLE IF EXISTS time_entries;

DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS projects;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS task_priority;

DROP TYPE IF EXISTS task_status;

DROP TYPE IF EXISTS user_role;

-- +goose StatementEnd
