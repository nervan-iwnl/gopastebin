package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"

	firebase "firebase.google.com/go/v4"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type firebaseStore struct {
	bucket *storage.BucketHandle
	prefix string
}

func NewFirebaseStore(ctx context.Context, credPath, bucketName, prefix string) (FileStore, error) {
	if credPath == "" || bucketName == "" {
		return nil, fmt.Errorf("firebase creds or bucket not set")
	}

	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase init: %w", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase storage: %w", err)
	}

	bkt, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("bucket: %w", err)
	}

	return &firebaseStore{
		bucket: bkt,
		prefix: prefix,
	}, nil
}

func (f *firebaseStore) objectName(folder, name string) string {
	if folder != "" {
		return path.Join(f.prefix, folder, name+".txt")
	}
	return path.Join(f.prefix, name+".txt")
}

func (f *firebaseStore) Upload(ctx context.Context, folder, name, content string) (string, error) {
	obj := f.bucket.Object(f.objectName(folder, name))
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, bytes.NewBufferString(content)); err != nil {
		_ = w.Close()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return obj.ObjectName(), nil
}

func (f *firebaseStore) Download(ctx context.Context, folder, name string) (string, error) {
	obj := f.bucket.Object(f.objectName(folder, name))
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *firebaseStore) DownloadByPath(ctx context.Context, storagePath string) (string, error) {
	obj := f.bucket.Object(storagePath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *firebaseStore) DeleteByPath(ctx context.Context, storagePath string) error {
	obj := f.bucket.Object(storagePath)
	return obj.Delete(ctx)
}
