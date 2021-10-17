package ssdp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

// Additional time to wait for responses beyond provided MX value
const ssdpTimeout = time.Second * 1

// SearchResponse represents the response from an SSDP search request.
type SearchResponse struct {
	Location string
	ST       string
	USN      string
	AL       string
}

func Search(st string, mx int, laddr *net.UDPAddr) ([]*SearchResponse, error) {
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request, err := newSearchRequest(st, mx)
	if err != nil {
		return nil, err
	}

	requestBytes, err := httputil.DumpRequest(request, false)
	if err != nil {
		return nil, err
	}

	_, err = conn.WriteToUDP(requestBytes, ssdpUdpAddr)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(mx)*time.Second + ssdpTimeout
	conn.SetDeadline(time.Now().Add(timeout))
	bufReader := bufio.NewReader(conn)
	sResps := []*SearchResponse{}
	for {
		response, err := http.ReadResponse(bufReader, request)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			break
		}
		if err != nil {
			return nil, err
		}

		sResp := SearchResponse{
			Location: response.Header.Get("Location"),
			ST:       response.Header.Get("ST"),
			USN:      response.Header.Get("USN"),
			AL:       response.Header.Get("AL"),
		}

		sResps = append(sResps, &sResp)
	}

	return sResps, nil
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
