package ssdp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type Service struct {
	Name     string
	Type     string
	Location string
}

func Advertise(ctx context.Context, services []Service, errorsChan chan<- error) {
	defer close(errorsChan)

	conn, err := net.ListenMulticastUDP("udp4", nil, ssdpUdpAddr)
	if err != nil {
		errorsChan <- err
		return
	}

	data := make([]byte, 4096)
	for {
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			errorsChan <- err
			continue
		}
		fmt.Println("Request received")

		reader := bytes.NewReader(data[:n])

		request, err := http.ReadRequest(bufio.NewReader(reader))
		if err != nil {
			errorsChan <- err
			continue
		}

		mxStr := request.Header.Get("MX")
		if mxStr == "" {
			errorsChan <- fmt.Errorf("read from %s: mx header missing", addr)
			continue
		}
		mx, err := strconv.Atoi(mxStr)
		if err != nil || !(mx >= 1 && mx <= 5) {
			errorsChan <- fmt.Errorf("read from %s: invalid mx value", addr)
			continue
		}
		delay := time.Duration(rand.Float64() * float64(mx) * float64(time.Second))
		time.Sleep(delay)

		st := request.Header.Get("ST")

		for _, service := range services {
			fmt.Println("HERE")
			fmt.Println(st)
			if st != "ssdp:all" && st != service.Type {
				continue
			}

			resp := &http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Ext":      []string{""},
					"Location": []string{service.Location},
					"ST":       []string{service.Type},
					"USN":      []string{service.Name},
				},
			}
			respBytes, err := httputil.DumpResponse(resp, true)
			if err != nil {
				errorsChan <- err
				continue
			}

			fmt.Println(string(respBytes))

			_, err = conn.WriteToUDP(respBytes, addr)
			if err != nil {
				errorsChan <- err
				continue
			}
		}
	}
}
