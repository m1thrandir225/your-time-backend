-- name: CreateTask :one
INSERT INTO tasks (
    title,
    description,
    due_date,
    reminder_date,
    user_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks 
WHERE id = $1 LIMIT 1;

-- name: GetTasksByUser :many
SELECT * FROM tasks 
WHERE user_id = $1;
