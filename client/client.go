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
	Url         string    `json:"url"`
	Provider    string    `json:"provider"`
	Client      *string   `json:"client"`
	ContentType *string   `json:"contentType"`
	UserAgent   *string   `json:"userAgent"`
	Cookies     *[]Cookie `json:"cookies"`
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
	ID         string  `json:"id"`
	Index      int     `json:"index"`
	Downloaded int64   `json:"downloaded"`
	Size       int64   `json:"size"`
	Progress   float64 `json:"progress"`
	Done       bool    `json:"done"`
}

type OnProgress = func(progress Progress, err error)

type Rapid interface {
	// will return the download progress for each chunk.
	// if error happens during the way, the progress will contains zero value for each properties
	// and the caller should return the error and end the callback
	Listen(progress OnProgress)

	Fetch(req Request) (*Download, error)
	Download(id string) error
	Resume(id string) error
	Restart(id string) error
	Stop(id string) error
	Pause(id string) error
}

type RapidCloser interface {
	Close() error
}
