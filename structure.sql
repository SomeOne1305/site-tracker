-- Users table
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100)  NOT NULL,
    last_name  VARCHAR(100)  NOT NULL,
    email      VARCHAR(255)  UNIQUE NOT NULL,
    is_verified BOOLEAN      DEFAULT FALSE,
    password   VARCHAR(255)  NOT NULL,
    created_at TIMESTAMPTZ   DEFAULT NOW(),
    updated_at TIMESTAMPTZ   DEFAULT NOW()
);

-- Projects table
CREATE TABLE projects (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_token VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    owner_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- Paths table
-- path is unique PER PROJECT, not globally (e.g. two projects can both have "/login")
CREATE TABLE paths (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path       VARCHAR(2048) NOT NULL,
    project_id UUID          NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ   DEFAULT NOW(),
    updated_at TIMESTAMPTZ   DEFAULT NOW(),
    UNIQUE (project_id, path)   -- was wrongly UNIQUE(path) globally before
);

-- Visitors table
-- One row per (path, session). visit_count increments on every hit.
-- All stats are derived live from this table — no aggregation needed.
CREATE TABLE visitors (
    id          UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    path_id     UUID              NOT NULL REFERENCES paths(id) ON DELETE CASCADE,
    session_id  VARCHAR(255)      NOT NULL,
    ip_address  VARCHAR(45)       NOT NULL,
    user_agent  TEXT,
    country     VARCHAR(100),
    first_visit TIMESTAMPTZ       DEFAULT NOW(),
    last_visit  TIMESTAMPTZ       DEFAULT NOW(),
    visit_count INTEGER           DEFAULT 1,   -- was INT32 (invalid in PostgreSQL)
    UNIQUE (path_id, session_id)
);

-- Sessions table (authentication)
CREATE TABLE sessions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT        UNIQUE NOT NULL,
    user_agent         TEXT,
    ip_address         VARCHAR(45),
    location           VARCHAR(255),
    expires_at         TIMESTAMPTZ NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW(),
    revoked_at         TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_projects_owner_id            ON projects(owner_id);
CREATE INDEX idx_paths_project_id             ON paths(project_id);
CREATE INDEX idx_visitors_path_id             ON visitors(path_id);
CREATE INDEX idx_visitors_first_visit         ON visitors(first_visit);   -- used in all date-range filters
CREATE INDEX idx_sessions_user_id             ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token_hash  ON sessions(refresh_token_hash);