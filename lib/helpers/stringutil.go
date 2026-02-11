package helpers

import (
	"log"
	"strconv"
	"strings"
)

var htmlEsc = strings.NewReplacer(
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;",
	"\"", "&#34;",
	"'", "&#39;",
	"\n", "<br>",
	"\r", "",
)

// ParseFloat parses a float64 from a string with error logging.
func ParseFloat(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if len(val) > 0 && err != nil {
		log.Printf("error parse float: '%s'", val)
	}

	return res
}

// FormatFloat formats a float64 as a string.
func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 6, 64)
}

// ParseInt64 parses an int64 from a string with error logging.
func ParseInt64(val string) int64 {
	res, err := strconv.ParseInt(val, 32, 64)
	if len(val) > 0 && err != nil {
		log.Printf("error parse int: '%s'", val)
	}

	return res
}

// FormatInt64 formats an int64 as a string.
func FormatInt64(val int64) string {
	return strconv.FormatInt(val, 32)
}

// ReplaceToHtml replaces special characters with HTML entities.
func ReplaceToHtml(val string) string {
	return htmlEsc.Replace(val)
}

// HideEmail partially hides an email address for privacy.
func HideEmail(email string) string {
	nameLen := strings.LastIndex(email, "@")

	if nameLen < 0 {
		return ""
	}

	if nameLen < 3 {
		return "***" + email[nameLen:]
	}

	return email[:1] + "***" + email[nameLen-1:]
}
