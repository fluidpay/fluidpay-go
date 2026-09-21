package fluidpay

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrNotFound is returned when the gateway answers a single-record lookup
// with an empty result set.
var ErrNotFound = errors.New("fluidpay: not found")

// Error is returned when the gateway rejects a request, either with a non-2xx
// HTTP status or with a "failed" envelope. Declined transactions are not
// errors: they come back as a successful envelope with a decline response
// code (see Transaction.Approved).
type Error struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Status is the envelope status, typically "failed" or "error".
	Status string
	// Msg is the human readable message from the gateway, for example
	// "bad request error: invalid Postal Code".
	Msg string
	// CorrelationID is the x-correlation-id header. Include it when
	// contacting support.
	CorrelationID string
	// Method and Path identify the request that failed.
	Method string
	Path   string
	// Body is the raw response body, useful when Msg is empty.
	Body []byte

	header http.Header
}

// Error implements the error interface.
func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "fluidpay: %s /%s: %d", e.Method, strings.TrimLeft(e.Path, "/"), e.StatusCode)
	if text := http.StatusText(e.StatusCode); text != "" {
		b.WriteString(" " + text)
	}
	if e.Msg != "" {
		b.WriteString(": " + e.Msg)
	} else if e.Status != "" {
		b.WriteString(": status " + e.Status)
	}
	if e.CorrelationID != "" {
		b.WriteString(" (correlation id " + e.CorrelationID + ")")
	}
	return b.String()
}

// IsUnauthorized reports whether the gateway rejected the credentials. Common
// causes are a missing or deleted key, a public key used server-side, an IP
// or URL restriction on the key, or the wrong environment.
func (e *Error) IsUnauthorized() bool { return e.StatusCode == http.StatusUnauthorized }

// IsNotFound reports whether the resource does not exist.
func (e *Error) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsBadRequest reports whether the gateway rejected the request payload.
func (e *Error) IsBadRequest() bool { return e.StatusCode == http.StatusBadRequest }

// Response returns the metadata of the HTTP response that produced the
// error, including the correlation ID, status code and headers.
func (e *Error) Response() *APIResponse {
	return &APIResponse{
		Method:        e.Method,
		Path:          e.Path,
		StatusCode:    e.StatusCode,
		CorrelationID: e.CorrelationID,
		Status:        e.Status,
		Msg:           e.Msg,
		Header:        e.header,
		raw:           e.Body,
	}
}

func (e *Error) correlationID() string { return e.CorrelationID }

// decodeError is returned when the gateway answered but the SDK could not
// decode the body. It keeps the response so the correlation ID survives.
type decodeError struct {
	resp *APIResponse
	err  error
}

func (e *decodeError) Error() string {
	msg := fmt.Sprintf("fluidpay: %s /%s: decoding response: %v", e.resp.Method, strings.TrimLeft(e.resp.Path, "/"), e.err)
	if e.resp.CorrelationID != "" {
		msg += " (correlation id " + e.resp.CorrelationID + ")"
	}
	return msg
}

func (e *decodeError) Unwrap() error { return e.err }

func (e *decodeError) correlationID() string { return e.resp.CorrelationID }

// correlated is implemented by errors that carry a correlation ID.
type correlated interface{ correlationID() string }

// CorrelationID returns the x-correlation-id associated with err, or "" when
// err did not come from a gateway response (for example a network failure
// or a validation error raised before any request was sent). It works on
// *Error and on decoding errors, including when they are wrapped.
func CorrelationID(err error) string {
	var c correlated
	if errors.As(err, &c) {
		return c.correlationID()
	}
	return ""
}

// AsError extracts a *Error from err, if there is one.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// IsUnauthorized reports whether err is a gateway authentication failure.
func IsUnauthorized(err error) bool {
	e, ok := AsError(err)
	return ok && e.IsUnauthorized()
}

// IsNotFound reports whether err indicates a missing resource, either a 404
// from the gateway or ErrNotFound.
func IsNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	e, ok := AsError(err)
	return ok && e.IsNotFound()
}

func errNilRequest(what string) error {
	return fmt.Errorf("fluidpay: %s request must not be nil", what)
}
