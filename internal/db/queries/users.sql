-- name: CreateUser :one
INSERT INTO users(
    user_name,
    name,
    email,
    password_hash
) VALUES (
    $1, $2, $3, $4
) RETURNING user_id, user_name, email, name;


-- name: FindUserByID :one
SELECT user_id, user_name, email, name
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
RETURNING user_id, user_name, name, email, updated_at;

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
RETURNING user_id, user_name, email, name, updated_at;

-- name: ListUsers :many
SELECT user_id,user_name, email, name
FROM users
WHERE
    deleted_at IS NULL
    AND (
        sqlc.narg(search)::TEXT IS NULL OR
        user_name ILIKE '%' || sqlc.narg(search) || '%' OR
        email ILIKE '%' || sqlc.narg(search) || '%' OR
        name ILIKE '%' || sqlc.narg(search) || '%'
    )
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'name'  AND sqlc.arg(order_by)::text = 'asc'  THEN name       END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'name'  AND sqlc.arg(order_by)::text = 'desc' THEN name       END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'email' AND sqlc.arg(order_by)::text = 'asc'  THEN email      END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'email' AND sqlc.arg(order_by)::text = 'desc' THEN email      END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'created_at' AND sqlc.arg(order_by)::text = 'asc' THEN created_at END ASC,
    created_at DESC
LIMIT sqlc.arg('limit_val')::int
OFFSET sqlc.arg('offset_val')::int;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE
    deleted_at IS NULL
    AND (
        sqlc.narg('search')::text IS NULL OR 
        sqlc.narg('search')::text = '' OR
        user_name ILIKE '%' || sqlc.narg('search') || '%' OR
        email ILIKE '%' || sqlc.narg('search') || '%' OR
        name ILIKE '%' || sqlc.narg('search') || '%'
    );

-- name: GetByEmail :one
SELECT user_id, user_name, name, email, password_hash
FROM users 
WHERE EMAIL = sqlc.arg(email) and deleted_at IS NULL;