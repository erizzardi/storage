package util

type Row struct {
	Uuid     string `json:"uuid"`
	FileName string `json:"name"`
	Bucket   string `json:"bucket,omitempty"`
}
