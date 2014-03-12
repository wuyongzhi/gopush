package utils

import (
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	http.Client
}

func TimeoutDialer(cTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		//		if (responseHeaderTimeout > 0) {
		//			conn.SetDeadline(time.Now().Add(rwTimeout))
		//		}
		return conn, nil
	}
}

func NewHttpClient(maxIdleConnsPerHost int,
	connectTimeout time.Duration,
	responseHeaderTimeout time.Duration,
	disableKeepAlives bool) *HttpClient {

	c := HttpClient{}
	c.Transport = &http.Transport{
		Dial: TimeoutDialer(connectTimeout),
		ResponseHeaderTimeout: responseHeaderTimeout,
		DisableKeepAlives:     disableKeepAlives,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
	}

	return &c
}
