-- name: CreateCocktailCategory :one
INSERT INTO cocktail_categories (
  cocktail_id,
  category_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetCocktailCategory :one
SELECT * from cocktail_categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoriesForCocktail :many
SELECT c.id, name, tag, created_at from categories c
JOIN cocktail_categories cc ON cc.category_id = c.id
WHERE cc.cocktail_id = $1
GROUP BY cc.cocktail_id, cc.category_id, c.id;

-- name: GetCocktailsForCategory :many
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
JOIN cocktail_categories cc on cc.cocktail_id = c.id
WHERE cc.category_id = $1
GROUP BY cc.category_id, cc.cocktail_id, c.id;

-- name: DeleteCocktailCategory :exec
DELETE FROM cocktail_categories
WHERE id = $1;
