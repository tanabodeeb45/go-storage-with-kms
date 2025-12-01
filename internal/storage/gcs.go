package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

func ReadInboxFile(ctx context.Context, bucketName, object string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	rc, err := client.Bucket(bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

func WriteVaultFile(ctx context.Context, bucketName, object string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w := client.Bucket(bucketName).Object(object).NewWriter(ctx)
	defer w.Close()

	_, err = w.Write(data)
	return err
}

func DeleteInboxFile(ctx context.Context, bucketName, object string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Bucket(bucketName).Object(object).Delete(ctx)
}
