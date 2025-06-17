-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) 
    REFERENCES users(id) 
    ON DELETE CASCADE
);

ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
    ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);

ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;
DROP TABLE refresh_tokens;
ALTER TABLE users DROP COLUMN hashed_password;
DROP TABLE chirps;
DROP TABLE users;