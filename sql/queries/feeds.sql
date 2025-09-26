-- name: CreateFeed :one
insert into feeds (id, created_at,updated_at,name,url,user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedByURL :one
select * from feeds where url = $1;

-- name: CreateFeedFollow :one
with inserted_feed_follow as (insert into feed_follows (id,created_at,updated_at,user_id,feed_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
)
select inserted_feed_follow.*,
f.name as feed_name,
u.name as user_name
from inserted_feed_follow
join feeds f on inserted_feed_follow.feed_id = f.id 
join users u on inserted_feed_follow.user_id = u.id ;

-- name: GetFeedFollowsForUser :many
select feed_follows.*,
f.name as feed_name,
u.name as user_name
from feed_follows
join feeds f on feed_follows.feed_id = f.id
join users u on feed_follows.user_id = u.id
where u.name = $1;

-- name: DeleteFeedFollowForUser :exec
delete from feed_follows where user_id = $1 and feed_id = $2;

-- name: MarkFeedFetched :exec
update feeds
set updated_at = $1, last_fetched_at = $1
where id = $2;

-- name: GetNextFeedToFetch :one
select * from feeds order by last_fetched_at asc nulls first limit 1;