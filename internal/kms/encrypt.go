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

var kmsClient *kms.KeyManagementClient

func getClient(ctx context.Context) (*kms.KeyManagementClient, error) {
	if kmsClient != nil {
		return kmsClient, nil
	}

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	kmsClient = client
	return kmsClient, nil
}

func Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	client, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	env := godotenv.Load()
	if env != nil {
		return nil, env
	}

	keyName := os.Getenv("GCP_KMS_KEY")

	fmt.Println(keyName)

	if keyName == "" {
		log.Fatal("GCP_KMS_KEY environment variable not set")
	}

	req := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: plaintext,
	}

	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Ciphertext, nil
}
