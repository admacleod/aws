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
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	server "github.com/admacleod/aws/internal"
)

func TestLogger(t *testing.T) {
	testBody := "test"
	testMethod := "GET"
	testRequestURI := "http://test.example.com"
	testStatusCode := http.StatusCreated
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(testStatusCode)
		if _, err := io.WriteString(w, testBody); err != nil {
			t.Errorf("could not write testBody in testHandler: %v", err)
		}
	})

	req := httptest.NewRequest(testMethod, testRequestURI, nil)
	w := httptest.NewRecorder()

	var output bytes.Buffer
	middleware := server.CombinedLogFormatLogger(&output)
	start := time.Now()

	middleware(testHandler).ServeHTTP(w, req)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	if string(body) != testBody {
		t.Errorf("returned body is incorrect: expected=%s, got=%s", testBody, body)
	}

	if testStatusCode != res.StatusCode {
		t.Errorf("returned status code is incorrect: expected=%d, got=%d", testStatusCode, res.StatusCode)
	}

	expected := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
		strings.Split(req.RemoteAddr, ":")[0],
		start.Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.RequestURI,
		req.Proto,
		res.StatusCode,
		len(body),
		req.Referer(),
		req.UserAgent(),
	)
	got := output.String()

	if got != expected {
		t.Errorf("error with log output:\nexpect=\"%s\"\nactual=\"%s\"", expected, got)
	}
}

func TestServerLoggerOption(t *testing.T) {
	testLogger := log.New(os.Stdout, "test: ", log.LUTC)
	testSrv := server.New(
		server.ErrorLog(testLogger),
	)

	if !reflect.DeepEqual(testLogger, testSrv.ErrorLog) {
		t.Errorf("incorrect server error log: expected=%v, got=%v", testLogger, testSrv.ErrorLog)
	}
}
