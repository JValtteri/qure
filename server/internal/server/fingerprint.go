package server

import (
	"fmt"
	"net"
	"net/http"
)

func Fingerprint(request *http.Request) string {
	headers := ""
	headers = addCommonHeaders(headers, request)
	headers = addNamedHeaders(headers, request)
	headers = addTLSheaders(headers, request)
	headers = addAddressHeader(headers, request)
	return headers
}

func addCommonHeaders(headers string, request *http.Request) string {
	headers = headers+request.UserAgent()
	return headers
}

func addNamedHeaders(headers string, request *http.Request) string {
	namedHeaders := []string{
		"Accept-Language",
		"Accept-Encoding",
		"Referer",
		"Origin",
		"Accept",
		"Connection",
		"Upgrade-Insecure-Requests",
		"Sec-CH-UA",
	}
	for _, headerName := range namedHeaders {
		header := request.Header.Get(headerName)
		if header != "" {
			headers = headers + header
		}
	}
	return headers
}

func addTLSheaders(headers string, request *http.Request) string {
	if request.TLS != nil {
		tlsVersion := fmt.Sprintf("%v", request.TLS.Version)
		serverName := request.TLS.ServerName
		headers = headers + tlsVersion + serverName
	}
	return headers
}

func addAddressHeader(headers string, request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	headers = headers + host
	return headers
}
