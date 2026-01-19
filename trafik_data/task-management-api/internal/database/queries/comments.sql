-- name: CreateComment :one
INSERT INTO comments (todo_id, content)
VALUES ($1, $2)
RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM comments
WHERE id = $1;

-- name: ListCommentsByTodoID :many
SELECT * FROM comments
WHERE todo_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateComment :one
UPDATE comments
SET content = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

-- name: CountCommentsByTodoID :one
SELECT COUNT(*) FROM comments
WHERE todo_id = $1;
