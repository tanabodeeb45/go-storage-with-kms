# go-storage-with-kms

## Overview
This project demonstrates a tiny pipeline for Google Cloud Storage. Objects dropped in an `inbox/` folder are encrypted with Cloud KMS and re-uploaded to the same bucket under a `vault/` prefix with a `.enc` suffix. After a successful encryption cycle, the original inbox file is deleted so only the encrypted blob remains.

```
 inbox/myfile.txt  --read-->  KMS encrypt  --write-->  vault/myfile.txt.enc
                      \                               /
                       ---------- delete inbox -------
```

## Components
- `cmd/process_inbox` exposes `ProcessInbox`, intended for Cloud Functions / Cloud Run events. It reads the `GCSEvent`, filters for `inbox/` keys, orchestrates encryption, writes to `vault/`, deletes the source, and logs the operation.
- `internal/kms` creates a global KMS client (using `GCP_KMS_KEY`) and exposes `Encrypt(ctx, data)` which wraps `cloud.google.com/go/kms`.
- `internal/storage` provides helpers for the inbox (`ReadInboxFile`, `DeleteInboxFile`) and vault (`WriteVaultFile`) object operations.

## Environment
| Variable | Description |
|----------|-------------|
| `GCP_BUCKET_NAME` | Bucket hosting the `inbox/` and `vault/` objects. |
| `GCP_KMS_KEY` | Full resource ID of the Cloud KMS key used for encryption. Required at runtime. |

The Google Cloud credentials must allow access to Storage (read/write/delete on the bucket) and to use the specified KMS key.

## How it works
1. A Cloud Storage notification invokes `ProcessInbox` with the object metadata.
2. If the object name does not start with `inbox/`, the function no-ops.
3. Inbox files are downloaded, encrypted via Cloud KMS, and re-uploaded to `vault/<name>.enc`.
4. The original inbox object is deleted to keep only the encrypted version.
5. Success is logged via the standard logger.

## Development
- Requires Go 1.24.
- Dependencies are managed through `go.mod`; run `go mod tidy` if you add or remove imports.
- For local testing, ensure `GCP_BUCKET_NAME` and `GCP_KMS_KEY` are set and that your gcloud identity has the necessary IAM permissions.
