package main

import (
	"errors"
	"fmt"
	"github.com/alisher0594/reddit-old/internal/data/entitys"
	"github.com/alisher0594/reddit-old/internal/validator"
	"net/http"
)

const (
	voteUp   = 1
	voteDown = -1
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Link        string `json:"link"`
		SubredditID *int64 `json:"subreddit_id"`
		Content     string `json:"content"`
		Promoted    bool   `json:"promoted"`
		NSFW        bool   `json:"nsfw"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &entitys.Post{
		Title:       input.Title,
		Author:      input.Author,
		Link:        input.Link,
		SubredditID: input.SubredditID,
		Content:     input.Content,
		Promoted:    input.Promoted,
		NSFW:        input.NSFW,
	}

	v := validator.New()
	if post.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Insert(r.Context(), post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/post/%d", post.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Posts.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) upPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	score, err := app.models.Posts.Vote(r.Context(), id, voteUp)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"score": score}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) downPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	score, err := app.models.Posts.Vote(r.Context(), id, voteDown)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"score": score}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	post, err := app.models.Posts.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Author      *string `json:"author"`
		Link        *string `json:"link"`
		SubredditID *int64  `json:"subreddit_id"`
		Content     *string `json:"content"`
		Promoted    *bool   `json:"promoted"`
		NSFW        *bool   `json:"nsfw"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Author != nil {
		post.Author = *input.Author
	}
	if input.Link != nil {
		post.Link = *input.Link
	}
	if input.SubredditID != nil {
		post.SubredditID = input.SubredditID
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.Promoted != nil {
		post.Promoted = *input.Promoted
	}
	if input.NSFW != nil {
		post.NSFW = *input.NSFW
	}

	v := validator.New()

	if post.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Update(r.Context(), post)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		entitys.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 27, v)
	input.Filters.Sort = app.readString(qs, "sort", "-score")
	input.Filters.SortSafelist = []string{"id", "score", "-score"}

	if input.Filters.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	posts, metadata, err := app.models.Posts.GetAll(r.Context(), input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"posts": posts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Posts.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entitys.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "post successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
