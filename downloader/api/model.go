package api

type ProgressBar struct {
	ID         string `json:"id"`
	Index      int    `json:"index"`
	Downloaded int64  `json:"downloaded"`
	Progress   int64  `json:"progress"`
	Size       int64  `json:"size"`
	Done       bool   `json:"done"`
}
