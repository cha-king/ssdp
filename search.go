package ssdp

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
)

// SearchResponse represents the response from an SSDP search request.
type SearchResponse struct {
	Location          string
	SearchTarget      string
	UniqueServiceName string
}

func Search(ctx context.Context, searchTarget string, mx int, laddr *net.UDPAddr, responses chan<- SearchResponse, errorsChan chan<- error) {
	defer close(responses)
	defer close(errorsChan)

	if !(mx >= 1 && mx <= 5) {
		errorsChan <- errors.New("mx must be between 1 and 5")
		return
	}

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

	request, err := newSearchRequest(searchTarget, mx)
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
			Location:          response.Header.Get("Location"),
			SearchTarget:      response.Header.Get("ST"),
			UniqueServiceName: response.Header.Get("USN"),
		}
		responses <- sResp
	}
}

func newSearchRequest(searchTarget string, mx int) (*http.Request, error) {
	request, err := http.NewRequest("M-SEARCH", "*", nil)
	if err != nil {
		return nil, err
	}

	request.Host = ssdpUdpAddr.String()

	request.Header.Set("MAN", `"ssdp:discover"`)
	request.Header.Set("ST", searchTarget)
	request.Header.Set("MX", strconv.Itoa(mx))

	return request, nil
}
