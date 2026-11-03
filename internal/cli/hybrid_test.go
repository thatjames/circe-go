package cli

import (
	"bytes"
	"testing"

	"github.com/tink-crypto/tink-go/v2/hybrid"
)

// hybridTemplates are one template of each hybrid family and AEAD.
var hybridTemplates = []string{
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_GCM",
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_CTR_HMAC_SHA256_RAW",
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM",
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_CHACHA20_POLY1305_RAW",
	"DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_256_GCM",
}

func TestHybridEveryFamilyAndEncoding(t *testing.T) {
	plain := []byte("a secret for the service")
	for _, name := range hybridTemplates {
		priv := newKeyset(t, name)
		pub := publicOf(t, priv)
		privFile, pubFile := keysetFile(t, priv), keysetFile(t, pub)
		dec, err := hybrid.NewHybridDecrypt(priv)
		if err != nil {
			t.Fatal(err)
		}
		enc, err := hybrid.NewHybridEncrypt(pub)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range encodings {
			out := mustRun(t, plain, "hybrid", "encrypt", "--keyset", pubFile, "--encoding", e, "--context-info", "ctx")
			if got, err := dec.Decrypt(decodeAs(t, e, out), []byte("ctx")); err != nil || !bytes.Equal(got, plain) {
				t.Errorf("%s, %s: tink-go decrypts circe's ciphertext to %q, %v", name, e, got, err)
			}
			if got := mustRun(t, []byte(out), "hybrid", "decrypt", "--keyset", privFile, "--encoding", e, "--context-info", "ctx"); got != string(plain) {
				t.Errorf("%s, %s: circe decrypts to %q", name, e, got)
			}
		}
		// circe decrypts tink-go's ciphertext, with the context info from a file.
		ct, err := enc.Encrypt(plain, []byte("file ctx"))
		if err != nil {
			t.Fatal(err)
		}
		ctxFile := dataFile(t, "ctx", []byte("file ctx"))
		if got := mustRun(t, ct, "hybrid", "decrypt", "--keyset", privFile, "--context-info-file", ctxFile); got != string(plain) {
			t.Errorf("%s: circe decrypts tink-go's ciphertext to %q", name, got)
		}
		r := run(t, ct, "hybrid", "decrypt", "--keyset", privFile, "--context-info", "wrong")
		if r.code != 3 || r.stderr != "circe: decryption failed\n" || r.stdout != "" {
			t.Errorf("%s: wrong context info: exit %d, stderr %q; want 3", name, r.code, r.stderr)
		}
	}
}

func TestHybridWrongKindOfKeyset(t *testing.T) {
	priv := newKeyset(t, "DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM")
	privFile, pubFile := keysetFile(t, priv), keysetFile(t, publicOf(t, priv))
	r := run(t, []byte("x"), "hybrid", "encrypt", "--keyset", privFile)
	want := "circe: " + privFile + " is a private keyset; circe hybrid encrypt needs the public keyset: make it with circe create-public-keyset\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("encrypt with the private keyset: exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
	r = run(t, []byte("x"), "hybrid", "decrypt", "--keyset", pubFile)
	want = "circe: " + pubFile + " is a public keyset; circe hybrid decrypt needs the private keyset\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("decrypt with the public keyset: exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}

func TestHybridContextInfoFlagsExclusive(t *testing.T) {
	priv := newKeyset(t, "DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM")
	pubFile := keysetFile(t, publicOf(t, priv))
	r := run(t, []byte("x"), "hybrid", "encrypt", "--keyset", pubFile, "--context-info", "a", "--context-info-file", pubFile)
	if r.code != 2 {
		t.Errorf("both context info flags: exit %d, want 2", r.code)
	}
}
