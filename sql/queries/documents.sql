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

-- name: GetDocumentsByOwner :many
SELECT id, name FROM documents
WHERE owner_id = $1 ;


-- name: GetDocumentsByUser :many
SELECT d.id, d.name FROM documents d
LEFT JOIN document_permissions p
ON p.document_id = d.id 
WHERE p.user_id = $1 OR d.owner_id = $1; 




-- name: UpdateDocument :exec
UPDATE documents 
SET updated_at = NOW()
WHERE id = $1;

-- name: UpdateDocumentName :exec
UPDATE documents SET updated_at = NOW(), name = $1
WHERE id = $2;


-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;


-- name: GetDocumentOwnerId :one
SELECT owner_id from documents WHERE id = $1 ;



-- name: DeleteDocument :exec
DELETE from documents WHERE id = $1 ;


-- name: GetDocumentOwner :one
SELECT u.email, u.id 
FROM documents d
INNER JOIN users u 
ON d.owner_id = u.id
WHERE d.id = $1 ;