package fb

import (
	"context"
	"fmt"
	"io"
)

func GetFileFromFirebase(userID string, path string) (string, error) {
	ctx := context.Background()

	if storageClient == nil {
		return "", fmt.Errorf("клиент Firebase Storage не инициализирован")
	}

	storagePath := fmt.Sprintf("%s/%s.txt", userID, path)

	obj := bucket.Object(storagePath)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("ошибка создания reader для файла: %v", err)
	}
	defer reader.Close() 

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %v", err)
	}


	return string(data), nil
}
