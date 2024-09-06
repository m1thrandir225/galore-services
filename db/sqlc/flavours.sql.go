// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: flavours.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createFlavour = `-- name: CreateFlavour :one
INSERT INTO flavours (
  name
) VALUES (
  $1
) RETURNING id, name, created_at
`

func (q *Queries) CreateFlavour(ctx context.Context, name string) (Flavour, error) {
	row := q.db.QueryRow(ctx, createFlavour, name)
	var i Flavour
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const deleteFlavour = `-- name: DeleteFlavour :exec
DELETE FROM flavours 
WHERE id = $1
`

func (q *Queries) DeleteFlavour(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteFlavour, id)
	return err
}

const getAllFlavours = `-- name: GetAllFlavours :many
SELECT id, name, created_at FROM flavours
`

func (q *Queries) GetAllFlavours(ctx context.Context) ([]Flavour, error) {
	rows, err := q.db.Query(ctx, getAllFlavours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Flavour{}
	for rows.Next() {
		var i Flavour
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

const getFlavourId = `-- name: GetFlavourId :one
SELECT id, name, created_at FROM flavours
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFlavourId(ctx context.Context, id uuid.UUID) (Flavour, error) {
	row := q.db.QueryRow(ctx, getFlavourId, id)
	var i Flavour
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const getFlavourName = `-- name: GetFlavourName :one
SELECT id, name, created_at FROM flavours
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetFlavourName(ctx context.Context, name string) (Flavour, error) {
	row := q.db.QueryRow(ctx, getFlavourName, name)
	var i Flavour
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const updateFlavour = `-- name: UpdateFlavour :one
UPDATE flavours
SET name = $2
WHERE id = $1
RETURNING id, name, created_at
`

type UpdateFlavourParams struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (q *Queries) UpdateFlavour(ctx context.Context, arg UpdateFlavourParams) (Flavour, error) {
	row := q.db.QueryRow(ctx, updateFlavour, arg.ID, arg.Name)
	var i Flavour
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}
