package main

import (
	"fmt"

	"phone/pkg/conf"
)

func main() {
	c := conf.Load()
	fmt.Printf("hello world %v", c)
}
