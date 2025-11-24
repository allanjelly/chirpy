-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (gen_random_uuid(),NOW(),NOW(),$1,$2)
RETURNING *;

-- name: DeleteUsers :exec
DELETE from users;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email=$2, password=$3, updated_at=NOW() where id=$1
RETURNING *;

-- name: UpgradeUser_is_red :one
UPDATE users SET is_chirpy_red = true where id=$1 RETURNING *;