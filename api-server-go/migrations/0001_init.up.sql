-- Active: 1756357892734@@127.0.0.1@5432@minivercel
CREATE TYPE deployment_status AS ENUM (
  'NOT_STARTED',
  'QUEUED',
  'IN_PROGRESS',
  'READY',
  'FAIL'
);

CREATE TABLE projects (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid ()::text,
    name TEXT NOT NULL,
    git_url TEXT NOT NULL,
    subdomain TEXT NOT NULL,
    custom_domain TEXT,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_projects_user_id ON projects (user_id);

CREATE TABLE deployments (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid ()::text,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    status deployment_status NOT NULL DEFAULT 'NOT_STARTED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
