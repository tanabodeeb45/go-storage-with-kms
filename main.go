package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	kms "go-storage-with-kms/internal/kms"
	stg "go-storage-with-kms/internal/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("ProcessInbox", ProcessInbox)
}

type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func ProcessInbox(ctx context.Context, e event.Event) error {
	var data GCSEvent
	if err := e.DataAs(&data); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	bucketName := data.Bucket
	inboxKey := data.Name

	if !strings.HasPrefix(inboxKey, "inbox/") {
		log.Println("Not inbox file => skip")
		return nil
	}

	fileData, err := stg.ReadInboxFile(ctx, bucketName, inboxKey)
	if err != nil {
		return fmt.Errorf("read inbox: %w", err)
	}

	encrypted, err := kms.Encrypt(ctx, fileData)
	if err != nil {
		return fmt.Errorf("encrypt error: %w", err)
	}

	vaultKey := strings.Replace(inboxKey, "inbox/", "vault/", 1) + ".enc"

	err = stg.WriteVaultFile(ctx, bucketName, vaultKey, encrypted)
	if err != nil {
		return fmt.Errorf("write vault: %w", err)
	}

	err = stg.DeleteInboxFile(ctx, bucketName, inboxKey)
	if err != nil {
		return fmt.Errorf("delete inbox: %w", err)
	}

	log.Printf("Encrypted %s → %s\n", inboxKey, vaultKey)
	return nil
}
