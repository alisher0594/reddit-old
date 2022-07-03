package entitys

import (
	"fmt"
	"github.com/alisher0594/reddit-old/internal/validator"
	"strings"
	"time"
)

const prefix = "t2_"

// Post ...
type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Link        string    `json:"link"`
	SubredditID *int64    `json:"-"`
	Subreddit   *SubPost  `json:"subreddit,omitempty"`
	Content     string    `json:"content"`
	Score       int64     `json:"score"`
	Promoted    bool      `json:"promoted"`
	NSFW        bool      `json:"nsfw"`
	Version     int32     `json:"version"`
}

// Validate ...
func (p *Post) Validate(v *validator.Validator) {
	v.Check(p.Title != "", "title", "must be provided")
	v.Check(len(p.Title) <= 500, "title", "must not be more than 500 bytes long")
	validateAuthor(v, p.Author)

	v.Check(validator.IsValidLink(p.Link), "link", "must be a valid URL")

	v.Check(p.Content != "", "content", "must be provided")
	//v.Check(validator.Matches(p.Content, validator.LinkRX), "content", "cannot have link")
}

func validateAuthor(v *validator.Validator, author string) {
	v.Check(author != "", "author", "must be provided")
	v.Check(strings.HasPrefix(author, prefix), "author", fmt.Sprintf("must have prefix: %s", prefix))
	v.Check(len(author) == 11, "author", "must be 8 bytes long")
	v.Check(validator.IsLowercase(author), "author", "should only contain lowercase letters and numbers")
}

// SubPost ...
type SubPost struct {
	ID        *int64     `json:"id"`
	CreatedAt *time.Time `json:"-"`
	Title     *string    `json:"title"`
	Author    *string    `json:"author"`
	Link      *string    `json:"link"`
	Content   *string    `json:"content"`
	Score     *int64     `json:"score"`
	Promoted  *bool      `json:"promoted"`
	NSFW      *bool      `json:"nsfw"`
	Version   *int32     `json:"version"`
}
