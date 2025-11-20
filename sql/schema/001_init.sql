-- +goose Up
CREATE TABLE IF NOT EXISTS snippets (
    snippet_id uuid PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS snippets;