package phone

import (
	"errors"
	"fmt"
	"regexp"
)

const COMMON_EXTENSIONS = `(ext|ex|x|xt|#|:)+[^0-9]*([-0-9]{1,})*#?$`
const COMMON_NUMBER = `[0-9]{1,}$`
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

func (p Phone) Valid(s string) bool {
	_, err := parse(s)
	if err != nil {
		return false
	}
	return true
}

func parse(s string) (*Country, error) {
	if s == "" {
		return nil, nil
	}
	sub, _ := extractExtension(s)
	sub = normalize(sub)
	args, err := SplitToParts(sub)
	if err != nil {
		return nil, err
	}
	c, err := New(args...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func extractExtension(s string) (string, string) {
	re := regexp.MustCompile(COMMON_EXTENSIONS)
	subbed := re.FindString(s)
	if subbed != "" {
		re = regexp.MustCompile(COMMON_EXTENSIONS)
		_s := re.FindString(subbed)
		return s, _s
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

func SplitToParts(s string) (args []string, err error) {
	c := detectCountry(s)
	if c != nil {
		//countryCode := c.CountryCode
		//string = string.gsub(country.country_code_regexp, "0")
		re := c.CountryCodeRegexp()
		s = re.ReplaceAllString(s, "0")
	}

	if c == nil {
		e := fmt.Sprint("Must enter country code or set default country code")
		err = errors.New(e)
		return nil, err
	}

	format := c.DetectFormat(s)
	if format == "" {
		return nil, err
	}

	//parts = string.match formats(country)[format]
	sh, real := c.Formats()

	switch format {
	case "short":
		p := sh.FindAllString(s, -1)
		args = append(args, p[2])
		args = append(args, p[1])
		args = append(args, c.CountryCode)
	case "really_short":
		re := real.FindAllString(s, -1)
		args = append(args, re[1])
		args = append(args, c.AreaCode)
		args = append(args, c.CountryCode)
	}

	return args, nil
}
