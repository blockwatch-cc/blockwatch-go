// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPError retains HTTP status
type HTTPError interface {
	error
	Status() string  // e.g. "200 OK"
	StatusCode() int // e.g. 200
	Body() []byte
	Unmarshal(val interface{}) error
}

// RateLimitError helps manage rate limit errors
type RateLimitError interface {
	HTTPError
	Wait(context.Context) error
	Done() <-chan struct{}
	Deadline() time.Duration
}

func IsRateLimited(err error) (RateLimitError, bool) {
	e, ok := err.(RateLimitError)
	return e, ok
}

type httpError struct {
	status int
	body   []byte
	text   string
}

func (e *httpError) Status() int {
	return e.status
}

func (e *httpError) Error() string {
	bodyText := e.text
	if bodyText == "" {
		bodyText = strings.TrimRight(string(e.body[:256]), "\x00")
		bodyText = strings.Replace(bodyText, "\n", "", -1)
		bodyText = strings.Replace(bodyText, "  ", " ", -1)
	}
	return fmt.Sprintf("%d %s: %s", e.status, http.StatusText(e.status), bodyText)
}

func (e *httpError) Unmarshal(val interface{}) error {
	return json.Unmarshal([]byte(e.body), val)
}

type errRateLimited struct {
	httpError
	deadline time.Time
	done     chan struct{}
}

func (e *errRateLimited) timeout() {
	select {
	case <-time.After(time.Until(e.deadline)):
		close(e.done)
	}
}

func newErrRateLimited(err *httpError, until int64) *errRateLimited {
	e := &errRateLimited{
		httpError: *err,
		deadline:  time.Unix(until, 0),
		done:      make(chan struct{}),
	}
	if e.deadline.After(time.Now()) {
		go e.timeout()
	} else {
		close(e.done)
	}
	return e
}

func (e *errRateLimited) Error() string {
	return e.Error()
}

func (e *errRateLimited) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.done:
		return nil
	}
}

func (e *errRateLimited) Done() <-chan struct{} {
	return e.done
}

func (e *errRateLimited) Deadline() time.Duration {
	return e.deadline.Sub(time.Now().UTC())
}

type Error struct {
	Code      int    `json:"code"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Scope     string `json:"scope"`
	Detail    string `json:"detail"`
	RequestId string `json:"requestId"`
	Reason    string `json:"reason"`
}

func (e *Error) Error() string {
	s := make([]string, 0)
	if e.Status != 0 {
		s = append(s, "status="+strconv.Itoa(e.Status))
	}
	if e.Code != 0 {
		s = append(s, "code="+strconv.Itoa(e.Code))
	}
	if e.Scope != "" {
		s = append(s, "scope="+e.Scope)
	}
	s = append(s, "message=\""+e.Message+"\"")
	if e.Detail != "" {
		s = append(s, "detail=\""+e.Detail+"\"")
	}
	if e.RequestId != "" {
		s = append(s, "request-id="+e.RequestId)
	}
	if e.Reason != "" {
		s = append(s, "reason=\""+e.Reason+"\"")
	}
	return strings.Join(s, " ")
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func (e Errors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	return e.Errors[0].Error()
}

type apiError struct {
	*httpError
	errors Errors
}

func (e apiError) Error() string {
	return e.errors.Error()
}

func (e *apiError) Errors() []Error {
	return e.errors.Errors
}
