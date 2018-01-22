package collopi

import "net/http"

// Client contains the http client and the base url
type Client struct {
	http    *http.Client
	url     string
	debugID string
}

// NewClient creates a client that will access urls
// from given base url
func NewClient(url string, transport ...*http.Transport) *Client {
	client := &Client{
		http: &http.Client{},
		url:  url,
	}
	if len(transport) > 1 {
		panic("maximum one transport allowed")
	}
	if len(transport) == 1 {
		client.http.Transport = transport[0]
	}
	return client
}

// Debug will enable debugging for this client, identifying it with given
// clientID
func (c *Client) Debug(debugID string) {
	c.debugID = debugID
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
