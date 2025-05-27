BEGIN;

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR NOT NULL,
    body JSONB NOT NULL
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR NOT NULL,
    admin_user VARCHAR NOT NULL
);

COMMIT;
