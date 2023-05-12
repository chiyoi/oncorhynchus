package kitsune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseError struct {
	StatusCode int
	Message    string
}

func (re *ResponseError) Error() string {
	return fmt.Sprintf("[%d] %s", re.StatusCode, re.Message)
}

func RoundTrip(endpoint string, req any, resp any) (err error) {
	var hr *http.Response
	if _, ok := req.(Empty); ok {
		if hr, err = http.Get(endpoint); err != nil {
			return
		}
	} else {
		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(req); err != nil {
			return
		}

		hr, err = http.Post(endpoint, "application/json", &buf)
		if err != nil {
			return
		}
	}

	if hr.StatusCode != http.StatusOK {
		re := &ResponseError{
			StatusCode: hr.StatusCode,
		}
		err = re

		if data, err := io.ReadAll(hr.Body); err != nil {
			re.Message = fmt.Sprintf("cannot read body: %s", err)
		} else {
			re.Message = string(data)
		}

		return
	}

	if _, ok := resp.(Empty); ok {
		return
	}

	d := json.NewDecoder(hr.Body)
	return d.Decode(resp)
}

type Empty struct{}
