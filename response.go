package ghttp

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	HttpResponse *http.Response
	Body         string
	RawBody      []byte
	Code         int
	Duration     time.Duration
}

func NewResponse(resp *http.Response, body []byte, duration time.Duration) *Response {
	response := &Response{
		HttpResponse: resp,
		RawBody:      body,
		Code:         resp.StatusCode,
		Body:         string(body),
		Duration:     duration,
	}
	return response
}

func (r *Response) Unmarshal(v interface{}) error {
	return json.Unmarshal(r.RawBody, v)
}

func (r *Response) ToString() string {
	return r.Body
}
