-- name: CreatePermission :exec
INSERT INTO document_permissions (user_id, document_id, role)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUsersFromDocument :many
SELECT u.email ,u.id , d.role
FROM document_permissions d
INNER JOIN users u 
ON u.id = d.user_id
WHERE document_id = $1 ;


-- name: UpdatePermission :exec
UPDATE document_permissions SET role = $1 WHERE user_id = $2 AND document_id = $3 ;

-- name: DeletePermission :exec
DELETE from document_permissions WHERE user_id = $1 AND document_id = $2; 


-- name: GetUserPermission :one
SELECT role FROM document_permissions 
WHERE user_id = $1 AND document_id = $2;