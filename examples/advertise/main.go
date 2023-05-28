package main

import (
	"context"
	"fmt"

	"github.com/cha-king/ssdp"
)

var services = []ssdp.Service{
	{
		Type:     "uuid:3d858ef0-5047-47df-b5a8-62234e314ac0",
		Name:     "uuid:3d858ef0-5047-47df-b5a8-62234e314ac0",
		Location: "https://127-0-0-1.3d858ef0-5047-47df-b5a8-62234e314ac0.local.cha-king.com:8081",
	},
	{
		Type:     "verkada:hubdevice",
		Name:     "uuid:3d858ef0-5047-47df-b5a8-62234e314ac0::verkada:hubdevice",
		Location: "https://127-0-0-1.3d858ef0-5047-47df-b5a8-62234e314ac0.local.cha-king.com:8081",
	},
}

func main() {
	ctx := context.Background()
	errorsChan := make(chan error)
	go ssdp.Advertise(ctx, services, errorsChan)
	for err := range errorsChan {
		fmt.Println(err)
	}
}
