-- +goose Up
-- +goose StatementBegin
CREATE TYPE todo_status AS ENUM ('pending', 'in_progress', 'completed');
CREATE TYPE todo_priority AS ENUM ('low', 'medium', 'high');

CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status todo_status NOT NULL DEFAULT 'pending',
    priority todo_priority NOT NULL DEFAULT 'medium',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_priority ON todos(priority);
CREATE INDEX idx_todos_category_id ON todos(category_id);
CREATE INDEX idx_todos_due_date ON todos(due_date);
CREATE INDEX idx_todos_created_at ON todos(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todos;
DROP TYPE IF EXISTS todo_status;
DROP TYPE IF EXISTS todo_priority;
-- +goose StatementEnd
