-- name: DeleteFeedFollow :exec
delete from feed_follows where users_id = $1 and feeds_id = $2;
