-- name: CreateCocktailCategory :one 
INSERT INTO cocktail_categories (
  cocktail_id,
  category_id 
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetCategoriesForCocktail :many 
SELECT * from categories c
JOIN cocktail_categories cc ON cc.category_id = c.id 
WHERE cc.cocktail_id = $1
GROUP BY cc.user_id;

-- name: GetCocktailsForCategory :many
SELECT * FROM cocktails c 
JOIN cocktail_categories cc on cc.cocktail_id 

-- name: DeleteCocktailCategory :exec
DELETE FROM cocktail_categories
WHERE id = $1;
