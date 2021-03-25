package phone

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	c := Load()
	fmt.Printf("hello wolrd %v", c)
}
