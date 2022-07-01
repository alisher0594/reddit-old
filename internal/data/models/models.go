package models

import (
	"context"
	"database/sql"
	"github.com/alisher0594/reddit-old/internal/data/entitys"
	"github.com/alisher0594/reddit-old/internal/data/models/postges"
)

type Posts interface {
	Insert(ctx context.Context, movie *entitys.Post) error
	Get(ctx context.Context, id int64) (*entitys.Post, error)
	GetAll(ctx context.Context, filters entitys.Filters) ([]*entitys.Post, entitys.Metadata, error)
	Update(ctx context.Context, movie *entitys.Post) error
	Vote(ctx context.Context, id int64, vote int64) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type Models struct {
	Posts Posts
}

func New(db *sql.DB) Models {
	return Models{
		Posts: postges.PostModel{DB: db},
	}
}

//func NewMock() Models {
//	return Models{
//		Posts: mock.PostModel{},
//	}
//}
