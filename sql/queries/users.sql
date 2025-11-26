-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;


-- name: UpdateUsersPassword :exec
UPDATE users
SET hashed_password = $1 
WHERE id = $2 ;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 ;

