-- name: GenerateFeaturedForToday :many
WITH today_featured AS (
    SELECT COUNT(*) AS count
    FROM daily_featured_cocktails
    WHERE DATE(created_at) = CURRENT_DATE
),
     not_recently_featured AS (
         SELECT c.id
         FROM cocktails c
         WHERE NOT EXISTS (
             SELECT 1
             FROM daily_featured_cocktails df
             WHERE df.cocktail_id = c.id
               AND df.created_at > CURRENT_DATE - INTERVAL '7 days'
         )
     )
INSERT INTO daily_featured_cocktails (cocktail_id)
SELECT random_cocktails.id
FROM (
         SELECT nrf.id
         FROM not_recently_featured nrf
         ORDER BY RANDOM()
         LIMIT 10
     ) random_cocktails
WHERE (SELECT count FROM today_featured) < 10
RETURNING *;

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