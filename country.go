package phone

import (
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
	N1Length                   string
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

func detectCountry(s, defaultCode string) (c *Country) {
	_c := Country{}
	for k, v := range Countries {
		re := fmt.Sprintf("^[+]%s", k)
		matched, _ := regexp.MatchString(re, s)
		if matched {
			_c = v
		}
	}
	if _c.CountryCode == "" {
		_c = Countries[defaultCode]
	}
	c = &_c

	return c
}
func FindByCountryCode(code string) *Country {
	c := Countries[code]
	return &c
}

func FindByCountryIsoCode(isCcode string) (c *Country) {
	for _, v := range Countries {
		if isCcode == strings.ToLower(v.Char3Code) {
			c = &v
		}
	}
	return c
}

func (c Country) CountryCodeRegexp() *regexp.Regexp {
	exp := fmt.Sprintf("^[+]%s", c.CountryCode)
	re, _ := regexp.Compile(exp)

	return re
}

func (c Country) Formats() (*regexp.Regexp, *regexp.Regexp) {
	numberRegex := fmt.Sprintf("([0-9]{1,%s})$", c.MaxNumLength)
	short := regexp.MustCompile(fmt.Sprintf("^0?(%s)%s", c.AreaCode, numberRegex))
	reallyShort := regexp.MustCompile(fmt.Sprintf("^%s", numberRegex))

	return short, reallyShort
}

func (c Country) DetectFormat(stringWithNumber string) string {
	sh, real := c.Formats()
	var arr []string

	if sh.MatchString(stringWithNumber) {
		arr = append(arr, "short")
	}
	if real.MatchString(stringWithNumber) {
		arr = append(arr, "really_short")
	}
	if len(arr) > 1 {
		return "really_short"
	}
	if len(arr) == 0 {
		return "short"
	}

	return arr[0]
}
