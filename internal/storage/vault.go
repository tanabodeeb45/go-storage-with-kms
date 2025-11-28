package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func WriteVaultFile(ctx context.Context, bucketName string, vaultKey string, encrypted []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w := client.Bucket(bucketName).Object(vaultKey).NewWriter(ctx)
	defer w.Close()

	_, err = w.Write(encrypted)
	if err != nil {
		return fmt.Errorf("cannot write vault file: %w", err)
	}

	return nil
}
