package helper

import (
	"net/http"
	"strings"
)

func ToSamesite(val string) http.SameSite {
	switch strings.ToLower(val) {
	case "no_restriction":
		return http.SameSiteDefaultMode
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	default:
		return http.SameSiteDefaultMode
	}
}
