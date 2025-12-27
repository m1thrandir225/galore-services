-- name: LikeCocktail :one
INSERT INTO liked_cocktails (
  cocktail_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetLikedCocktails :many
SELECT c.id,
        c.name,
        c.is_alcoholic,
        c.glass,
        c.image,
        c.embedding,
        c.instructions,
        c.ingredients,
        c.created_at
from cocktails c
JOIN  liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1
GROUP BY lc.user_id, c.id, lc.id;

-- name: IsCocktailLiked :one
SELECT EXISTS (
    SELECT 1
    FROM liked_cocktails lc
    WHERE lc.user_id = $1 AND lc.cocktail_id = $2
) AS is_liked;

-- name: GetLikedCocktail :one
SELECT c.id,
        c.name,
        c.is_alcoholic,
        c.glass,
        c.image,
        c.embedding,
        c.instructions,
        c.ingredients,
        c.created_at
from cocktails c
JOIN liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1 and lc.cocktail_id = $2
GROUP BY lc.user_id, lc.cocktail_id, c.id, lc.id;

-- name: UnlikeCocktail :exec
DELETE FROM liked_cocktails
WHERE cocktail_id = $1 AND user_id = $2;
