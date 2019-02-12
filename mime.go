package chttp

const JSON = "application/json"
const XML = "application/xml"
const XHTML = "application/html+xml"
const FORM = "application/x-www-form-urlencoded"
const UPLOAD = "multipart/form-data"
const PLAIN = "text/plain"
const JS = "text/javascript"
const HTML = "text/html"
const YAML = "application/x-yaml"
const CSV = "text/csv"

var Mimes = map[string]string{
	"json":       JSON,
	"xml":        XML,
	"form":       FORM,
	"plain":      PLAIN,
	"text":       PLAIN,
	"upload":     UPLOAD,
	"html":       HTML,
	"xhtml":      XHTML,
	"js":         JS,
	"javascript": JS,
	"yaml":       YAML,
	"csv":        CSV,
}

// GetFullMime Get the full Mime Type name from a "short name".
func GetFullMime(shortName string) string {
	if _, ok := Mimes[shortName]; ok {
		shortName = Mimes[shortName]
	}

	return shortName
}

// SupportsMimeType Determine whether it supports
func SupportsMimeType(shortName string) bool {
	if _, ok := Mimes[shortName]; ok {
		return true
	}
	return false
}
