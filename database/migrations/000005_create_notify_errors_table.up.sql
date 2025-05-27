BEGIN;

CREATE TABLE notification_errors (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    body JSONB NOT NULL,
    error TEXT NOT NULL
);

COMMIT;
