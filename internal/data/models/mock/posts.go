package mock

import (
	"context"
	"github.com/alisher0594/reddit-old/internal/data/entitys"
)

// PostModel ...
type PostModel struct{}

// Insert ...
func (p PostModel) Insert(ctx context.Context, post *entitys.Post) error {
	return nil
}

// Get ...
func (p PostModel) Get(ctx context.Context, id int64) (*entitys.Post, error) {
	return &entitys.Post{
		ID:       13,
		Title:    "Mocked post",
		Author:   "t2_fed2ere13",
		Link:     "https://old.reddit.com/user/fed2ere13",
		Content:  "content of mocked post",
		Score:    100500,
		Promoted: false,
		NSFW:     false,
		Version:  1,
	}, nil
}

// Update ...
func (p PostModel) Update(ctx context.Context, post *entitys.Post) error {
	return nil
}

// Vote ...
func (p PostModel) Vote(ctx context.Context, id int64, vote int64) (int64, error) {
	return 0, nil
}

// Delete ...
func (p PostModel) Delete(ctx context.Context, id int64) error {
	return nil
}

// GetAll ...
func (p PostModel) GetAll(ctx context.Context, filters entitys.Filters) ([]*entitys.Post, entitys.Metadata, error) {
	_, err := p.getAllPromoted(ctx, filters)
	if err != nil {
		return nil, entitys.Metadata{}, err
	}

	return nil, entitys.Metadata{}, nil
}

func (p PostModel) getAllPromoted(ctx context.Context, filters entitys.Filters) ([]*entitys.Post, error) {
	return nil, nil
}
