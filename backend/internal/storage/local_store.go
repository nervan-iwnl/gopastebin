package storage

import (
	"context"
	"os"
	"path/filepath"
)

type localStore struct {
	base string
}

func NewLocalStore(base string) FileStore {
	_ = os.MkdirAll(base, 0755)
	return &localStore{base: base}
}

func (l *localStore) Upload(ctx context.Context, folder, name, content string) (string, error) {
	dir := filepath.Join(l.base, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, name+".txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return filepath.Join(folder, name+".txt"), nil
}

func (l *localStore) Download(ctx context.Context, folder, name string) (string, error) {
	path := filepath.Join(l.base, folder, name+".txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (l *localStore) DownloadByPath(ctx context.Context, storagePath string) (string, error) {
	path := filepath.Join(l.base, storagePath)
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (l *localStore) DeleteByPath(ctx context.Context, storagePath string) error {
	path := filepath.Join(l.base, storagePath)
	return os.Remove(path)
}
