package phone

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Country struct {
	Number                     string `yaml:"number"`
	Name                       string `yaml:"name"`
	CountryCode                string `yaml:"country_code"`
	Char2Code                  string `yaml:"char_2_code"`
	Char3Code                  string `yaml:"char_3_code"`
	AreaCode                   string `yaml:"area_code"`
	MaxNumLength               string `yaml:"max_num_length"`
	NationalDialingPrefix      string `yaml:"national_dialing_prefix"`
	InternationalDialingPrefix string `yaml:"international_dialing_prefix"`
	Extension                  string `yaml:"extension"`
}

var Countries map[string]Country

func init() {
	Countries = Load()
}

func Load() map[string]Country {
	filename, _ := filepath.Abs("./data/phone/countries.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	var c map[string]Country

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return c
}

func detectCountry(s string) (c *Country) {
	for _, v := range Countries {
		r := v.CountryCodeRegexp()
		if r.MatchString(s) {
			c = &v
		}
	}

	return c
}
func FindByCountryCode(code string) Country {
	return Countries[code]
}

func FindByCountryIsoCode(isocode string) (c Country) {
	for _, v := range Countries {
		if isocode == strings.ToLower(v.Char3Code) {
			c = v
		}
	}
	return c
}

func (c Country) CountryCodeRegexp() *regexp.Regexp {
	exp := fmt.Sprintf("/^[+]%s/", c.CountryCode)
	re, _ := regexp.Compile(exp)

	return re
}

func (c Country) Formats() (*regexp.Regexp, *regexp.Regexp) {
	numberRegex := fmt.Sprintf("([0-9]{1,%s})$", c.MaxNumLength)
	short := regexp.MustCompile(fmt.Sprintf("/^0?(%s)%s/", c.AreaCode, numberRegex))
	reallyShort := regexp.MustCompile(fmt.Sprintf("/^%s", numberRegex))

	return short, reallyShort
}

func (c Country) DetectFormat(stringWithNumber string) string {
	sh, real := c.Formats()
	var arr []string

	if sh.MatchString(stringWithNumber) {
		arr = append(arr, "short")
	}
	if real.MatchString(stringWithNumber) {
		arr = append(arr, "reall_short")
	}
	if len(arr) > 1 {
		return "reall_short"
	}
	if len(arr) == 0 {
		return ""
	}

	return arr[0]
}

func New(args ...string) (input *Country, err error) {
	input = ArgsToCountry(args...)

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

func ArgsToCountry(args ...string) *Country {
	c := &Country{}
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
		c.Extension = args[3]
	}

	return c
}

func (c Country) AreaCodeLong() string {
	if c.AreaCode != "" {
		return fmt.Sprintf("0%s", c.AreaCode)
	}
	return ""
}
