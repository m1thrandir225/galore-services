// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: category_flavour.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createCategoryFlavour = `-- name: CreateCategoryFlavour :one
INSERT INTO category_flavour (
  category_id,
  flavour_id
) VALUES (
  $1,
  $2
) RETURNING id, category_id, flavour_id
`

type CreateCategoryFlavourParams struct {
	CategoryID uuid.UUID `json:"category_id"`
	FlavourID  uuid.UUID `json:"flavour_id"`
}

func (q *Queries) CreateCategoryFlavour(ctx context.Context, arg CreateCategoryFlavourParams) (CategoryFlavour, error) {
	row := q.db.QueryRow(ctx, createCategoryFlavour, arg.CategoryID, arg.FlavourID)
	var i CategoryFlavour
	err := row.Scan(&i.ID, &i.CategoryID, &i.FlavourID)
	return i, err
}

const deleteCategoryFlavour = `-- name: DeleteCategoryFlavour :exec
DELETE from category_flavour
WHERE id = $1
`

func (q *Queries) DeleteCategoryFlavour(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCategoryFlavour, id)
	return err
}

const getCategoriesFromLikedFlavours = `-- name: GetCategoriesFromLikedFlavours :many
SELECT DISTINCT c.id, c.name, c.tag, c.created_at
FROM liked_flavours lf
JOIN category_flavour cf ON lf.flavour_id = cf.flavour_id
JOIN categories c ON cf.category_id = c.id
WHERE lf.user_id = $1
`

func (q *Queries) GetCategoriesFromLikedFlavours(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	rows, err := q.db.Query(ctx, getCategoriesFromLikedFlavours, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Tag,
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

const getCategoryFlavour = `-- name: GetCategoryFlavour :one
SELECT id, category_id, flavour_id FROM category_flavour
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCategoryFlavour(ctx context.Context, id uuid.UUID) (CategoryFlavour, error) {
	row := q.db.QueryRow(ctx, getCategoryFlavour, id)
	var i CategoryFlavour
	err := row.Scan(&i.ID, &i.CategoryID, &i.FlavourID)
	return i, err
}
