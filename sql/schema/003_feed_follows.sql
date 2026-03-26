-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not null,
    users_id UUID NOT NULL,
    FOREIGN KEY (users_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    feeds_id UUID NOT NULL,
    FOREIGN KEY (feeds_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE,
    CONSTRAINT uq_users_feeds UNIQUE (users_id, feeds_id)
);

-- +goose Down
DROP TABLE feed_follows;
