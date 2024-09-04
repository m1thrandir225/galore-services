// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: cocktail_categories.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createCocktailCategory = `-- name: CreateCocktailCategory :one
INSERT INTO cocktail_categories (
  cocktail_id,
  category_id 
) VALUES (
  $1,
  $2
) RETURNING id, cocktail_id, category_id
`

type CreateCocktailCategoryParams struct {
	CocktailID uuid.UUID `json:"cocktail_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (q *Queries) CreateCocktailCategory(ctx context.Context, arg CreateCocktailCategoryParams) (CocktailCategory, error) {
	row := q.db.QueryRow(ctx, createCocktailCategory, arg.CocktailID, arg.CategoryID)
	var i CocktailCategory
	err := row.Scan(&i.ID, &i.CocktailID, &i.CategoryID)
	return i, err
}

const deleteCocktailCategory = `-- name: DeleteCocktailCategory :exec
DELETE FROM cocktail_categories
WHERE id = $1
`

func (q *Queries) DeleteCocktailCategory(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCocktailCategory, id)
	return err
}

const getCategoriesForCocktail = `-- name: GetCategoriesForCocktail :many
SELECT c.id, name, created_at from categories c
JOIN cocktail_categories cc ON cc.category_id = c.id 
WHERE cc.cocktail_id = $1
GROUP BY cc.cocktail_id, cc.category_id, c.id
`

func (q *Queries) GetCategoriesForCocktail(ctx context.Context, cocktailID uuid.UUID) ([]Category, error) {
	rows, err := q.db.Query(ctx, getCategoriesForCocktail, cocktailID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCocktailCategory = `-- name: GetCocktailCategory :one
SELECT id, cocktail_id, category_id from cocktail_categories
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCocktailCategory(ctx context.Context, id uuid.UUID) (CocktailCategory, error) {
	row := q.db.QueryRow(ctx, getCocktailCategory, id)
	var i CocktailCategory
	err := row.Scan(&i.ID, &i.CocktailID, &i.CategoryID)
	return i, err
}

const getCocktailsForCategory = `-- name: GetCocktailsForCategory :many
SELECT c.id, name, is_alcoholic, glass, image, instructions, ingredients, created_at FROM cocktails c 
JOIN cocktail_categories cc on cc.cocktail_id = c.id
WHERE cc.category_id = $1
GROUP BY cc.category_id, cc.cocktail_id, c.id
`

func (q *Queries) GetCocktailsForCategory(ctx context.Context, categoryID uuid.UUID) ([]Cocktail, error) {
	rows, err := q.db.Query(ctx, getCocktailsForCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cocktail{}
	for rows.Next() {
		var i Cocktail
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsAlcoholic,
			&i.Glass,
			&i.Image,
			&i.Instructions,
			&i.Ingredients,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
