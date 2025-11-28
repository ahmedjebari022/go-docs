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

-- name: GetDocumentsByUser :many
SELECT id, name FROM documents
WHERE owner_id = $1 ;



-- name: UpdateDocument :exec
UPDATE documents 
SET updated_at = NOW()
WHERE id = $1;

-- name: UpdateDocumentName :exec
UPDATE documents SET updated_at = NOW(), name = $1
WHERE id = $2;


-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;