-- name: CreatePost :one
insert into posts (id, created_at,updated_at,title,url,description,published_at,feed_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
select p.*
from posts p 
join feed_follows f on p.feed_id = f.feed_id 
where f.user_id = $1 
order by p.created_at desc 
limit $2;