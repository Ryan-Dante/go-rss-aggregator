package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the headers
// Example:
// Authorization: ApiKey {key}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentcation found")
	}

	// Takes the value of the Authorization header and splits it by space
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid authorization header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("invalid first part of auth header")
	}
	return vals[1], nil
}
