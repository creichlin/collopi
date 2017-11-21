package collopi

import "net/http"

// Client contains the http client and the base url
type Client struct {
	http *http.Client
	url  string
}

// NewClient creates a client that will access urls
// from given base url
func NewClient(url string) *Client {
	client := &Client{
		http: &http.Client{},
		url:  url,
	}
	return client
}

// GET creates a GET request
func (c *Client) GET() *Request {
	return c.request().Method("GET")
}

// DELETE creates a DELETE request
func (c *Client) DELETE() *Request {
	return c.request().Method("DELETE")
}

// POST creates a POST request
func (c *Client) POST() *Request {
	return c.request().Method("POST")
}

// PUT creates a PUT request
func (c *Client) PUT() *Request {
	return c.request().Method("PUT")
}

// Method creates a request using given method string
func (c *Client) Method(m string) *Request {
	return c.request().Method(m)
}

func (c *Client) request() *Request {
	return newRequest(c)
}
