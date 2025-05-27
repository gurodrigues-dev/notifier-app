BEGIN;

CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(100) NOT NULL,
    target_id VARCHAR(100) NOT NULL,
    "group" VARCHAR(100) NOT NULL
);

ALTER TABLE tokens
ADD CONSTRAINT unique_admin_user UNIQUE (admin_user);

COMMIT;
