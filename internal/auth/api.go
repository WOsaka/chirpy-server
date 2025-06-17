package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	headerValue := headers.Get("Authorization")
	if headerValue == "" {
		return "", errors.New("missing API key in Authorization header")
	}

	fields := strings.Fields(headerValue)
	if len(fields) == 0 {
		return "", errors.New("invalid Authorization header format")
	}
	apiKey := fields[len(fields)-1]
	return apiKey, nil
}