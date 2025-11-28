package storage

type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}
