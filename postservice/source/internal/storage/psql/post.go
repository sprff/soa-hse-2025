package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"postservice/internal/utils"
	"social/shared/models"
	"time"

	"github.com/lib/pq"
)

type psqlPost struct {
	ID          models.PostID  `json:"id" db:"id"`
	Title       string         `json:"title" db:"title"`
	Description string         `json:"description" db:"description"`
	CreatorID   string         `json:"creator_id" db:"creator_id"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	IsPrivate   bool           `json:"is_private" db:"is_private"`
	Tags        pq.StringArray `json:"tags" db:"tags"`
}

func convertToPsql(p models.Post) psqlPost {
	return psqlPost{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		CreatorID:   p.CreatorID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		IsPrivate:   p.IsPrivate,
		Tags:        p.Tags,
	}
}
func convertToModels(p psqlPost) models.Post {
	return models.Post{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		CreatorID:   p.CreatorID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		IsPrivate:   p.IsPrivate,
		Tags:        p.Tags,
	}
}

// CreatePost implements storage.CommonRepository.
func (p *PsqlStorage) CreatePost(ctx context.Context, postArg models.Post) (models.PostID, error) {
	slog.DebugContext(ctx, "PSQL CreatePost")
	post := convertToPsql(postArg)
	post.ID = models.PostID(utils.GenerateUUID())
	slog.DebugContext(ctx, "Gen ID", "id", post.ID)

	_, err := p.db.NamedExec(`
	INSERT INTO posts
	(
	id, title, description, creator_id,
	is_private, tags, created_at, updated_at
	)
	VALUES
	(:id, :title, :description, :creator_id,
	:is_private, :tags, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, &post)
	if err != nil {
		return models.PostID(""), fmt.Errorf("can't insert post: %w", err)
	}

	return post.ID, nil
}

// GetPostByID implements storage.CommonRepository.
func (p *PsqlStorage) GetPostByID(ctx context.Context, id models.PostID) (models.Post, error) {
	slog.DebugContext(ctx, "PSQL GetPostByID")

	var post psqlPost
	err := p.db.Get(&post, "SELECT * FROM posts WHERE id=$1", id)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return models.Post{}, models.ErrPostNotFound
		default:
			return models.Post{}, fmt.Errorf("can't select: %w", err)
		}
	}
	return convertToModels(post), nil
}

// GetPosts implements storage.CommonRepository.
func (p *PsqlStorage) GetPosts(ctx context.Context, offset int, limit int) ([]models.Post, error) {
	slog.DebugContext(ctx, "PSQL GetPostByID")

	var posts []psqlPost
	err := p.db.Select(&posts, "SELECT * FROM posts offset $1 limit $2", offset, limit)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return []models.Post{}, models.ErrPostNotFound
		default:
			return nil, fmt.Errorf("can't select: %w", err)
		}
	}
	if len(posts) == 0 {
		return []models.Post{}, models.ErrPostNotFound
	}

	cposts := make([]models.Post, len(posts))
	for i := range posts {
		cposts[i] = convertToModels(posts[i])
	}
	return cposts, nil
}

// UpdatePost implements storage.CommonRepository.
func (p *PsqlStorage) UpdatePost(ctx context.Context, postArg models.Post) error {
	slog.DebugContext(ctx, "PSQL UpdatePost")
	post := convertToPsql(postArg)
	res, err := p.db.NamedExec(`
	UPDATE posts
	SET
	id=:id,title=:title, description=:description,
	is_private=:is_private, tags=:tags, updated_at=CURRENT_TIMESTAMP
	WHERE id=:id`, &post)
	if err != nil {
		return fmt.Errorf("can't insert user: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected")
	}
	if affected == 0 {
		return models.ErrPostNotFound
	}
	return nil
}

// UpdatePost implements storage.CommonRepository.
func (p *PsqlStorage) DeletePost(ctx context.Context, id models.PostID) error {
	slog.DebugContext(ctx, "PSQL DeletePost")
	res, err := p.db.Exec(`DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("can't insert user: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected")
	}
	if affected == 0 {
		return models.ErrPostNotFound
	}
	return nil
}
