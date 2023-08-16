package entry

import (
	"net/http"
	"time"

	"github.com/rapid-downloader/rapid/helper"
)

type (
	cookie struct {
		Name       string
		Value      string
		Path       string
		Domain     string
		Expires    time.Time
		RawExpires string
		MaxAge     int
		Secure     bool
		HttpOnly   bool
		SameSite   string
		Raw        string
		Unparsed   []string
	}

	browserRequest struct {
		Url         string `json:"url"`
		ContentType string `json:"contentType"`
		UserAgent   string `json:"userAgent"`
		cookies     []cookie
	}
)

func (r *browserRequest) toOptions() []Options {
	options := make([]Options, 0)

	cookies := make([]*http.Cookie, len(r.cookies))
	for i, cookie := range r.cookies {
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
			SameSite:   helper.ToSamesite(cookie.SameSite),
			Raw:        cookie.Raw,
			Unparsed:   cookie.Unparsed,
		}
	}

	options = append(options,
		AddCookies(cookies),
		AddHeaders(Headers{
			"Content-Type": r.ContentType,
			"User-Agent":   r.UserAgent,
		}),
	)

	return options
}
