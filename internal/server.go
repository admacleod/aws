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

// Package internal provides a higher level wrapper around the standard library http.Server.
//
// This wrapping allows for the addition of helper functions that can be used to trivialise
// the setup of a secure web server.
//
// The default http.Server is close to being safe to expose directly to the internet but misses some important settings:
// Timeouts, TLS settings, and Response headers.
//
// Creating a safe, modern, web server with this package is as easy as:
//
//	srv := server.New(
//		server.Timeout(120 * time.Second),
//		server.TLS(server.ModerniseTLS(&tls.Config{})),
//		server.Handle(server.SecureHeaders(handler)),
//	)
//
// Additional configurations are also provided to simplify server creation in general.
package internal

import (
	"net/http"
	"time"
)

// Server defines a http server that allows for extension of the standard http.Server struct.
type Server struct {
	http.Server
}

// Option is a function that will apply some option to a Server object.
type Option func(*Server)

// New creates a new Server with the passed Options applied to it.
//
// If no options are passed then a default server implementation is returned.
func New(opts ...Option) *Server {
	srv := &Server{}

	for _, o := range opts {
		o(srv)
	}

	return srv
}

// Timeout creates a server.Option function that will set the passed time.Duration as the
// ReadTimeout, WriteTimeout, and IdleTimeout for the server.
func Timeout(timeout time.Duration) Option {
	return func(srv *Server) {
		srv.ReadTimeout = timeout
		srv.WriteTimeout = timeout
		srv.IdleTimeout = timeout
	}
}

// Handle creates a server.Option function that will set the passed http.Handler to the server as the Handler.
func Handle(handler http.Handler) Option {
	return func(srv *Server) {
		srv.Handler = handler
	}
}
