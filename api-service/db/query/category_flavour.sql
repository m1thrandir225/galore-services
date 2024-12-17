-- name: CreateCategoryFlavour :one
INSERT INTO category_flavour (
  category_id,
  flavour_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetCategoryFlavour :one
SELECT * FROM category_flavour
WHERE id = $1 LIMIT 1;

-- name: GetCategoriesFromLikedFlavours :many
SELECT DISTINCT c.id, c.name, c.tag, c.created_at
FROM liked_flavours lf
JOIN category_flavour cf ON lf.flavour_id = cf.flavour_id
JOIN categories c ON cf.category_id = c.id
WHERE lf.user_id = $1;

-- name: DeleteCategoryFlavour :exec
DELETE from category_flavour
WHERE id = $1;
