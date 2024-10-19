package fb

import (
	"context"
	"fmt"
	"io"
)

func GetFileFromFirebase(userID string, path string) (string, error) {
	ctx := context.Background()

	if storageClient == nil {
		return "", fmt.Errorf("firebase Storage client has not been initialized")
	}

	storagePath := fmt.Sprintf("%s/%s.txt", userID, path)

	obj := bucket.Object(storagePath)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating a reader for a file: %v", err)
	}
	defer reader.Close() 

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("file reading error: %v", err)
	}


	return string(data), nil
}
