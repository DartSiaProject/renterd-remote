package models

type BucketObject struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	FileName string `json:"filename"`
}
