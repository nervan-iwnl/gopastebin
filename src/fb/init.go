package fb

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var storageClient *storage.Client
var bucket *storage.BucketHandle

func InitFirebaseApp() error {
	ctx := context.Background()

	opt := option.WithCredentialsFile("firebase_settings.json")

	var err error
	storageClient, err = storage.NewClient(ctx, opt)
	if err != nil {
		return fmt.Errorf("firebase Storage client initialization error: %v", err)
	}

	bucketName := os.Getenv("FIREBASE_BUCKET_NAME")
	if bucketName == "" {
		return fmt.Errorf("bucket name is not specified in the .env file")
	}

	bucket = storageClient.Bucket(bucketName)
	if bucket == nil {
		return fmt.Errorf("error receiving the bucket %s", bucketName)
	}

	return nil
}
