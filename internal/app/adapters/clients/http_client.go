package adapters

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func NewInsecureHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxConnsPerHost:     10,
		MaxIdleConnsPerHost: 10,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: transport, Timeout: 15 * time.Second}
}
