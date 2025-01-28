CREATE TABLE IF NOT EXISTS recorder.claims (
    digest TEXT NOT NULL,
    headers JSONB NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
