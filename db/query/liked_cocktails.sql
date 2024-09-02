-- name: LikeCocktail :one 
INSERT INTO liked_cocktails (
  cocktail_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: UnlikeCocktail :exec
DELETE FROM liked_cocktails 
WHERE cocktail_id = $1 AND user_id = $2;


