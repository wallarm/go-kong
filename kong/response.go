package kong

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response is a Kong Admin API response. It wraps http.Response.
type Response struct {
	*http.Response
	// other Kong specific fields
}

func newResponse(res *http.Response) *Response {
	return &Response{Response: res}
}

func messageFromBody(b []byte) string {
	s := struct {
		Message string
	}{}

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Sprintf("<failed to parse response body: %v>", err)
	}

	return s.Message
}

func hasError(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 399 {
		return nil
	}

	body, _ := io.ReadAll(res.Body) // TODO error in error?
	return NewAPIError(res.StatusCode, messageFromBody(body))
}
