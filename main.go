package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
)

type (
	DownloadItem struct {
		ByExtensionId    string `json:"byExtensionId"`
		ByExtensionName  string `json:"byExtensionName"`
		BytesReceived    int64  `json:"bytesReceived"`
		CanResume        bool   `json:"canResume"`
		CookieStoreId    string `json:"cookieStoreId"`
		Danger           string `json:"danger"`
		EndTime          string `json:"endTime"`
		Error            string `json:"error"`
		EstimatedEndTime string `json:"estimatedEndTime"`
		Exists           bool   `json:"exists"`
		Filename         string `json:"filename"`
		FileSize         int64  `json:"fileSize"`
		ID               int    `json:"id"`
		Incognito        bool   `json:"incognito"`
		Mime             string `json:"mime"`
		Paused           bool   `json:"paused"`
		Referrer         string `json:"referrer"`
		StartTime        string `json:"startTime"`
		State            string `json:"state"`
		TotalBytes       int64  `json:"totalBytes"`
		Url              string `json:"url"`
	}

	Cookie struct {
		Name  string
		Value string

		Path       string    // optional
		Domain     string    // optional
		Expires    time.Time // optional
		RawExpires string    // for reading cookies only

		// MaxAge=0 means no 'Max-Age' attribute specified.
		// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
		// MaxAge>0 means Max-Age attribute present and given in seconds
		MaxAge   int
		Secure   bool
		HttpOnly bool
		SameSite string
		Raw      string
		Unparsed []string // Raw text of unparsed attribute-value pairs
	}

	Body struct {
		Item      DownloadItem `json:"item"`
		Cookies   []Cookie     `json:"cookies"`
		UserAgent string       `json:"userAgent"`
	}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer runtime.GC()

		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			log.Println("Error parsing body:", err)
		}

		var body Body
		if err := json.Unmarshal(buffer.Bytes(), &body); err != nil {
			log.Println("Error unmarshalling body:", err)
		}

		cookies := make([]*http.Cookie, len(body.Cookies))
		for i, cookie := range body.Cookies {
			cookies[i] = &http.Cookie{
				Name:       cookie.Name,
				Value:      cookie.Value,
				Path:       cookie.Path,
				Domain:     cookie.Domain,
				Expires:    cookie.Expires,
				RawExpires: cookie.RawExpires,
				MaxAge:     cookie.MaxAge,
				Secure:     cookie.Secure,
				HttpOnly:   cookie.HttpOnly,
				SameSite:   0,
				Raw:        cookie.Raw,
				Unparsed:   cookie.Unparsed,
			}
		}

		entry, err := entry.Fetch(
			body.Item.Url,
			entry.AddCookies(cookies),
			entry.AddHeaders(entry.Headers{
				"User-Agent":   body.UserAgent,
				"Content-Type": body.Item.Mime,
			}),
		)

		if err != nil {
			log.Println("Error fetching url:", err)
			return
		}

		dl := downloader.New(downloader.Default)
		if watcher, ok := dl.(downloader.Watcher); ok {
			watcher.Watch(func(data ...interface{}) {
				log.Println(data...)
			})
		}

		if err := dl.Download(entry); err != nil {
			log.Printf("Error downloading %s: %v", entry.Name(), err.Error())
		}
	})

	fmt.Println("Listening for browser extension...")
	http.ListenAndServe(":6969", nil)
}
