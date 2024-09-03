-- name: CreateCocktailCategory :one 
INSERT INTO cocktail_categories (
  cocktail_id,
  category_id 
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetCategoriesForCocktail :many 
SELECT c.id, name, created_at from categories c
JOIN cocktail_categories cc ON cc.category_id = c.id 
WHERE cc.cocktail_id = $1
GROUP BY cc.cocktail_id, cc.category_id, c.id;

-- name: GetCocktailsForCategory :many
SELECT c.id, name, is_alcoholic, glass, image, instructions, ingredients, created_at FROM cocktails c 
JOIN cocktail_categories cc on cc.cocktail_id = c.id
WHERE cc.category_id = $1
GROUP BY cc.category_id, cc.cocktail_id, c.id;

-- name: DeleteCocktailCategory :exec
DELETE FROM cocktail_categories cc
WHERE cc.id = $1;
