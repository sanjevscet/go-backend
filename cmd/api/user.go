package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sanjevscet/go-backend.git/internal/store"
)

var userKey string = "user"

type createUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=40"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=40"`
}

type FollowerUser struct {
	UserID int64 `json:"user_id"`
}

// func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload createUserPayload

// 	if err := readJson(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	user := &store.User{
// 		Username: payload.Username,
// 		Email:    payload.Email,
// 		Password: payload.Password,
// 	}

// 	if err := app.store.Users.Create(context.Background(), user); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// 	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*store.User)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r.Context())

	// TODO: revert back to auth user id once auth is implemented
	var payload FollowerUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ForeignKeyViolated):
		case errors.Is(err, store.CheckConstraintViolated):
		case errors.Is(err, store.UniqueKeyConstraintViolated):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UnFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowerUser := getUserFromContext(r.Context())

	// TODO: revert back to auth user id once auth is implemented
	var payload FollowerUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.UnFollow(ctx, unFollowerUser.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ForeignKeyViolated):
		case errors.Is(err, store.CheckConstraintViolated):
		case errors.Is(err, store.UniqueKeyConstraintViolated):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userId")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if err := app.store.Users.Activate(r.Context(), token); err != nil {

		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getUserFromContext(ctx context.Context) *store.User {
	return ctx.Value(userKey).(*store.User)
}
