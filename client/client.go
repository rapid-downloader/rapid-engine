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
	ID               string
	Name             string
	Url              string
	Provider         string
	Size             int64
	Type             string
	Chunklen         int
	Resumable        bool
	Progress         float64
	Expired          bool
	DownloadedChunks []int64
	TimeLeft         int
	Speed            int
	Status           string
	Date             time.Time
}

type Progress struct {
	ID         string
	Index      int
	Downloaded int64
	Size       int64
	Progress   float64
	Done       bool
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
