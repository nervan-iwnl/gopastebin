package storage

import "context"

type FileStore interface {
	Upload(ctx context.Context, folder, name, content string) (string, error)
	Download(ctx context.Context, folder, name string) (string, error)
	DownloadByPath(ctx context.Context, storagePath string) (string, error)
	DeleteByPath(ctx context.Context, storagePath string) error
}
