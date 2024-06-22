// GO API Testing Package
//
// Copyright (c) 2024 Tvative.
// All rights reserved.
//
// Source code and its usage are governed by
// the MIT license.

package apitest

import (
	"net/url"
	"strings"
)

// GenPathParam generates a path parameter string from a map of key-value pairs.
// The key-value pairs are URL-encoded and joined with "&".
// The resulting string is prefixed with "?".
// If the input map is empty, an empty string is returned.
// The function returns an error if "URL-encoding" fails.
//
// Example:
//
//	param, err := GenPathParam(map[string]string{"key1": "value1", "key2": "value2"})
//
// returns:
//
//	"?key1=value1&key2=value2"
func GenPathParam(params map[string]string) (string, error) {
	const emptyParams = ""
	if len(params) == 0 {
		return emptyParams, nil
	}

	var parts []string
	for key, value := range params {
		encodedKey := url.QueryEscape(key)
		encodedValue := url.QueryEscape(value)
		parts = append(parts, encodedKey+"="+encodedValue)
	}

	return "?" + strings.Join(parts, "&"), nil
}
