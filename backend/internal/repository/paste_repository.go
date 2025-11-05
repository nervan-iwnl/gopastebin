package repository

import (
	"context"
	"time"

	"gopastebin/internal/domain"

	"gorm.io/gorm"
)

type PasteRepository interface {
	Create(ctx context.Context, p *domain.Paste) error
	Update(ctx context.Context, p *domain.Paste) error
	Delete(ctx context.Context, p *domain.Paste) error
	FindBySlug(ctx context.Context, slug string) (*domain.Paste, error)
	IsSlugFree(ctx context.Context, slug string) (bool, error)
	FindByUser(ctx context.Context, userID uint) ([]domain.Paste, error)
	FindByUserAndFolder(ctx context.Context, userID uint, folder string) ([]domain.Paste, error)
	DistinctFolders(ctx context.Context, userID uint) ([]string, error)
	RecentPublic(ctx context.Context, limit int) ([]domain.Paste, error)
	PublicOfUser(ctx context.Context, username string, limit int) ([]domain.Paste, error)
	RecentPublicPaged(ctx context.Context, limit, offset int) ([]domain.Paste, error)
	FindByUserPaged(ctx context.Context, userID uint, limit, offset int) ([]domain.Paste, error)
}

type pasteRepo struct {
	db *gorm.DB
}

func NewPasteRepository(db *gorm.DB) PasteRepository {
	return &pasteRepo{db: db}
}

func (r *pasteRepo) Create(ctx context.Context, p *domain.Paste) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pasteRepo) Update(ctx context.Context, p *domain.Paste) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *pasteRepo) Delete(ctx context.Context, p *domain.Paste) error {
	return r.db.WithContext(ctx).Delete(p).Error
}

func (r *pasteRepo) FindBySlug(ctx context.Context, slug string) (*domain.Paste, error) {
	var p domain.Paste
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pasteRepo) IsSlugFree(ctx context.Context, slug string) (bool, error) {
	var c int64
	err := r.db.WithContext(ctx).Model(&domain.Paste{}).Where("slug = ?", slug).Count(&c).Error
	return c == 0, err
}

func (r *pasteRepo) FindByUser(ctx context.Context, userID uint) ([]domain.Paste, error) {
	var list []domain.Paste
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *pasteRepo) FindByUserAndFolder(ctx context.Context, userID uint, folder string) ([]domain.Paste, error) {
	var list []domain.Paste
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if folder == "" {
		q = q.Where("folder = '' OR folder IS NULL")
	} else {
		q = q.Where("folder = ?", folder)
	}
	err := q.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *pasteRepo) DistinctFolders(ctx context.Context, userID uint) ([]string, error) {
	var folders []string
	err := r.db.WithContext(ctx).
		Model(&domain.Paste{}).
		Where("user_id = ?", userID).
		Distinct().
		Pluck("folder", &folders).Error
	return folders, err
}

func (r *pasteRepo) RecentPublic(ctx context.Context, limit int) ([]domain.Paste, error) {
	var list []domain.Paste
	err := r.db.WithContext(ctx).
		Where("is_public = ? AND (expires_at IS NULL OR expires_at > ?)", true, time.Now()).
		Order("created_at DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *pasteRepo) PublicOfUser(ctx context.Context, username string, limit int) ([]domain.Paste, error) {
	var list []domain.Paste
	err := r.db.WithContext(ctx).
		Table("pastes").
		Joins("JOIN users ON pastes.user_id = users.id").
		Where("users.username = ? AND pastes.is_public = ? AND (pastes.expires_at IS NULL OR pastes.expires_at > ?)",
			username, true, time.Now()).
		Order("pastes.created_at DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// ↓↓↓ вот тут были ошибки у тебя — ставим *pasteRepo
func (r *pasteRepo) RecentPublicPaged(ctx context.Context, limit, offset int) ([]domain.Paste, error) {
	var res []domain.Paste
	q := r.db.WithContext(ctx).
		Where("is_public = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)
	if err := q.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *pasteRepo) FindByUserPaged(ctx context.Context, userID uint, limit, offset int) ([]domain.Paste, error) {
	var res []domain.Paste
	q := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)
	if err := q.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

