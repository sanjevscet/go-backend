package main

import (
	"net/http"

	"github.com/sanjevscet/go-backend.git/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := store.PaginateFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	parsedFeedQuery, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(parsedFeedQuery); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	ctx := r.Context()

	posts, err := app.store.Posts.GetUserFeed(ctx, fq, 4)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}
}
