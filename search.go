package ssdp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
)

// SearchResponse represents the response from an SSDP search request.
type SearchResponse struct {
	Location string
	ST       string
	USN      string
}

func Search(ctx context.Context, st string, mx int, laddr *net.UDPAddr, responses chan<- SearchResponse, errorsChan chan<- error) {
	defer close(responses)
	defer close(errorsChan)

	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		errorsChan <- err
		return
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	request, err := newSearchRequest(st, mx)
	if err != nil {
		errorsChan <- err
		return
	}

	requestBytes, err := httputil.DumpRequest(request, false)
	if err != nil {
		errorsChan <- err
		return
	}

	_, err = conn.WriteToUDP(requestBytes, ssdpUdpAddr)
	if err != nil {
		errorsChan <- err
		return
	}

	bufReader := bufio.NewReader(conn)

	for {
		response, err := http.ReadResponse(bufReader, request)
		if errors.Is(err, net.ErrClosed) {
			return
		} else if err != nil {
			errorsChan <- err
			continue
		}

		sResp := SearchResponse{
			Location: response.Header.Get("Location"),
			ST:       response.Header.Get("ST"),
			USN:      response.Header.Get("USN"),
		}
		responses <- sResp
	}
}

func newSearchRequest(st string, mx int) (*http.Request, error) {
	request, err := http.NewRequest("M-SEARCH", "*", nil)
	if err != nil {
		return nil, err
	}

	request.Host = ssdpUdpAddr.String()

	request.Header.Set("MAN", `"ssdp:discover"`)
	request.Header.Set("ST", fmt.Sprintf("ssdp:%s", st))
	request.Header.Set("MX", strconv.Itoa(mx))

	return request, nil
}
