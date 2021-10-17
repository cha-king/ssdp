package main

import (
	"fmt"

	"github.com/cha-king/ssdp"
)

func main() {
	responses, err := ssdp.Search("all", 3, nil)
	if err != nil {
		panic(err)
	}
	for _, response := range responses {
		fmt.Printf("%+v\n\n", response)
	}
}
