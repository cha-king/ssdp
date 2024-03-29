package ssdp

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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

func Advertise(ctx context.Context, services []Service, errs chan<- error) {
	defer close(errs)

	conn, err := net.ListenMulticastUDP("udp4", nil, ssdpUdpAddr)
	if err != nil {
		errs <- err
		return
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	data := make([]byte, 4096)
	for {
		n, addr, err := conn.ReadFromUDP(data)
		if errors.Is(err, net.ErrClosed) {
			return
		} else if err != nil {
			errs <- err
			continue
		}

		request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data[:n])))
		if err != nil {
			errs <- err
			continue
		}

		go handleRequest(services, conn, addr, request, errs)
	}
}

func handleRequest(services []Service, conn *net.UDPConn, addr *net.UDPAddr, request *http.Request, errs chan<- error) {
	mxStr := request.Header.Get("MX")
	if mxStr == "" {
		errs <- fmt.Errorf("read from %s: mx header missing", addr)
		return
	}
	mx, err := strconv.Atoi(mxStr)
	if err != nil || !(mx >= 1 && mx <= 5) {
		errs <- fmt.Errorf("read from %s: invalid mx value", addr)
		return
	}
	delay := time.Duration(rand.Float64() * float64(mx) * float64(time.Second))
	time.Sleep(delay)

	st := request.Header.Get("ST")
	if st == "" {
		errs <- fmt.Errorf("read from %s: st header missing", addr)
		return
	}

	fmt.Println(addr, st)

	for _, service := range services {
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
			errs <- err
			continue
		}

		_, err = conn.WriteToUDP(respBytes, addr)
		if err != nil {
			errs <- err
			continue
		}
	}
}
