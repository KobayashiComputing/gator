-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not null,
    url TEXT UNIQUE NOT NULL,
    name TEXT not NULL,
    users_id UUID NOT NULL,
    FOREIGN KEY (users_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
