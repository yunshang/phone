package phone

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const COMMON_EXTENSIONS = `(ext|ex|x|xt|#|:)+[^0-9]*([-0-9]{1,})*#?$`
const COMMON_NUMBER = `[0-9]{1,}$`
const COMMON_EXTRAS = `(\(0\)|[^0-9+]|^\+?00?)`

var FMTENUM = []string{"default", "default_with_extension", "europe", "us"}

const FORMAT_TOKENS = `(%[caAnflx])`

var namedFormat = map[string]string{
	"default":                "+%c%a%n",
	"default_with_extension": "+%c%a%n%x",
	"europe":                 "+%c (0) %a %f %l",
	"us":                     "(%a) %f-%l",
}

var COMMON_EXTRAS_REPLACEMENTS = map[string]string{
	"(0)": "+",
	"00":  "+",
	"+00": "+",
	"+0":  "+",
}

var DefaultCountryCode string
var DefaultAreaCode string

type Phone struct {
	NamedFormats       string
	N1Length           string
	Number             string `yaml:"number"`
	CountryCode        string `yaml:"country_code"`
	AreaCode           string `yaml:"area_code"`
	Extension          string `yaml:"extension"`
	DefaultCountryCode string
	DefaultAreaCode    string
}

func Valid(s string) bool {
	_, err := Parse(s)
	if err != nil {
		return false
	}
	return true
}

func New(args []string) (input *Phone, err error) {
	input = ArgsToCountry(args...)

	if input.N1Length == "" {
		input.N1Length = "3"
	}
	if input.CountryCode == "" {
		input.CountryCode = DefaultCountryCode
	}

	if input.AreaCode == "" {
		input.AreaCode = DefaultAreaCode
	}

	if strings.Trim(input.Number, "\t \n") == "" {
		err = errors.New("Must enter number")
	}
	if strings.Trim(input.AreaCode, "\t \n") == "" {
		err = errors.New("Must enter area code or set default")
	}
	if strings.Trim(input.CountryCode, "\t \n") == "" {
		err = errors.New("Must enter country code or set default")
	}

	return input, err
}

func ArgsToCountry(args ...string) *Phone {
	c := &Phone{}
	switch len(args) {
	case 1:
		c.Number = args[0]
	case 2:
		c.Number = args[0]
		c.AreaCode = args[1]
	case 3:
		c.Number = args[0]
		c.AreaCode = args[1]
		c.CountryCode = args[2]
	case 4:
		c.Number = args[0]
		c.AreaCode = args[1]
		c.CountryCode = args[2]
		c.Extension = args[3]
	}

	return c
}

func Parse(s string) (*Phone, error) {
	if s == "" {
		return nil, nil
	}
	sub, e := extractExtension(s)
	sub = normalize(sub)
	args, err := SplitToParts(sub)
	if err != nil {
		return nil, err
	}
	c, err := New(args)
	if err != nil {
		return nil, err
	}
	c.Extension = e
	return c, nil
}

func extractExtension(s string) (string, string) {
	re := regexp.MustCompile(COMMON_EXTENSIONS)
	subbed := re.FindString(s)
	if subbed != "" {
		re = regexp.MustCompile(COMMON_EXTENSIONS)
		s = re.ReplaceAllString(s, "")
		return s, subbed
	} else {
		return s, ""
	}
}

func normalize(stringWithNumber string) string {
	re := regexp.MustCompile(COMMON_EXTRAS)
	match := re.FindAllString(stringWithNumber, -1)
	var s string
	for _, m := range match {
		s = COMMON_EXTRAS_REPLACEMENTS[m]
		stringWithNumber = re.ReplaceAllString(stringWithNumber, s)
	}
	return stringWithNumber
}

func SplitToParts(s string) (args []string, err error) {
	c := detectCountry(s, DefaultCountryCode)
	fmt.Printf("country is %v \n", c)
	if c != nil {
		re := c.CountryCodeRegexp()
		s = re.ReplaceAllString(s, "0")
		c.CountryCode = "+" + c.CountryCode
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

	exp := fmt.Sprintf("%s", c.AreaCode)
	r, _ := regexp.Compile(exp)
	areaCode := r.FindString(s)

	nExp := fmt.Sprintf("^0*(%s)", c.AreaCode)
	n, _ := regexp.Compile(nExp)
	number := n.ReplaceAllString(s, "")

	args = append(args, number)
	args = append(args, areaCode)
	args = append(args, c.CountryCode)

	return args, nil
}

func (c *Phone) ToS() string {
	return c.format("default")
}

func (c *Phone) Number1() string {
	data := []byte(c.Number)
	i, err := strconv.Atoi(c.N1Length)
	if err != nil {
		panic(err)
	}
	str := string(data[0:i])

	return str
}
func (c *Phone) Number2() string {
	data := []byte(c.Number)
	i, err := strconv.Atoi(c.N1Length)
	if err != nil {
		panic(err)
	}
	l := len(data) - i - 1
	str := string(data[l:])

	return str
}

//Formats the phone number.
func (c *Phone) format(fmt string) (s string) {
	if contains(FMTENUM, fmt) {
		s = c.FormatNumber(namedFormat[fmt])
	} else {
		s = c.FormatNumber(fmt)
	}
	return s
}

func (c Phone) AreaCodeLong() string {
	if c.AreaCode != "" {
		return fmt.Sprintf("0%s", c.AreaCode)
	}
	return ""
}

func (c *Phone) FormatNumber(fm string) string {
	var replacements = map[string]string{
		"%c": c.CountryCode,
		"%a": c.AreaCode,
		"%A": c.AreaCodeLong(),
		"%n": c.Number,
		"%f": c.Number1(),
		"%l": c.Number2(),
		"%x": c.Extension,
	}
	re := regexp.MustCompile(FORMAT_TOKENS)
	match := re.FindAllString(fm, -1)
	fmt.Printf("match is %v \n", match)
	for _, m := range match {
		_s := replacements[m]
		fm = strings.Replace(fm, m, _s, 1)
	}
	fm = RemoveUselessPlus(fm)

	return fm
}

func RemoveUselessPlus(s string) string {
	re := regexp.MustCompile(`^(\+ \+)|^(\+\+)`)
	s = re.ReplaceAllString(s, "+")
	return s
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SetDefaultCountryCode(code string) string {
	DefaultCountryCode = code
	return code
}

func SetDefaultAreaCode(code string) string {
	DefaultAreaCode = code
	return code
}
