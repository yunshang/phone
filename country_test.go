package phone

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	c := Load()
	fmt.Printf("country %v", c)
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
	c, err := Parse("+385915125486")
	fmt.Printf("phone is %v", c)
	fmt.Printf("error is %v", err)
}

func TestNormalize(t *testing.T) {
	c := normalize("+00385915125486")
	fmt.Printf("string is %v", c)
}
