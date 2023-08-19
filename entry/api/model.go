package entry

import (
	"net/http"
	"time"

	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/helper"
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
		SameSite string    `json:"sameSite"`
	}

	request struct {
		Url         string `json:"url"`
		ContentType string `json:"contentType"`
		UserAgent   string `json:"userAgent"`
		cookies     []cookie
	}

	cliRequest struct {
		Url string `json:"url"`
	}
)

func (r *request) toOptions() []entry.Options {
	options := make([]entry.Options, 0)

	cookies := make([]*http.Cookie, len(r.cookies))
	for i, cookie := range r.cookies {
		cookies[i] = &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  cookie.Expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: helper.ToSamesite(cookie.SameSite),
		}
	}

	setting := setting.Get()

	options = append(options,
		entry.UseSetting(setting),
		entry.AddCookies(cookies),
		entry.AddHeaders(entry.Headers{
			"Content-Type": r.ContentType,
			"User-Agent":   r.UserAgent,
		}),
	)

	return options
}
