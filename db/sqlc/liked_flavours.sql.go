// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: liked_flavours.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getLikedFlavour = `-- name: GetLikedFlavour :one
select id ,name, created_at from flavours f 
join liked_flavours lf ON f.id = lf.flavour_id 
where lf.user_id = $1 and lf.flavour_id = $2
group by lf.flavour_id, lf.user_id, f.id
`

type GetLikedFlavourParams struct {
	UserID    uuid.UUID `json:"user_id"`
	FlavourID uuid.UUID `json:"flavour_id"`
}

func (q *Queries) GetLikedFlavour(ctx context.Context, arg GetLikedFlavourParams) (Flavour, error) {
	row := q.db.QueryRow(ctx, getLikedFlavour, arg.UserID, arg.FlavourID)
	var i Flavour
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const getUserLikedFlavours = `-- name: GetUserLikedFlavours :many
SELECT id, name, created_at from flavours f 
JOIN liked_flavours lf ON f.id = lf.flavour_id
WHERE lf.user_id = $1
GROUP BY lf.flavour_id, lf.user_id, f.id
`

func (q *Queries) GetUserLikedFlavours(ctx context.Context, userID uuid.UUID) ([]Flavour, error) {
	rows, err := q.db.Query(ctx, getUserLikedFlavours, userID)
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

const likeFlavour = `-- name: LikeFlavour :one
INSERT INTO liked_flavours (
  flavour_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING flavour_id, user_id
`

type LikeFlavourParams struct {
	FlavourID uuid.UUID `json:"flavour_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) LikeFlavour(ctx context.Context, arg LikeFlavourParams) (LikedFlavour, error) {
	row := q.db.QueryRow(ctx, likeFlavour, arg.FlavourID, arg.UserID)
	var i LikedFlavour
	err := row.Scan(&i.FlavourID, &i.UserID)
	return i, err
}

const unlikeFlavour = `-- name: UnlikeFlavour :exec
DELETE FROM liked_flavours 
WHERE flavour_id = $1 AND user_id = $2
`

type UnlikeFlavourParams struct {
	FlavourID uuid.UUID `json:"flavour_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) UnlikeFlavour(ctx context.Context, arg UnlikeFlavourParams) error {
	_, err := q.db.Exec(ctx, unlikeFlavour, arg.FlavourID, arg.UserID)
	return err
}
