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
	"crypto/tls"
	"reflect"
	"testing"

	server "github.com/admacleod/aws/internal"
)

func TestTLSModerniseConfig(t *testing.T) {
	testConfig := &tls.Config{}

	server.ModerniseTLS(testConfig)

	if !reflect.DeepEqual([]tls.CurveID{
		tls.X25519,
		tls.CurveP256,
		tls.CurveP384,
	}, testConfig.CurvePreferences) {
		t.Errorf("incorrect curve preferences: expected=%v, got=%v", []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		}, testConfig.CurvePreferences)
	}

	if testConfig.MinVersion != uint16(tls.VersionTLS12) {
		t.Errorf("incorrect minimum tls setting: expected=%v, got=%v", uint16(tls.VersionTLS12), testConfig.MinVersion)
	}

	if !reflect.DeepEqual([]uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}, testConfig.CipherSuites) {
		t.Errorf("incorrect cipher suite preferences: expected=%v, got=%v", []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}, testConfig.CipherSuites)
	}
}

func TestTLSServerOption(t *testing.T) {
	testConfig := &tls.Config{}
	server.ModerniseTLS(testConfig)
	testSrv := server.New(
		server.TLS(testConfig),
	)

	if !reflect.DeepEqual(testConfig, testSrv.TLSConfig) {
		t.Errorf("incorrect tls config: expected=%v, got=%v", testConfig, testSrv.TLSConfig)
	}
}
