-- name: CreatePasswordToken :one
INSERT INTO password_tokens (token, user_id)
VALUES (sqlc.arg(token)::text, sqlc.arg(user_id)::bigint)
RETURNING *;

-- name: GetPasswordToken :one
SELECT * FROM password_tokens
WHERE id = sqlc.arg(id)::bigint
  AND user_id = sqlc.arg(user_id)::bigint
  AND created_at >= sqlc.arg(created_after)::timestamptz;

-- name: DeletePasswordTokensByUser :exec
DELETE FROM password_tokens
WHERE user_id = $1;

-- name: UpdatePasswordTokenCreatedAt :execrows
UPDATE password_tokens
SET created_at = sqlc.arg(created_at)::timestamptz
WHERE id = sqlc.arg(id)::bigint;

-- name: CountPasswordTokensByUser :one
SELECT COUNT(*) FROM password_tokens
WHERE user_id = $1;
