package inter

import (
	"crypto/tls"
	"net/http"
)

const (
	defaultApiBaseUri = "https://cdpj.partners.bancointer.com.br"
)

type Client struct {
	*http.Client
	apiBaseUrl string
}

func NewClient(c tls.Certificate) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{c},
		},
	}

	return &Client{
		Client:     &http.Client{Transport: tr},
		apiBaseUrl: defaultApiBaseUri,
	}
}
