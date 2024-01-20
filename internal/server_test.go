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
	"reflect"
	"testing"
	"time"

	server "github.com/admacleod/aws/internal"
)

func TestServerInstantiation(t *testing.T) {
	testSrv := server.New(
		server.Timeout(10 * time.Second),
	)

	for _, tt := range []struct {
		name     string
		got      time.Duration
		expected time.Duration
	}{
		{"ReadTimeout", testSrv.ReadTimeout, 10 * time.Second},
		{"WriteTimeout", testSrv.WriteTimeout, 10 * time.Second},
		{"IdleTimeout", testSrv.IdleTimeout, 10 * time.Second},
	} {
		if tt.got != tt.expected {
			t.Errorf("incorrect %s: expected=%v, got=%v", tt.name, tt.expected, tt.got)
		}
	}
}

func TestServerHandlerOption(t *testing.T) {
	testHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))
	testSrv := server.New(
		server.Handle(testHandler),
	)

	if reflect.DeepEqual(testHandler, testSrv.Handler) {
		t.Errorf("incorrect server handler: expected=%v, got=%v", testHandler, testSrv.Handler)
	}
}
