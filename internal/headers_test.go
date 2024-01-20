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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	server "github.com/admacleod/aws/internal"
)

func TestSecureHeaders(t *testing.T) {
	testBody := "test"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, testBody); err != nil {
			t.Errorf("could not write testBody in testHandler: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "http://test.example.com", nil)
	w := httptest.NewRecorder()

	server.SecureHeaders(testHandler).ServeHTTP(w, req)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	if string(body) != testBody {
		t.Errorf("returned body is incorrect: expected=%s, got=%s", testBody, body)
	}

	for header, expected := range map[string]string{
		"Content-Security-Policy":   server.CSP,
		"Referrer-Policy":           "no-referrer",
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
	} {
		got := res.Header.Get(header)
		if expected != got {
			t.Errorf("%s header is incorrect: expected=%s, got=%s", header, expected, got)
		}
	}
}
