package ghttp

import (
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultTimeout = 8 * time.Second

var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
	//Transport: &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		InsecureSkipVerify: true,
	//	},
	//},
}

type Request struct {
	method       string
	contentType  string
	expectedType string
	uri          string
	payload      interface{}
	header       http.Header
	RequestId    string

	client *http.Client
}

func NewRequest(options ...Option) *Request {
	request := &Request{
		header:    make(http.Header),
		RequestId: uuid.New().String(),
		client:    DefaultClient,
	}

	for _, o := range options {
		o.Apply(request)
	}

	return request
}

func Post(uri string, payload interface{}) *Request {
	return NewRequest().Post(uri, payload)
}

func Get(uri string, payload interface{}) *Request {
	return NewRequest().Get(uri, payload)
}

func (r *Request) Post(uri string, payload interface{}) *Request {
	return r.ContentType(FORM).Method(POST).Uri(uri).Body(payload)
}

func (r *Request) Get(uri string, payload interface{}) *Request {
	return r.ContentType(FORM).Method(GET).Uri(uri).Body(payload)
}

func (r *Request) Send() (*Response, error) {
	params := ""

	switch r.payload.(type) {
	case string:
		params = r.payload.(string)
	case []byte:
		params = string(r.payload.([]byte))
	case Params:
		params = r.payload.(Params).Encode()
	case url.Values:
		params = r.payload.(url.Values).Encode()
	default:
	}
	var request *http.Request
	var err error
	if r.method == GET {
		url := r.uri
		if params != "" {
			url = r.uri + "?" + params
		}
		request, err = http.NewRequest(r.method, url, nil)
	} else {
		request, err = http.NewRequest(r.method, r.uri, strings.NewReader(params))
	}

	if err != nil {
		return nil, err
	}

	if r.contentType != "" {
		request.Header.Set("Content-Type", r.contentType)
	}

	for key, values := range r.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	start := time.Now()

	client := r.getClient()
	rep, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}

	response := NewResponse(rep, body, time.Now().Sub(start))

	response.RequestId = r.RequestId

	return response, nil
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

func (r *Request) AddHeader(name string, value string) *Request {
	r.header.Add(name, value)
	return r
}

func (r *Request) SetHeader(name string, value string) *Request {
	r.header.Set(name, value)
	return r
}

func (r *Request) getClient() *http.Client {
	if r.client == nil {
		r.client = DefaultClient
	}
	return r.client
}

// An Option configures a Request.
type Option interface {
	Apply(*Request)
}

// OptionFunc is a function that configures a Request.
type OptionFunc func(*Request)

// Apply calls f(Request)
func (f OptionFunc) Apply(r *Request) {
	f(r)
}

// WithClient can be used to set the client of a Request to the given value.
func WithClient(client *http.Client) Option {
	return OptionFunc(func(r *Request) {
		r.client = client
	})
}

func WithTransport(transport http.RoundTripper) Option {
	return OptionFunc(func(r *Request) {
		r.getClient().Transport = transport
	})
}

func WithContentType(mime string) Option {
	return OptionFunc(func(r *Request) {
		r.ContentType(mime)
	})
}
