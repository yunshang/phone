package main

import (
	"fmt"
	"phone/country"
)

func main() {
	c := country.Load()
	fmt.Printf("hello world %v", c)
}
