-- name: GetPostsForUser :many
with feed_list as (
        select 
            feed_follows.*, 
            feeds.name as feed_name,
            users.name as user_name
        from feed_follows
        inner join feeds on feeds.id = feed_follows.feeds_id
        inner join users on users.id = feed_follows.users_id
        where feed_follows.users_id=$1
)

select 
        posts.id as post_id, 
        posts.published_at as pub_date, 
        posts.title as title, 
        posts.url as url,
        posts.description as description 
    from feed_list
    inner join posts on posts.feeds_id = feed_list.feeds_id
    order by posts.published_at DESC
    limit $2
;
