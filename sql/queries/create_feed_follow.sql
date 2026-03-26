-- name: CreateFeedFollow :one
with inserted_feed_follow as (
        INSERT INTO feed_follows (
            id, 
            created_at, 
            updated_at, 
            users_id,
            feeds_id
        )
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
select 
    inserted_feed_follow.*, 
    feeds.name as feed_name,
    users.name as user_name
from inserted_feed_follow
inner join feeds on feeds.id = inserted_feed_follow.feeds_id
inner join users on users.id = inserted_feed_follow.users_id
;
