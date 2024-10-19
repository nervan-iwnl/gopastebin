package fb

import (
	"context"
	"fmt"
	"io"
)

func UploadFileToFirebase(userID string, path string, fileContent string) (string, error) {
	ctx := context.Background()

	if storageClient == nil {
		return "", fmt.Errorf("клиент Firebase Storage не инициализирован")
	}

	storagePath := fmt.Sprintf("%s/%s.txt", userID, path)

	obj := bucket.Object(storagePath)
	writer := obj.NewWriter(ctx)

	if _, err := io.WriteString(writer, fileContent); err != nil {
		return "", fmt.Errorf("ошибка записи файла: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("ошибка закрытия writer: %v", err)
	}

	return storagePath, nil
}
