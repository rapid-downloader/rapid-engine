package client

import "time"

type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Path     string    `json:"path"`
	Domain   string    `json:"domain"`
	Expires  time.Time `json:"expirationDate"`
	Secure   bool      `json:"secure"`
	HttpOnly bool      `json:"httpOnly"`
	SameSite string    `json:"sameSite"`
}

type Request struct {
	Url       string    `json:"url"`
	Provider  string    `json:"provider"`
	Client    *string   `json:"client"`
	MimeType  *string   `json:"mimeType"`
	UserAgent *string   `json:"userAgent"`
	Cookies   *[]Cookie `json:"cookies"`
}

type Download struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Url              string    `json:"url"`
	Provider         string    `json:"provider"`
	Size             int64     `json:"size"`
	Type             string    `json:"type"`
	Chunklen         int       `json:"chunklen"`
	Resumable        bool      `json:"resumable"`
	Progress         float64   `json:"progress"`
	Expired          bool      `json:"expired"`
	DownloadedChunks []int64   `json:"downloadedChunks"`
	TimeLeft         int       `json:"timeLeft"`
	Speed            int       `json:"speed"`
	Status           string    `json:"status"`
	Date             time.Time `json:"date"`
}

type Progress struct {
	ID     string          `json:"id"`
	Done   bool            `json:"done"`
	Chunks []ChunkProgress `json:"chunks"`
}

type ChunkProgress struct {
	Downloaded int64   `json:"downloaded"`
	Size       int64   `json:"size"`
	Progress   float64 `json:"progress"`
	Done       bool    `json:"done"`
}

type OnProgress = func(progress Progress, err error)
