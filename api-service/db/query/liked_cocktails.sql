-- name: LikeCocktail :one 
INSERT INTO liked_cocktails (
  cocktail_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetLikedCocktails :many
SELECt * from cocktails c
JOIN  liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1
GROUP BY lc.user_id;

-- name: GetLikedCocktail :one
SELECT * from cocktails c
JOIN liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1 and lc.cocktail_id = $2
GROUP BY lc.user_id, lc.cocktail_id, c.id;

-- name: UnlikeCocktail :exec
DELETE FROM liked_cocktails 
WHERE cocktail_id = $1 AND user_id = $2;
