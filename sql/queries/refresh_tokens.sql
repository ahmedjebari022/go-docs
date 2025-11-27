-- name: CreateToken :one
INSERT INTO refresh_tokens (token, expires_at, user_id)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;


-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;


