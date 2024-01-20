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

package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

func (lrw *loggerResponseWriter) WriteHeader(code int) {
	lrw.ResponseWriter.WriteHeader(code)
	lrw.Status = code
}

func (lrw *loggerResponseWriter) Write(bb []byte) (int, error) {
	length, err := lrw.ResponseWriter.Write(bb)
	lrw.ContentLength += length

	return length, err
}

// CombinedLogFormatLogger is a middleware generator function that will write an Apache Combined Log Format
// to the passed output Writer for all requests to the wrapped handler.
//
// The definition of the Combined Log Format can be found at: https://httpd.apache.org/docs/2.4/logs.html#combined
func CombinedLogFormatLogger(output io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := loggerResponseWriter{w, 200, 0}
			next.ServeHTTP(&lrw, r)
			fmt.Fprintf(output, "%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
				strings.Split(r.RemoteAddr, ":")[0], // Remove potential port number from remote address
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.RequestURI,
				r.Proto,
				lrw.Status,
				lrw.ContentLength,
				r.Referer(),
				r.UserAgent(),
			)
		})
	}
}

// ErrorLog creates a server.Option function that will apply the passed log.Logger to the server as the ErrorLog.
func ErrorLog(logger *log.Logger) Option {
	return func(srv *Server) {
		srv.ErrorLog = logger
	}
}
