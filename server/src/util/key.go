package util

import (
	"regexp"
	"strings"
)

var pattrenNoCache = regexp.MustCompile(`(\s+|^)no-cache(\s+|$)`)

// KeyNoCache ...
func KeyNoCache(key string) bool {
	return pattrenNoCache.MatchString(key)
}

// FormatKey ...
func FormatKey(key string) string {
	key = pattrenNoCache.ReplaceAllString(key, ` `)
	key = strings.TrimSpace(key)
	if key == `` {
		key = `<empty>`
	}
	return key
}
