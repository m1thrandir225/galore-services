-- name: LikeFlavour :one
INSERT INTO liked_flavours (
  flavour_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetLikedFlavour :one 
select * from flavours f 
join liked_flavours lf ON f.id = lf.flavour_id 
where lf.user_id = $1 and lf.flavour_id = $2
group by lf.flavour_id, lf.user_id, f.id;


-- name: UnlikeFlavour :exec
DELETE FROM liked_flavours 
WHERE flavour_id = $1 AND user_id = $2; 
