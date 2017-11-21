package collopi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request collects all the parameters for a request call.
// It can be instantiated from the client.
type Request struct {
	client       *Client
	method       string
	path         string
	queryValues  *url.Values
	body         interface{}
	target       interface{}
	acceptStati  map[int]bool
	authUser     string
	authPassword string
}

func newRequest(client *Client) *Request {
	return &Request{
		client:      client,
		queryValues: &url.Values{},
		acceptStati: map[int]bool{200: true},
	}
}

// Method sets the method verb used for the request
func (cr *Request) Method(method string) *Request {
	cr.method = method
	return cr
}

// Path will set the path of the request, the parts will be merged using / as separator.
func (cr *Request) Path(paths ...string) *Request {
	if len(paths) == 0 {
		return cr
	}
	if strings.HasPrefix(paths[0], "/") {
		panic("Path must not have / prefix")
	}

	cr.path = "/" + strings.Join(paths, "/")
	return cr
}

// Param adds a param to the url
func (cr *Request) Param(name string, value string) *Request {
	cr.queryValues.Add(name, value)
	return cr
}

// Body will set the body that is used for the request
// will be marshaled to json
func (cr *Request) Body(body interface{}) *Request {
	cr.body = body
	return cr
}

// Accept will add given code to accepted parameters. Default is 200.
// Each call adds another one.
func (cr *Request) Accept(code int) *Request {
	cr.acceptStati[code] = true
	return cr
}

// BasicAuth will decorate the url with basic auth things
func (cr *Request) BasicAuth(user string, password string) *Request {
	cr.authUser = user
	cr.authPassword = password
	return cr
}

// Target will sets a pointer which will be populated with the response body
// using json unmarshal
func (cr *Request) Target(target interface{}) *Request {
	cr.target = target
	return cr
}

func (cr *Request) readToTarget(reader io.Reader, code int) error {
	if cr.target == nil {
		return nil
	}
	if code == http.StatusNotFound {
		return nil
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, cr.target)
	if err != nil {
		str := string(bytes)
		if len(str) > 100 {
			str = str[:100] + "..."
		}
		return fmt.Errorf("Could not parse response as json %v (%v) int %Tv", err, str, cr.target)
	}
	return nil
}

func (cr *Request) buildBody() (io.Reader, error) {
	if cr.body == nil {
		return nil, nil
	}
	jsonBytes, err := json.Marshal(cr.body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jsonBytes)
	return reader, nil
}

// Do executes the request
func (cr *Request) Do() (int, error) {
	body, err := cr.buildBody()
	if err != nil {
		return 0, err
	}

	request, err := cr.buildRequest(body)
	if err != nil {
		return 0, err
	}

	resp, err := cr.client.http.Do(request)
	if resp != nil {
		defer resp.Body.Close() // nolint: errcheck
	}
	if err != nil {
		return 0, err
	}

	if !cr.acceptStati[resp.StatusCode] {
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return 0, fmt.Errorf("Response status is %v, failed to fetch body %v", resp.StatusCode, err2)
		}
		response := string(body)
		return 0, fmt.Errorf("Response error is %v for request %v %v, %v",
			resp.StatusCode, cr.method, cr.client.url+cr.path, response)
	}
	err = cr.readToTarget(resp.Body, resp.StatusCode)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func (cr *Request) buildRequest(body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(cr.method, cr.client.url+cr.path, body)
	if err != nil {
		return nil, err
	}

	if cr.authUser != "" && cr.authPassword != "" {
		request.SetBasicAuth(cr.authUser, cr.authPassword)
	}

	if cr.target != nil {
		request.Header.Add("Accept", "application/json")
	}
	if cr.body != nil {
		request.Header.Add("Content-Type", "application/json")
	}
	request.URL.RawQuery = cr.queryValues.Encode()

	return request, nil
}
