package ghttp

import (
	"bytes"
	"net/url"
)

type Params struct {
	params []*Param
}

type Param struct {
	Key   string
	Value []string
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v *Params) Get(key string) string {
	for _, param := range v.params {
		if param.Key == key {
			if len(param.Value) == 0 {
				return ""
			}
			return param.Value[0]
		}
	}
	return ""
}

// Has exists the key to value.
func (v *Params) Has(key string) bool {
	for _, param := range v.params {
		if param.Key == key {
			return true
		}
	}
	return false
}

// Index gets the index by key.
func (v *Params) Index(key string) int {
	for i, param := range v.params {
		if param.Key == key {
			return i
		}
	}
	return -1
}

// Set sets the key to value. It replaces any existing
// values.
func (v *Params) Set(key, value string) {
	if i := v.Index(key); i != -1 {
		v.params[i].Value = []string{value}
	} else {
		v.params = append(v.params, &Param{
			Key:   key,
			Value: []string{value},
		})
	}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v *Params) Add(key, value string) {
	if i := v.Index(key); i != -1 {
		v.params[i].Value = append(v.params[i].Value, value)
	} else {
		v.params = append(v.params, &Param{
			Key:   key,
			Value: []string{value},
		})
	}
}

// Del deletes the values associated with key.
func (v *Params) Del(key string) {
	if i := v.Index(key); i != -1 {
		v.params = append(v.params[:i], v.params[i+1:]...)
	}
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (v Params) Encode() string {
	if v.params == nil {
		return ""
	}
	var buf bytes.Buffer
	for _, param := range v.params {
		prefix := url.QueryEscape(param.Key) + "="
		for _, v := range param.Value {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
