package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound                 = errors.New("Resource not found")
	QueryTimeOutDuration        = time.Second * 5
	ForeignKeyViolated          = errors.New("Referenced key not found")
	CheckConstraintViolated     = errors.New("Check constraint violated")
	UniqueKeyConstraintViolated = errors.New("Unique key constraint violated")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		DeleteById(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
	}
	Comments interface {
		GetByPostId(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
