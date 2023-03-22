package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RoundTrip(endpoint string, req any, resp any) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("round trip: %w", err)
		}
	}()

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(req); err != nil {
		return
	}

	httpResp, err := http.Post(endpoint, "application/json", &buf)
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-ok response: %s", httpResp.Status)
		return
	}

	d := json.NewDecoder(httpResp.Body)
	return d.Decode(resp)
}
