-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = LOWER(sqlc.arg(email)::text);

-- name: CreateUser :one
INSERT INTO users (name, email, password, verified, admin)
VALUES (
    sqlc.arg(name)::text,
    LOWER(sqlc.arg(email)::text),
    sqlc.arg(password)::text,
    sqlc.arg(verified)::boolean,
    sqlc.arg(admin)::boolean
)
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password = sqlc.arg(password)::text
WHERE id = sqlc.arg(id)::bigint
RETURNING *;

-- name: SetUserVerified :one
UPDATE users
SET verified = TRUE
WHERE id = $1
RETURNING *;

-- name: SetUserAdmin :one
UPDATE users
SET admin = sqlc.arg(admin)::boolean
WHERE id = sqlc.arg(id)::bigint
RETURNING *;
