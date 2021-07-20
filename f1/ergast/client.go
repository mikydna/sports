package ergast

import "net/http"

const (
	DefaultHost = "https://ergast.com"
)

type Client struct {
	http *http.Client
	Host string
}

func NewClient(c *http.Client, host string) *Client {
	return &Client{c, host}
}

func DefaultClient() *Client {
	return NewClient(http.DefaultClient, DefaultHost)
}
