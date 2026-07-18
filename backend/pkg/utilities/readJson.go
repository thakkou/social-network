package utilities

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ReadJSONRequest[T any](r *http.Request) (T, error) {
	var dst T

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&dst)
		return dst, err
	}

	err := r.ParseForm()
	if err != nil {
		return dst, err
	}

	for key := range r.Form {
		err := json.Unmarshal([]byte(key), &dst)
		return dst, err
	}

	return dst, nil
}

func ReadJSONRequestIntoStruct(r *http.Request, dst any) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		return decoder.Decode(dst)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	for key := range r.Form {
		return json.Unmarshal([]byte(key), dst)
	}

	return nil
}
