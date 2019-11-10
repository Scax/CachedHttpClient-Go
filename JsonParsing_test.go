package CachedHttpClient

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"testing"
)

func TestJsonX509Certificate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"rsa.PublicKey"},
		{"ecdsa.PublicKey"},
		{"dsa.PublicKey"},
		{"ed25519.PublicKey"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var publicKey, privateKey interface{}
			var err error
			var certificate x509.Certificate
			switch test.name {
			case "rsa.PublicKey":
				privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				publicKey = &privateKey.(*rsa.PrivateKey).PublicKey
			case "ecdsa.PublicKey":
				key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				publicKey = &key.PublicKey
			case "ed25519.PublicKey":
				publicKeyO, _, err := ed25519.GenerateKey(rand.Reader)
				publicKey = &publicKeyO
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			default:
				t.Skip("not tested")
				return
			}

			certificate.PublicKey = publicKey

			jsonX509Certificate := NewJsonX509Certificate(&certificate)

			jsonBytes, err := json.Marshal(jsonX509Certificate)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			var recreatedJsonCert JsonX509Certificate
			err = json.Unmarshal(jsonBytes, &recreatedJsonCert)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			toCertificate := recreatedJsonCert.ToCertificate()
			equal := certificate.Equal(toCertificate)

			if !equal {
				t.Error("not equal")
				t.FailNow()
			}

		})
	}
}

func TestJsonTlsConnectionState_ToConnectionState(t *testing.T) {

	state := &JsonTlsConnectionState{}
	state.ToConnectionState()
	state = nil
	state.ToConnectionState()

}
