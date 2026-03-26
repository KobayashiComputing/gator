-- name: GetNextFeedToFetchSingle :one
select id, name, url, last_fetched_at from feeds order by last_fetched_at asc nulls first;