package cli

import (
	"strings"
	"testing"

	"github.com/thatjames/circe-go/internal/templates"
)

// primitives are every primitive a command group works on.
var primitives = []primitive{
	aeadPrimitive, daeadPrimitive, streamingAEADPrimitive, macPrimitive, prfPrimitive,
	signaturePrimitive, hybridPrimitive, jwtMACPrimitive, jwtSignaturePrimitive,
}

// TestEveryTemplateHasOnePrimitive checks the key type table against the
// key templates circe offers, with the public key types of private ones.
func TestEveryTemplateHasOnePrimitive(t *testing.T) {
	for _, name := range templates.Names() {
		url := template(t, name).GetTypeUrl()
		urls := []string{url}
		if strings.HasSuffix(url, "PrivateKey") {
			urls = append(urls, strings.TrimSuffix(url, "PrivateKey")+"PublicKey")
		}
		for _, u := range urls {
			n := 0
			for _, p := range primitives {
				if p.has(u) {
					n++
				}
			}
			if n != 1 {
				t.Errorf("%s: %s is in %d primitives, want 1", name, u, n)
			}
		}
	}
}

func TestCheckKeysetKind(t *testing.T) {
	priv := newKeyset(t, "ED25519")
	pub := publicOf(t, priv)
	for _, c := range []struct {
		h    string
		kind keysetKind
		want string // empty: no error
	}{
		{"private", privateKeyset, ""},
		{"public", publicKeyset, ""},
		{"private", publicKeyset, "k.json is a private keyset; circe signature verify needs the public keyset: make it with circe create-public-keyset"},
		{"public", privateKeyset, "k.json is a public keyset; circe signature verify needs the private keyset"},
	} {
		h := priv
		if c.h == "public" {
			h = pub
		}
		err := checkKeyset(h, "k.json", "circe signature verify", signaturePrimitive, c.kind)
		switch {
		case c.want == "" && err != nil:
			t.Errorf("%s keyset, kind %d: %v", c.h, c.kind, err)
		case c.want != "" && (err == nil || err.Error() != c.want):
			t.Errorf("%s keyset, kind %d: error %v, want %q", c.h, c.kind, err, c.want)
		}
	}
}

func TestCheckKeysetPrimitiveBeforeKind(t *testing.T) {
	// A private hybrid keyset given to a signature command is the wrong
	// primitive, not the wrong kind.
	err := checkKeyset(newKeyset(t, "DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM"), "k.json", "circe signature verify", signaturePrimitive, publicKeyset)
	if err == nil || !strings.Contains(err.Error(), "is not a signature keyset (type.googleapis.com/google.crypto.tink.HpkePrivateKey)") {
		t.Errorf("error %v, want a wrong-primitive error", err)
	}
}
