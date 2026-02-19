package jupiter

import "fmt"

// APIError represents an HTTP error response from the Jupiter API.
// It preserves the original status code and raw response body,
// allowing callers to forward the error as-is.
type APIError struct {
	StatusCode int
	RawBody    []byte
	Method     string
	URL        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("call %s() on %s status code: %d", e.Method, e.URL, e.StatusCode)
}