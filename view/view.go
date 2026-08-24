package view

import "net/http"

type Response struct {
	ExtraFields map[string]any `json:"-"`
	Data        any            `json:"data"`
	Error       string         `json:"error,omitempty"`
	Message     string         `json:"message,omitempty"`
	Success     bool           `json:"success,omitempty"`
	HTTPCode    int            `json:"httpCode,omitempty"`
}

type JSON struct {
	Fields map[string]any
}

func View(json *JSON) *Response {
	return &Response{
		ExtraFields: json.Fields,
		Success:     true,
		HTTPCode:    http.StatusOK,
	}
}
