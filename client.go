package collopi

import "net/http"

type Client struct {
	http *http.Client
	url  string
}

func NewClient(url string) *Client {
	client := &Client{
		http: &http.Client{},
		url:  url,
	}
	return client
}

func (c *Client) GET() *Request {
	return c.request().Method("GET")
}

func (c *Client) DELETE() *Request {
	return c.request().Method("DELETE")
}

func (c *Client) POST() *Request {
	return c.request().Method("POST")
}

func (c *Client) PUT() *Request {
	return c.request().Method("PUT")
}

func (c *Client) request() *Request {
	return newRequest(c)
}
