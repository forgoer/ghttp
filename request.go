package ghttp

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	method       string
	contentType  string
	expectedType string
	uri          string
	payload      interface{}
	header       http.Header
	timeout      time.Duration
}

func NewRequest() *Request {
	request := &Request{
		header: make(http.Header),
	}
	return request
}

func Init(method string) *Request {
	return NewRequest().ContentType("form").Method(method).Timeout(8 * time.Second)
}

func Post(uri string, payload interface{}) *Request {
	return Init(POST).Uri(uri).Body(payload)
}

func Get(uri string, payload interface{}) *Request {
	return Init(GET).Uri(uri).Body(payload)
}

func (r *Request) Send() (*Response, error) {
	client := &http.Client{
		Timeout: r.timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	params := ""

	switch r.payload.(type) {
	case string:
		params = r.payload.(string)
	case Params:
		params = r.payload.(Params).Encode()
	case url.Values:
		params = r.payload.(url.Values).Encode()
	default:
	}
	var request *http.Request
	var err error
	if r.method == GET {
		request, err = http.NewRequest(r.method, r.uri+"?"+params, nil)
	} else {
		request, err = http.NewRequest(r.method, r.uri, strings.NewReader(params))
	}

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", r.contentType)
	for key, values := range r.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	start := time.Now()

	rep, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}

	return NewResponse(rep, body, time.Now().Sub(start)), nil
}

func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// ContentType Set ContentType
func (r *Request) ContentType(mime string) *Request {
	r.contentType = GetFullMime(mime)
	return r
}

func (r *Request) ExpectedType(mime string) *Request {
	r.expectedType = GetFullMime(mime)
	return r
}

func (r *Request) Uri(uri string) *Request {
	r.uri = uri
	return r
}

func (r *Request) Timeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

func (r *Request) Body(payload interface{}) *Request {
	r.payload = payload
	return r
}

func (r *Request) Mime(mime string) *Request {
	r.contentType = GetFullMime(mime)
	r.expectedType = GetFullMime(mime)
	return r
}

func (r *Request) AddHeader(name string, value string) *Request {
	r.header.Add(name, value)
	return r
}

func (r *Request) SetHeader(name string, value string) *Request {
	r.header.Set(name, value)
	return r
}
