package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

const ssdpAddress = "239.255.255.250:1900"
const timeout = 5 * time.Second

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

	conn.SetDeadline(time.Now().Add(timeout))
	bufReader := bufio.NewReader(conn)
	responses := []*http.Response{}
	for {
		response, err := http.ReadResponse(bufReader, request)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n\n\n", response)
		responses = append(responses, response)
	}
	fmt.Println(len(responses))
}
