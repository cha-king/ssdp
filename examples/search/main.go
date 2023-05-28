package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cha-king/ssdp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	responses := make(chan ssdp.SearchResponse)
	errorsChan := make(chan error)
	go ssdp.Search(ctx, ssdp.All, 3, nil, responses, errorsChan)

	for {
		select {
		case response, ok := <-responses:
			if !ok {
				return
			}
			j, err := json.MarshalIndent(response, "", "    ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		case err, ok := <-errorsChan:
			if !ok {
				return
			}
			fmt.Println(err)
		}
	}
}

func PrintStruct(val interface{}) {
	j, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
