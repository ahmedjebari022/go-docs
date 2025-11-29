-- name: CreateDocument :one
INSERT INTO documents (id, created_at, updated_at, owner_id, name)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3
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


-- name: GetDocumentOwner :one
SELECT owner_id from documents WHERE id = $1 ;



-- name: DeleteDocument :exec
DELETE from documents WHERE id = $1 ;