package postges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alisher0594/reddit-old/internal/data/entitys"
	"time"
)

type PostModel struct {
	DB *sql.DB
}

func (p PostModel) Insert(ctx context.Context, post *entitys.Post) error {
	query := `
        INSERT INTO posts (title, author, link, subreddit_id, content, promoted, nsfw) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, version`

	args := []interface{}{post.Title, post.Author, post.Link, post.SubredditID, post.Content, post.Promoted, post.NSFW}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.Version)
}

func (p PostModel) Get(ctx context.Context, id int64) (*entitys.Post, error) {
	if id < 1 {
		return nil, entitys.ErrRecordNotFound
	}

	query := `
		SELECT posts.id, posts.created_at, posts.title, posts.author, posts.link, posts.subreddit_id, posts.content,
		       posts.score, posts.promoted, posts.nsfw, posts.version, 
		       sub_post.id, sub_post.created_at, sub_post.title, sub_post.author, sub_post.link,
		       sub_post.content, sub_post.score, sub_post.promoted, sub_post.nsfw, sub_post.version
		FROM posts
		LEFT JOIN posts as sub_post ON sub_post.id = posts.subreddit_id
		WHERE posts.id = $1`

	var post entitys.Post
	var subPost entitys.SubPost
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Title,
		&post.Author,
		&post.Link,
		&post.SubredditID,
		&post.Content,
		&post.Score,
		&post.Promoted,
		&post.NSFW,
		&post.Version,

		&subPost.ID,
		&subPost.CreatedAt,
		&subPost.Title,
		&subPost.Author,
		&subPost.Link,
		&subPost.Content,
		&subPost.Score,
		&subPost.Promoted,
		&subPost.NSFW,
		&subPost.Version,
	)

	if post.SubredditID != nil {
		post.Subreddit = &subPost
	}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entitys.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (p PostModel) Update(ctx context.Context, post *entitys.Post) error {
	query := `
        UPDATE posts
        SET title = $1, author = $2, link = $3, subreddit_id = $4, content = $5, promoted = $6, nsfw = $7,
            version = version + 1
        WHERE id = $8 AND version = $9
        RETURNING version`

	args := []interface{}{
		post.Title,
		post.Author,
		post.Link,
		post.SubredditID,
		post.Content,
		post.Promoted,
		post.NSFW,
		post.ID,
		post.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return entitys.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (p PostModel) Vote(ctx context.Context, id int64, vote int64) (int64, error) {
	if id < 1 {
		return 0, entitys.ErrRecordNotFound
	}

	query := `
        UPDATE posts
        SET score = score + $1
        WHERE id = $2
		RETURNING score`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []interface{}{
		vote, id,
	}
	var score int64
	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&score)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, entitys.ErrEditConflict
		default:
			return 0, err
		}
	}

	return score, nil
}

func (p PostModel) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return entitys.ErrRecordNotFound
	}

	query := `
        DELETE FROM posts
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entitys.ErrRecordNotFound
	}

	return nil
}

func (p PostModel) GetAll(ctx context.Context, filters entitys.Filters) ([]*entitys.Post, entitys.Metadata, error) {
	promoted, err := p.getAllPromoted(ctx, filters)
	if err != nil {
		return nil, entitys.Metadata{}, err
	}

	query := fmt.Sprintf(`
        SELECT posts.id, posts.created_at, posts.title, posts.author, posts.link, posts.subreddit_id, posts.content,
		       posts.score, posts.promoted, posts.nsfw, posts.version, 
		       sub_post.id, sub_post.created_at, sub_post.title, sub_post.author, sub_post.link,
		       sub_post.content, sub_post.score, sub_post.promoted, sub_post.nsfw, sub_post.version
		FROM posts
		LEFT JOIN posts as sub_post ON sub_post.id = posts.subreddit_id
		WHERE posts.promoted = false
        ORDER BY posts.%s %s, posts.id ASC
        LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []interface{}{filters.Limit(), filters.Offset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, entitys.Metadata{}, err
	}

	defer rows.Close()

	posts := make([]*entitys.Post, 0, filters.Limit())

	for rows.Next() {
		var post entitys.Post
		var subPost entitys.SubPost

		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.Title,
			&post.Author,
			&post.Link,
			&post.SubredditID,
			&post.Content,
			&post.Score,
			&post.Promoted,
			&post.NSFW,
			&post.Version,

			&subPost.ID,
			&subPost.CreatedAt,
			&subPost.Title,
			&subPost.Author,
			&subPost.Link,
			&subPost.Content,
			&subPost.Score,
			&subPost.Promoted,
			&subPost.NSFW,
			&subPost.Version,
		)
		if err != nil {
			return nil, entitys.Metadata{}, err
		}

		if post.SubredditID != nil {
			post.Subreddit = &subPost
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, entitys.Metadata{}, err
	}

	if len(posts) >= 3 && len(promoted) >= 1 {

		if !posts[0].NSFW && !posts[1].NSFW {
			pr := []*entitys.Post{promoted[0]}

			posts = append(posts[:1], append(pr, posts[1:]...)...)
		}
	}

	if len(posts) > 16 && len(promoted) >= 2 {

		if !posts[14].NSFW && !posts[15].NSFW {
			pr := []*entitys.Post{promoted[1]}

			posts = append(posts[:15], append(pr, posts[15:]...)...)
		}
	}

	metadata := entitys.CalculateMetadata(len(posts), filters.Page, filters.PageSize)

	return posts, metadata, nil
}

func (p PostModel) getAllPromoted(ctx context.Context, filters entitys.Filters) ([]*entitys.Post, error) {
	query := `
        SELECT id, created_at, title, author, link, content, score, promoted, nsfw
		FROM posts
		WHERE promoted = true
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []interface{}{filters.PromotedPerPage(), filters.PromotedOffset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]*entitys.Post, 0, filters.PromotedPerPage())

	for rows.Next() {
		var post entitys.Post

		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.Title,
			&post.Author,
			&post.Link,
			&post.Content,
			&post.Score,
			&post.Promoted,
			&post.NSFW,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
