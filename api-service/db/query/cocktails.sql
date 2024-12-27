-- name: CreateCocktail :one
INSERT INTO cocktails (
  name,
  is_alcoholic,
  glass,
  image,
  instructions,
  ingredients,
  embedding
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
) RETURNING *;

-- name: GetCocktail :one
SELECT * FROM cocktails
WHERE id = $1 LIMIT 1;

-- name: GetCocktailAndSimilar :many
WITH target_cocktail AS (
    SELECT embedding
    FROM cocktails
    WHERE cocktails.id = $1
    LIMIT 1
)
SELECT c.*
FROM cocktails c, target_cocktail t
WHERE c.id != $1
ORDER BY c.embedding <=> t.embedding
LIMIT 10;


-- name: UpdateCocktail :one
UPDATE cocktails
SET name=$2, is_alcoholic=$3, glass=$4, image=$5, instructions=$6, ingredients=$7
WHERE id = $1
RETURNING *;

-- name: DeleteCocktail :exec 
DELETE FROM cocktails 
WHERE id = $1;

-- name: SearchCocktails :many
SELECT
    id,
    name,
    is_alcoholic,
    glass,
    image,
    instructions,
    ingredients,
    embedding,
    created_at
FROM cocktails
WHERE
    $1::TEXT IS NULL
   OR
    (name ILIKE '%' || $1::TEXT || '%')
   OR
    (ingredients::TEXT ILIKE '%' || $1::TEXT || '%');

-- name: GetHomescreenForUser :many
WITH UserLikedCategories AS (
    SELECT cf.category_id
    FROM liked_flavours lf
             JOIN category_flavour cf ON lf.flavour_id = cf.flavour_id
    WHERE lf.user_id = $1
    GROUP BY cf.category_id
    ORDER BY RANDOM()
    LIMIT 4
),
     CategoryCocktails AS (
         SELECT cc.cocktail_id, cc.category_id
         FROM cocktail_categories cc
         WHERE cc.category_id IN (SELECT category_id FROM UserLikedCategories)
     ),
     RankedCocktails AS (
         SELECT
             cc.cocktail_id,
             cc.category_id,
             ROW_NUMBER() OVER (PARTITION BY cc.category_id ORDER BY RANDOM()) AS rank
         FROM CategoryCocktails cc
     )
SELECT
    rc.category_id,
    array_agg(rc.cocktail_id)::uuid[] AS cocktails
FROM RankedCocktails rc
WHERE rc.rank <= 5
GROUP BY rc.category_id
HAVING COUNT(rc.cocktail_id) >= 2;
