package phone

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	c := Load()
	fmt.Printf("c lenth is  %d", len(c))
}

func TestExtractExtension(t *testing.T) {
	s1, s2 := extractExtension("+385915125486")
	fmt.Printf("s1 is %s and s2 is %s", s1, s2)
}

func TestValid(t *testing.T) {
	c := Valid("+385915125486")
	fmt.Printf("phone is %v", c)
}

func TestParse(t *testing.T) {
	c, err := Parse("+00385915125486")
	s := c.ToS()
	fmt.Printf("phone is %v \n", c)
	fmt.Printf("error is %v \n", err)
	fmt.Printf("phone string %v \n", s)
}

func TestNormalize(t *testing.T) {
	c := normalize("+00385915125486")
	fmt.Printf("string is %v", c)
}

func TestNew(t *testing.T) {
	args := []string{"5125486", "91", "385", "143"}
	c, err := New(args)
	fmt.Printf("error is %v \n", err)
	fmt.Printf("string is %v", c)
}

func TestFormat(t *testing.T) {
	c, err := Parse("+00385915125486x148")
	fmt.Printf("error is %v \n", err)
	f := c.format("%A/%f-%l")
	n := c.format("+ %c (%a) %n")
	europe := c.format("europe")
	us := c.format("us")
	ex := c.format("default_with_extension")
	fmt.Printf("c is %v \n", c)
	fmt.Printf("f is %v \n", f)
	fmt.Printf("n is %v \n", n)
	fmt.Printf("europe is %v \n", europe)
	fmt.Printf("us is %v \n", us)
	fmt.Printf("ex is %v \n", ex)
}
func TestSetDefault(t *testing.T) {
	SetDefaultAreaCode("47")
	SetDefaultCountryCode("385")
	c, err := Parse("451-588")
	fmt.Printf("error is %v \n", err)
	f := c.format("%A/%f-%l")
	n := c.format("+ %c (%a) %n")
	europe := c.format("europe")
	us := c.format("us")
	ex := c.format("default_with_extension")
	fmt.Printf("c is %v \n", c)
	fmt.Printf("f is %v \n", f)
	fmt.Printf("n is %v \n", n)
	fmt.Printf("europe is %v \n", europe)
	fmt.Printf("us is %v \n", us)
	fmt.Printf("ex is %v \n", ex)
}

func TestFindByCountryIsoCode(t *testing.T) {
	f := FindByCountryIsoCode("de")
	fmt.Printf("f is %v \n", f)
}

func TestFindByCountryCode(t *testing.T) {
	f := FindByCountryCode("385222222222")
	fmt.Printf("f is %v \n", f)
}
