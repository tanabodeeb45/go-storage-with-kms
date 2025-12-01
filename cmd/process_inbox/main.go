package main

import (
	"context"
	"fmt"
	kms "go-storage-with-kms/internal/kms"
	stg "go-storage-with-kms/internal/storage"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func ProcessInbox(ctx context.Context, e stg.GCSEvent) error {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bucketName := os.Getenv("GCP_BUCKET_NAME")

	if bucketName == "" {
		fmt.Println("Warning: GCP_BUCKET_NAME env is not set")
	}

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
