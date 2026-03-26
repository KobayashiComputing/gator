-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not null,
    name TEXT not NULL
);

-- +goose Down
DROP TABLE users;