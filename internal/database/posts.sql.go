// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT posts.id, posts.created_at, posts.updated_at, title, url, description, published_at, posts.feed_id, feeds_follows.id, feeds_follows.created_at, feeds_follows.updated_at, user_id, feeds_follows.feed_id, users.id, users.created_at, users.updated_at, name FROM posts
LEFT JOIN feeds_follows ON posts.feed_id = feeds_follows.feed_id
LEFT JOIN users ON feeds_follows.user_id = $1
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostsForUserRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
	ID_2        uuid.NullUUID
	CreatedAt_2 sql.NullTime
	UpdatedAt_2 sql.NullTime
	UserID      uuid.NullUUID
	FeedID_2    uuid.NullUUID
	ID_3        uuid.NullUUID
	CreatedAt_3 sql.NullTime
	UpdatedAt_3 sql.NullTime
	Name        sql.NullString
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.UserID,
			&i.FeedID_2,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertPost = `-- name: InsertPost :one
INSERT INTO posts (    
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id 
    )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type InsertPostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) InsertPost(ctx context.Context, arg InsertPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, insertPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}
