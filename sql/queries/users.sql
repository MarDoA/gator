-- name: CreateUser :one
insert into users (id, created_at,updated_at,name)
values (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
select * from users where name = $1;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: DeleteAllUsers :exec
delete from users;

-- name: GetUsers :many
select * from users;