-- name: CreateUser :one
INSERT INTO users(
    user_name,
    name,
    email,
    password_hash
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: FindUserByID :one
SELECT *
FROM users
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: UpdateUserByID :one
UPDATE users
SET 
    name          = COALESCE(sqlc.narg(name), name),
    email         = COALESCE(sqlc.narg(email), email),
    user_name     = COALESCE(sqlc.narg(user_name), user_name),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    updated_at    = NOW()
WHERE user_id = sqlc.arg(user_id) AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :execrows
UPDATE users
SET deleted_at = NOW()
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: HardDeleteUser :execrows
DELETE FROM users
WHERE user_id = $1;

-- name: RestoreUser :one
UPDATE users
SET deleted_at = NULL
WHERE user_id = sqlc.arg(user_id)
  AND deleted_at IS NOT NULL
RETURNING *;

