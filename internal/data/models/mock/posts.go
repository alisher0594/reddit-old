package mock

import (
	"github.com/alisher0594/reddit-old/internal/data/entitys"
)

// PostModel ...
type PostModel struct{}

func (m PostModel) Insert(movie *entitys.Post) error {
	return nil
}

func (m PostModel) Get(id int64) (*entitys.Post, error) {
	return nil, nil
}

func (m PostModel) Update(movie *entitys.Post) error {
	return nil
}

func (m PostModel) Delete(id int64) error {
	return nil
}
