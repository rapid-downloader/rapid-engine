package api

import (
	"net/http"
	"time"

	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	cookie struct {
		Name     string    `json:"name"`
		Value    string    `json:"value"`
		Path     string    `json:"path"`
		Domain   string    `json:"domain"`
		Expires  time.Time `json:"expirationDate"`
		Secure   bool      `json:"secure"`
		HttpOnly bool      `json:"httpOnly"`
		SameSite int       `json:"sameSite"`
	}

	request struct {
		Url         string   `json:"url"`
		Client      string   `json:"client"`
		Provider    string   `json:"provider"`
		ContentType string   `json:"contentType"`
		UserAgent   string   `json:"userAgent"`
		Cookies     []cookie `json:"cookies"`
	}

	queueRequest struct {
		Requests []request `json:"request"`
	}
)

func (r *request) toOptions() []entry.Options {
	options := make([]entry.Options, 0)

	cookies := make([]*http.Cookie, len(r.Cookies))
	for i, cookie := range r.Cookies {
		cookies[i] = &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  cookie.Expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: http.SameSite(cookie.SameSite),
		}
	}

	setting := setting.Get()

	options = append(options,
		entry.UseSetting(setting),
		entry.AddCookies(cookies),
		entry.UseDownloader(r.Provider),
		entry.AddHeaders(entry.Headers{
			"Content-Type": r.ContentType,
			"User-Agent":   r.UserAgent,
		}),
	)

	return options
}
