-- name: GetFeedByURL :one
select * from feeds where url = $1;
