package service

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"gopastebin/internal/domain"
	"gopastebin/internal/repository"
	"gopastebin/internal/storage"
	"gopastebin/pkg/slug"
)

type PasteService struct {
	repo repository.PasteRepository
	fs   storage.FileStore
}

type CreatePasteDTO struct {
	Title     string `json:"title" binding:"required,max=200"`
	Content   string `json:"content" binding:"required"`
	Extension string `json:"extension"`
	Folder    string `json:"folder"`
	IsPublic  *bool  `json:"is_public"`
	TTLMin    int    `json:"ttl_minutes"`
	Slug      string `json:"slug"` // кастомный слаг от юзера
}

type UpdatePasteDTO struct {
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	Extension *string `json:"extension"`
	Folder    *string `json:"folder"`
	IsPublic  *bool   `json:"is_public"`
	TTLMin    *int    `json:"ttl_minutes"`
}

func NewPasteService(r repository.PasteRepository, fs storage.FileStore) *PasteService {
	return &PasteService{repo: r, fs: fs}
}

func cleanFolder(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")
	s = path.Clean(s)
	if s == "." {
		return ""
	}
	return s
}

// ---------- CREATE USER PASTE с перезаписью по тому же slug ----------
func (s *PasteService) CreateUserPaste(ctx context.Context, user any, dto CreatePasteDTO) (*domain.Paste, error) {
	userID := domain.ExtractUserID(user)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	logicalFolder := cleanFolder(dto.Folder)
	fileFolder := fmt.Sprintf("%d", userID)
	if logicalFolder != "" {
		fileFolder = path.Join(fileFolder, logicalFolder)
	}

	// 1. определяем, какой слаг хотим
	slugVal := dto.Slug
	if slugVal == "" {
		var err error
		slugVal, err = s.generateUniqueSlug(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 2. пробуем найти пасту с таким слагом
	existing, err := s.repo.FindBySlug(ctx, slugVal)
	if err == nil {
		// паста есть
		if existing.UserID != userID {
			// чужая паста с таким же слагом
			return nil, errors.New("slug_taken")
		}

		// это наша старая паста → перезаписываем файл и обновляем метаданные
		storagePath, err := s.fs.Upload(ctx, fileFolder, slugVal, dto.Content)
		if err != nil {
			return nil, fmt.Errorf("upload_failed: %w", err)
		}

		existing.Title = dto.Title
		existing.StoragePath = storagePath
		existing.Extension = dto.Extension
		existing.Folder = logicalFolder
		existing.IsPublic = true
		if dto.IsPublic != nil {
			existing.IsPublic = *dto.IsPublic
		}
		// если раньше мы её "мягко удалили" (по сути — истекла), то сейчас оживим
		existing.ExpiresAt = nil
		if dto.TTLMin > 0 {
			ex := time.Now().Add(time.Duration(dto.TTLMin) * time.Minute)
			existing.ExpiresAt = &ex
		}

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// 3. пасты с таким слагом нет → создаём новую
	storagePath, err := s.fs.Upload(ctx, fileFolder, slugVal, dto.Content)
	if err != nil {
		return nil, fmt.Errorf("upload_failed: %w", err)
	}

	p := &domain.Paste{
		Title:       dto.Title,
		Slug:        slugVal,
		StoragePath: storagePath,
		Extension:   dto.Extension,
		UserID:      userID,
		Folder:      logicalFolder,
		IsPublic:    true,
	}

	if dto.IsPublic != nil {
		p.IsPublic = *dto.IsPublic
	}
	if dto.TTLMin > 0 {
		ex := time.Now().Add(time.Duration(dto.TTLMin) * time.Minute)
		p.ExpiresAt = &ex
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ---------- CREATE ANON PASTE с перезаписью по тому же slug ----------
func (s *PasteService) CreateAnonPaste(ctx context.Context, dto CreatePasteDTO) (*domain.Paste, error) {
	logicalFolder := cleanFolder(dto.Folder)
	fileFolder := "anon"
	if logicalFolder != "" {
		fileFolder = path.Join("anon", logicalFolder)
	}

	slugVal := dto.Slug
	if slugVal == "" {
		var err error
		slugVal, err = s.generateUniqueSlug(ctx)
		if err != nil {
			return nil, err
		}
	}

	// пытаемся найти такую же анонимную пасту
	existing, err := s.repo.FindBySlug(ctx, slugVal)
	if err == nil {
		// паста есть; если она анонимная (user_id == 0), перезапишем
		if existing.UserID != 0 {
			return nil, errors.New("slug_taken")
		}

		storagePath, err := s.fs.Upload(ctx, fileFolder, slugVal, dto.Content)
		if err != nil {
			return nil, err
		}

		existing.Title = dto.Title
		existing.StoragePath = storagePath
		existing.Extension = dto.Extension
		existing.Folder = logicalFolder
		existing.IsPublic = true
		if dto.IsPublic != nil {
			existing.IsPublic = *dto.IsPublic
		}
		existing.ExpiresAt = nil
		if dto.TTLMin > 0 {
			ex := time.Now().Add(time.Duration(dto.TTLMin) * time.Minute)
			existing.ExpiresAt = &ex
		}

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// не нашли — создаём
	storagePath, err := s.fs.Upload(ctx, fileFolder, slugVal, dto.Content)
	if err != nil {
		return nil, err
	}

	p := &domain.Paste{
		Title:       dto.Title,
		Slug:        slugVal,
		StoragePath: storagePath,
		Extension:   dto.Extension,
		UserID:      0,
		Folder:      logicalFolder,
		IsPublic:    true,
	}

	if dto.IsPublic != nil {
		p.IsPublic = *dto.IsPublic
	}
	if dto.TTLMin > 0 {
		ex := time.Now().Add(time.Duration(dto.TTLMin) * time.Minute)
		p.ExpiresAt = &ex
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ---------- GET + CONTENT ----------
func (s *PasteService) GetPasteWithContent(ctx context.Context, slug string) (*domain.Paste, string, error) {
	p, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, "", err
	}
	// "мягко удалённые" мы помечаем через ExpiresAt — тут как раз отсекаем
	if p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt) {
		return nil, "", errors.New("not_found")
	}

	content, err := s.fs.DownloadByPath(ctx, p.StoragePath)
	if err != nil {
		return nil, "", err
	}

	return p, content, nil
}

// ---------- SOFT DELETE ----------
func (s *PasteService) DeletePaste(ctx context.Context, user any, slugStr string) error {
	userID := domain.ExtractUserID(user)
	p, err := s.repo.FindBySlug(ctx, slugStr)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return errors.New("forbidden")
	}

	// файл можно удалить сразу, чтобы место не держать
	if p.StoragePath != "" {
		_ = s.fs.DeleteByPath(ctx, p.StoragePath)
	}

	// "мягко удаляем": ставим истечение на сейчас
	now := time.Now()
	p.ExpiresAt = &now

	return s.repo.Update(ctx, p)
}

// ---------- LISTS ----------
func (s *PasteService) GetMyPastes(ctx context.Context, user any) ([]domain.Paste, error) {
	uid := domain.ExtractUserID(user)
	return s.repo.FindByUser(ctx, uid)
}

func (s *PasteService) GetMyPastesInFolder(ctx context.Context, user any, folder string) ([]domain.Paste, error) {
	uid := domain.ExtractUserID(user)
	return s.repo.FindByUserAndFolder(ctx, uid, cleanFolder(folder))
}

func (s *PasteService) GetMyFolders(ctx context.Context, user any) ([]string, error) {
	uid := domain.ExtractUserID(user)
	return s.repo.DistinctFolders(ctx, uid)
}

func (s *PasteService) GetRecentPublic(ctx context.Context, limit int) ([]domain.Paste, error) {
	return s.repo.RecentPublic(ctx, limit)
}

// ---------- UPDATE ----------
func (s *PasteService) UpdatePaste(ctx context.Context, user any, slugStr string, dto UpdatePasteDTO) (*domain.Paste, error) {
	uid := domain.ExtractUserID(user)
	p, err := s.repo.FindBySlug(ctx, slugStr)
	if err != nil {
		return nil, err
	}
	if p.UserID != uid {
		return nil, errors.New("forbidden")
	}

	if dto.Title != nil {
		p.Title = *dto.Title
	}
	if dto.Extension != nil {
		p.Extension = *dto.Extension
	}
	if dto.IsPublic != nil {
		p.IsPublic = *dto.IsPublic
	}
	if dto.TTLMin != nil {
		if *dto.TTLMin <= 0 {
			p.ExpiresAt = nil
		} else {
			ex := time.Now().Add(time.Duration(*dto.TTLMin) * time.Minute)
			p.ExpiresAt = &ex
		}
	}
	if dto.Folder != nil {
		p.Folder = cleanFolder(*dto.Folder)
	}

	// если контент поменяли — перезаливаем
	if dto.Content != nil {
		fileFolder := "anon"
		if p.UserID != 0 {
			fileFolder = fmt.Sprintf("%d", p.UserID)
		}
		if p.Folder != "" {
			fileFolder = path.Join(fileFolder, p.Folder)
		}
		storagePath, err := s.fs.Upload(ctx, fileFolder, p.Slug, *dto.Content)
		if err != nil {
			return nil, err
		}
		p.StoragePath = storagePath
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ---------- SLUG GENERATION ----------
func (s *PasteService) generateUniqueSlug(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ {
		candidate := slug.Generate()
		free, err := s.repo.IsSlugFree(ctx, candidate)
		if err != nil {
			return "", err
		}
		if free {
			return candidate, nil
		}
	}
	return "", errors.New("cannot_generate_slug")
}

func (s *PasteService) GetMyPastesPaged(ctx context.Context, user any, limit, offset int) ([]domain.Paste, error) {
	uid := domain.ExtractUserID(user)
	return s.repo.FindByUserPaged(ctx, uid, limit, offset)
}

func (s *PasteService) GetRecentPublicPaged(ctx context.Context, limit, offset int) ([]domain.Paste, error) {
	return s.repo.RecentPublicPaged(ctx, limit, offset)
}
