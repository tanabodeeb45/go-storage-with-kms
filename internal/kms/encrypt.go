package kms

import (
	"context"
	"os"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
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

	keyName := os.Getenv("GCP_KMS_KEY")

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
