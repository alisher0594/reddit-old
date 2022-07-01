package entitys

import (
	"github.com/alisher0594/reddit-old/internal/validator"
	"testing"
)

func TestPostValidate(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		message string
		post    *Post
	}{
		{
			name: "valid Title",
			key:  "title",
			post: &Post{
				Title:   "test title",
				Author:  "t2_11qnzrqv",
				Link:    "https://google.com",
				Content: "some content text for testing",
			},
		},
		{
			name:    "invalid title",
			key:     "title",
			message: "must be provided",
			post: &Post{
				Title:   "",
				Author:  "t2_11qnzrqv",
				Link:    "https://google.com",
				Content: "some content text for testing",
			},
		},
		{
			name: "valid Author",
			key:  "author",
			post: &Post{
				Title:   "title",
				Author:  "t2_11qnzrqv",
				Link:    "https://google.com",
				Content: "some content text for testing",
			},
		},
		{
			name:    "invalid Author",
			key:     "author",
			message: "must have prefix: t2_",
			post: &Post{
				Title:   "title",
				Author:  "t@_11qnzrqv",
				Link:    "https://google.com",
				Content: "some content text for testing",
			},
		},
	}

	for _, tt := range testCases {
		v := validator.New()

		t.Run(tt.name, func(t *testing.T) {
			tt.post.Validate(v)
			if v.Errors[tt.key] != tt.message {
				t.Errorf("want: %v; got %v", tt.message, v.Errors[tt.key])
			}
		})
	}
}
