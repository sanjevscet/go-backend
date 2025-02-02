package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowedID int64  `json:"followed_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerId, userId int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503": // foreign key violation
				return ForeignKeyViolated

			case "23505": // unique key constraint key violation
				return UniqueKeyConstraintViolated

			case "23514": // check constraint key violation
				return CheckConstraintViolated
			default:
				return err
			}
		}
	}
	return err
}

func (s *FollowerStore) UnFollow(ctx context.Context, followerId, userId int64) error {
	query := `
		DELETE FROM followers WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	return err
}
