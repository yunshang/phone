package phone

import (
	"regexp"
)

const COMMON_EXTENSIONS = `/[ ]*(ext|ex|x|xt|#|:)+[^0-9]*\(*([-0-9]{1,})\)*#?$/i`
const COMMON_EXTRAS = `/(\(0\)|[^0-9+]|^\+?00?)/`

var COMMON_EXTRAS_REPLACEMENTS = map[string]string{
	"00":  "+",
	"+00": "+",
	"+0":  "+",
}

type Phone struct {
	DefaultCountryCode string
	DefaultAreaCode    string
	NamedFormats       string
	N1Length           string
}

func (p Phone) Valid(s string) {

}

func parse(s string) bool {
	if s == "" {
		return false
	}
	// sub, extension := extractExtension(s)
	sub, _ := extractExtension(s)
	sub = normalize(sub)
	return true
}

func extractExtension(s string) (string, string) {
	re := regexp.MustCompile(COMMON_EXTENSIONS)
	subbed := re.FindString(s)
	if subbed == "" {
		return s, ""
	} else {
		return subbed, ""
	}
}

func normalize(stringWithNumber string) string {
	re := regexp.MustCompile(COMMON_EXTRAS)
	match := re.FindAllString(stringWithNumber, -1)
	var s string
	for _, m := range match {
		s = COMMON_EXTRAS_REPLACEMENTS[m]
	}
	return s
}
