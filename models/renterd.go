package models

type BucketObject struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	FileName string `json:"filename"`
}

type BucketLargeObject struct {
	Bucket         string `json:"bucket"`
	Key            string `json:"key"`
	FileName       string `json:"filename"`
	FilePart       int    `json:"filepart"`
	FileSize       int    `json:"filesize"`
	RemainingParts int    `json:"remainingparts"`
}
