package kms

import (
	"context"
	"fmt"
	"log"
	"os"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/joho/godotenv"
)

var (
	kmsClient  *kms.KeyManagementClient
	KMSKeyName string
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var KMSKeyName = os.Getenv("GCP_KMS_KEY")

	if KMSKeyName == "" {
		fmt.Println("Warning: GCP_KMS_KEY env is not set")
	}

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to create kms client: %w", err))
	}
	kmsClient = client
}

func Encrypt(ctx context.Context, plain []byte) ([]byte, error) {
	if KMSKeyName == "" {
		return nil, fmt.Errorf("GCP_KMS_KEY env is empty")
	}

	req := &kmspb.EncryptRequest{
		Name:      KMSKeyName,
		Plaintext: plain,
	}

	resp, err := kmsClient.Encrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("kms encrypt error: %w", err)
	}

	return resp.Ciphertext, nil
}
