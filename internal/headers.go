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

import "net/http"

// CSP defines the Content-Security-Policy applied by the SecureHeaders function.
// The policy is very restrictive. Currently only allowing self-hosted CSS, images,
// and PDF documents. JavaScript, forms, and iframes are disallowed.
const CSP = "default-src 'none';" +
	"style-src 'self';" +
	"img-src 'self';" +
	"object-src 'self';" +
	"base-uri 'none';" +
	"form-action 'none';" +
	"frame-ancestors 'none';" +
	"plugin-types application/pdf"

// SecureHeaders is a http middleware for adding security headers to server responses.
// Applying the middleware will add the following header values, inspired by
// https://securityheaders.com, to responses from the wrapped handler.
//
//	Content-Security-Policy: [see CSP constant]
//	Referrer-Policy: no-referrer
//	Strict-Transport-Security: max-age=63072000; includeSubDomains
//	X-Content-Type-Options: nosniff
//	X-Frame-Options: DENY
//	X-XSS-Protection: 1; mode=block
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		w.Header().Set("Content-Security-Policy", CSP)
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
