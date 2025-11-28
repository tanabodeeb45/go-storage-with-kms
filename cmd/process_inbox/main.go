package main

import (
	"context"
	kms "go-storage-with-kms/internal/kms"
	stg "go-storage-with-kms/internal/storage"
	"log"
	"os"
	"strings"
)

var bucketName = os.Getenv("GCP_BUCKET_NAME")

func ProcessInbox(ctx context.Context, e stg.GCSEvent) error {
	inboxKey := e.Name

	if !strings.HasPrefix(inboxKey, "inbox/") {
		return nil
	}

	fileData, err := stg.ReadInboxFile(ctx, bucketName, inboxKey)
	if err != nil {
		return err
	}

	encrypted, err := kms.Encrypt(ctx, fileData)
	if err != nil {
		return err
	}

	vaultKey := strings.Replace(inboxKey, "inbox/", "vault/", 1) + ".enc"

	err = stg.WriteVaultFile(ctx, bucketName, vaultKey, encrypted)
	if err != nil {
		return err
	}

	err = stg.DeleteInboxFile(ctx, bucketName, inboxKey)
	if err != nil {
		return err
	}

	log.Println("Encrypted:", inboxKey)
	log.Println("Saved vault:", vaultKey)

	return nil
}

func main() {}
