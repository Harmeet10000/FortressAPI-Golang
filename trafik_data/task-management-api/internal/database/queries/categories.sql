-- name: CreateCategory :one
INSERT INTO categories (name, description, color)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    color = COALESCE($4, color),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;
