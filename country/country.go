package country

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Country struct {
	Name                       string `yaml:"name"`
	CountryCode                string `yaml:"country_code"`
	Char2Code                  string `yaml:"char_2_code"`
	Char3Code                  string `yaml:"char_3_code"`
	AreaCode                   string `yaml:"area_code"`
	MaxNumLength               string `yaml:"max_num_length"`
	NationalDialingPrefix      string `yaml:"national_dialing_prefix"`
	InternationalDialingPrefix string `yaml:"international_dialing_prefix"`
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
	re := regexp.MustCompile(exp)

	return re
}

func (c Country) Formats() (reg map[string]*regexp.Regexp) {
	number_regex := fmt.Sprintf("([0-9]{1,%d})$", c.MaxNumLength)
	short := regexp.MustCompile(fmt.Sprintf("/^0?(%s)%s/", c.AreaCode, number_regex))
	really_short := regexp.MustCompile(fmt.Sprintf("/^%s", number_regex))
	reg["short"] = short
	reg["really_short"] = really_short

	return reg
}

// def self.formats(country)
// area_code_regexp = country.area_code
// number_regex     = "([0-9]{1,#{country.max_num_length}})$".freeze

// {
//   # 047451588, 013668734
//   :short => Regexp.new("^0?(#{area_code_regexp})#{number_regex}"),
//   # 451588
//   :really_short => Regexp.new("^#{number_regex}")
// }
// end
