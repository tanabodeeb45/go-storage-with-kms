package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

func ReadInboxFile(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage client error: %w", err)
	}
	defer client.Close()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open inbox file '%s': %w", objectName, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("cannot read inbox file '%s': %w", objectName, err)
	}

	return data, nil
}

func DeleteInboxFile(ctx context.Context, bucketName, objectName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage client error: %w", err)
	}
	defer client.Close()

	return client.Bucket(bucketName).Object(objectName).Delete(ctx)
}
