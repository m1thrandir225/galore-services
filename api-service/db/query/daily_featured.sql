-- name: CreateDailyFeatured :one
INSERT INTO daily_featured_cocktails (
    cocktail_id
) VALUES (
 $1
) RETURNING *;

-- name: CheckWasCocktailFeatured :one
SELECT c.id
FROM cocktails c
JOIN daily_featured_cocktails dfc ON dfc.cocktail_id = c.id
WHERE dfc.created_at >= CURRENT_DATE - INTERVAL '7 days' AND c.id = $1;


-- name: GetDailyFeatured :many
SELECT c.id,
       c.name,
       c.is_alcoholic,
       c.glass,
       c.image,
       c.embedding,
       c.instructions,
       c.ingredients,
       c.created_at
FROM cocktails c
JOIN daily_featured_cocktails dfc ON dfc.cocktail_id = c.id
WHERE dfc.created_at >= CURRENT_DATE
  AND dfc.created_at < CURRENT_DATE + INTERVAL '1 day';

-- name: DeleteOlderFeatured :exec
DELETE FROM daily_featured_cocktails
WHERE created_at < NOW() - INTERVAL '7 days';