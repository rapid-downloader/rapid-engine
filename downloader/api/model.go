package api

type ProgressBar struct {
	ID         string `json:"id"`
	Index      int    `json:"index"`
	Downloaded int64  `json:"downloaded"`
	Size       int64  `json:"size"`
	Progress   int    `json:"progress"`
	Done       bool   `json:"done"`
}
