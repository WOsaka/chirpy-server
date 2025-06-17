-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(), 
    NOW(), 
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: SetPassword :exec
UPDATE users
SET hashed_password = $1
WHERE email = $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at) 
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;

-- name: UpdateUserCredentials :one
UPDATE users
SET email = $1,
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE id = $1;

-- name: SetChirpyRedByID :exec
UPDATE users 
SET is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;