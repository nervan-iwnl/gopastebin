package fb

import (
	"context"
	"fmt"
	"io"
)

func UploadFileToFirebase(userID string, path string, fileContent string) (string, error) {
	ctx := context.Background()

	if storageClient == nil {
		return "", fmt.Errorf("the Firebase Storage client is not initialized")
	}

	storagePath := fmt.Sprintf("%s/%s.txt", userID, path)

	obj := bucket.Object(storagePath)
	writer := obj.NewWriter(ctx)

	if _, err := io.WriteString(writer, fileContent); err != nil {
		return "", fmt.Errorf("file recording error: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	return storagePath, nil
}
