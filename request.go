package chttp

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	method       string
	contentType  string
	expectedType string
	uri          string
	payload      interface{}
}

func NewRequest() *Request {
	request := &Request{}
	request = request.ContentType("form")
	return request
}

func Init(method string) *Request {
	request := &Request{}
	request.ContentType("form").Method(method)
	return request
}

func Post(uri string, payload interface{}) *Request {
	return Init("POST").Uri(uri).Body(payload)
}

func (r *Request) Send() (*Response, error) {

	client := &http.Client{}
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
	var response *Response

	req, err := http.NewRequest(r.method, r.uri, strings.NewReader(params))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", r.contentType)

	rep, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer rep.Body.Close()

	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return response, err
	}

	return NewResponse(rep, body), nil
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

func (r *Request) Body(payload interface{}) *Request {
	r.payload = payload
	return r
}

func (r *Request) Mime(mime string) *Request {
	r.contentType = GetFullMime(mime)
	r.expectedType = GetFullMime(mime)
	return r
}
