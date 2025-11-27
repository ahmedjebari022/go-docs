-- name: CreateDocument :one
INSERT INTO documents (id, created_at, updated_at, owner_id, name, document_url)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;