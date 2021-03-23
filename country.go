package phoner

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name                       string `yaml:"name"`
	CountryCode                string `yaml:"country_code"`
	Char2Code                  string `yaml:"char_2_code"`
	Char3Code                  string `yaml:"char_3_code"`
	AreaCode                   string `yaml:"area_code"`
	MaxNumLength               string `yaml:"max_num_length"`
	NationalDialingPrefix      string `yaml:"national_dialing_prefix"`
	InternationalDialingPrefix string `yaml:"international_dialing_prefix"`
}

func (a Config) load() Config {
	filename, _ := filepath.Abs("data/phone/country.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	var config map[string]Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}
