package service

import (
	"context"

	"gopastebin/internal/storage"
)

type DynamicFileStore struct {
	settings *AppSettingsService
	local    storage.FileStore
	firebase storage.FileStore
}

func NewDynamicFileStore(
	settings *AppSettingsService,
	local storage.FileStore,
	firebase storage.FileStore,
) *DynamicFileStore {
	return &DynamicFileStore{
		settings: settings,
		local:    local,
		firebase: firebase,
	}
}

func (d *DynamicFileStore) pick() storage.FileStore {
	cur := d.settings.GetStorageProvider()
	if cur == "firebase" && d.firebase != nil {
		return d.firebase
	}
	return d.local
}

func (d *DynamicFileStore) Upload(ctx context.Context, folder, name, content string) (string, error) {
	return d.pick().Upload(ctx, folder, name, content)
}

func (d *DynamicFileStore) Download(ctx context.Context, folder, name string) (string, error) {
	return d.pick().Download(ctx, folder, name)
}

func (d *DynamicFileStore) DownloadByPath(ctx context.Context, p string) (string, error) {
	return d.pick().DownloadByPath(ctx, p)
}

func (d *DynamicFileStore) DeleteByPath(ctx context.Context, p string) error {
	return d.pick().DeleteByPath(ctx, p)
}
