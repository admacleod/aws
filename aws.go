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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	server "github.com/admacleod/aws/internal"

	"golang.org/x/crypto/acme/autocert"
)

const (
	usage = `Usage: %s [OPTION] HOSTNAME ...
Serve the current directory over HTTPS using ACME certificates for HOST(s).

`
	narg = `%[1]s: missing host operand
Try '%[1]s -h' for more information.
`
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	var (
		certDir string
	)
	flag.StringVar(&certDir, "c", "../certs", "certificate directory")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(flag.CommandLine.Output(), narg, os.Args[0])
		os.Exit(2)
	}

	// Configure TLS and certificate management
	mgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(flag.Args()...),
		Cache:      autocert.DirCache(certDir),
	}
	tlsCfg := mgr.TLSConfig()
	server.ModerniseTLS(tlsCfg)

	// Setup our handler
	mux := &http.ServeMux{}
	mw := server.ChainMiddleware(
		server.SecureHeaders,
		server.CombinedLogFormatLogger(os.Stdout),
	)
	handler := http.FileServer(http.Dir("."))
	mux.Handle("/", mw(handler))

	timeout := 10 * time.Second
	errLog := log.New(os.Stderr, "aws: ", log.LstdFlags)
	// We need two servers, one for HTTP redirect and the other for HTTPS
	srv := server.New(
		server.Timeout(timeout),
		server.ErrorLog(errLog),
		server.Handle(mgr.HTTPHandler(nil)),
	)
	srvTLS := server.New(
		server.Timeout(timeout),
		server.ErrorLog(errLog),
		server.Handle(mux),
		server.TLS(tlsCfg),
	)

	// Spool up and listen for errors
	e := make(chan error, 1)
	go func() {
		e <- srv.ListenAndServe()
	}()
	go func() {
		e <- srvTLS.ListenAndServeTLS("", "")
	}()
	err := <-e
	if err != nil {
		errLog.Fatalf("%v", err)
	}
}
