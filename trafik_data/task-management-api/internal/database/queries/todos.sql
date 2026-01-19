-- name: CreateTodo :one
INSERT INTO todos (title, description, status, priority, category_id, due_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTodoByID :one
SELECT t.*, c.name as category_name
FROM todos t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.id = $1;

-- name: ListTodos :many
SELECT t.*, c.name as category_name
FROM todos t
LEFT JOIN categories c ON t.category_id = c.id
ORDER BY t.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListTodosByStatus :many
SELECT t.*, c.name as category_name
FROM todos t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.status = $1
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTodosByCategory :many
SELECT t.*, c.name as category_name
FROM todos t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.category_id = $1
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTodo :one
UPDATE todos
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    priority = COALESCE($5, priority),
    category_id = COALESCE($6, category_id),
    due_date = COALESCE($7, due_date),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTodoStatus :one
UPDATE todos
SET status = $2,
    completed_at = CASE WHEN $2 = 'completed' THEN CURRENT_TIMESTAMP ELSE NULL END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: CountTodos :one
SELECT COUNT(*) FROM todos;

-- name: CountTodosByStatus :one
SELECT COUNT(*) FROM todos
WHERE status = $1;
