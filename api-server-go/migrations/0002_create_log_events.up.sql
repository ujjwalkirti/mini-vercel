-- Create log_events table matching ClickHouse schema
CREATE TABLE IF NOT EXISTS log_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    deployment_id TEXT NOT NULL,
    log TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Index for fast retrieval by deployment_id (most common query)
CREATE INDEX idx_log_events_deployment_id ON log_events (deployment_id);

-- Composite index for time-range queries
CREATE INDEX idx_log_events_deployment_timestamp ON log_events (deployment_id, timestamp DESC);

-- Foreign key constraint for referential integrity
ALTER TABLE log_events
ADD CONSTRAINT fk_log_events_deployment FOREIGN KEY (deployment_id) REFERENCES deployments (id) ON DELETE CASCADE;
