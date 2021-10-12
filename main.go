package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const ssdpAddress = "239.255.255.250:1900"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp4", ssdpAddress)
	if err != nil {
		panic(err)
	}

	// server, err := net.ListenMulticastUDP("udp4", nil, udpAddr)
	server, err := net.ListenUDP("udp4", &net.UDPAddr{IP: nil, Port: 1900})
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, server)

	conn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request, err := http.NewRequest("M-SEARCH", "*", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("ST", "ssdp:all")
	request.Header.Set("MAN", "ssdp:discover")
	request.Header.Set("MX", "120")
	raw, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))

	n, err := conn.Write(raw)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)

	time.Sleep(1 * time.Minute)
}
