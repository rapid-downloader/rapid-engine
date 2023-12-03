package helper

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-json"
)

func UnmarshalBody(r *http.Response, v interface{}) error {
	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		return err
	}

	if err := json.Unmarshal(buffer.Bytes(), &v); err != nil {
		return err
	}

	return nil
}
