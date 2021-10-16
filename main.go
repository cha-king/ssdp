package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"

	"golang.org/x/net/ipv4"
)

const ssdpAddress = "239.255.255.250:1900"

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp4", ssdpAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println(udpAddr)

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request, err := http.NewRequest("M-SEARCH", "*", nil)
	if err != nil {
		panic(err)
	}
	request.Host = ssdpAddress
	request.Header.Set("ST", "ssdp:all")
	request.Header.Set("MAN", `"ssdp:discover"`)
	request.Header.Set("MX", "3")
	raw, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))

	_, err = conn.WriteToUDP(raw, udpAddr)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)
	n, _, addr, err := pConn.ReadFrom(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
	fmt.Println(string(buffer[:n]))
}
