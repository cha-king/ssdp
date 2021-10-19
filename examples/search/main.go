package main

import (
	"encoding/json"
	"fmt"

	"github.com/cha-king/ssdp"
)

func main() {
	responses, err := ssdp.Search("all", 3, nil)
	if err != nil {
		panic(err)
	}
	for _, response := range responses {
		j, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))

		// fmt.Printf("%+v\n\n", response)
	}
}

func PrintStruct(val interface{}) {
	j, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
