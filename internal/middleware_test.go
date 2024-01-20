// Copyright (c) Alisdair MacLeod <copying@alisdairmacleod.co.uk>
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package internal_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	server "github.com/admacleod/aws/internal"
)

func TestMiddleware(t *testing.T) {
	var testMiddlewareCalls []string

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			testMiddlewareCalls = append(testMiddlewareCalls, "middleware1")
			next.ServeHTTP(w, r)
		})
	}
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			testMiddlewareCalls = append(testMiddlewareCalls, "middleware2")
			next.ServeHTTP(w, r)
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMiddlewareCalls = append(testMiddlewareCalls, "handler")
	})

	req := httptest.NewRequest("GET", "http://test.example.com", nil)
	w := httptest.NewRecorder()

	testMiddleware := server.ChainMiddleware(middleware1, middleware2)
	testMiddleware(handler).ServeHTTP(w, req)

	expect := []string{
		"middleware2",
		"middleware1",
		"handler",
	}
	if !reflect.DeepEqual(expect, testMiddlewareCalls) {
		t.Errorf("incorrect middleware calls: expected=%v, got=%v", expect, testMiddlewareCalls)
	}
}
