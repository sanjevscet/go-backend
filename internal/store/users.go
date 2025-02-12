package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("Duplicate user email")
	ErrDuplicateUsername = errors.New("Duplicate user name")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, password, email)
	VALUES($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query,
		user.Username,
		user.Password.hash,
		user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {

	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		// create user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			log.Println(err)
			return err
		}
		log.Println("no error 111")

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userId int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) values($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	log.Println(query, token, userId, time.Now().Add(invitationExp))

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.is_active 
	FROM users u 
	JOIN user_invitations ui on u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users set username = $1, email = $2, is_active = $3 WHERE id =  $4
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `
		DELETE FROM user_invitations WHERE user_id =  $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find the user belongs
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// get user by Id
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// remove entry from user_invitations table
		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}
