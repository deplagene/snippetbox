-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id uuid PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;