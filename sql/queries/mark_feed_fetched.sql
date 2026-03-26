-- name: MarkFeedFetched :one
update feeds
set  last_fetched_at = LOCALTIMESTAMP
where id = $1
returning *;
