package cli

import (
	"bytes"
	"testing"

	"github.com/tink-crypto/tink-go/v2/daead"
)

func TestDAEADRoundTripEveryEncoding(t *testing.T) {
	h := newKeyset(t, "AES256_SIV")
	ks := keysetFile(t, h)
	d, err := daead.New(h)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("column value")
	want, err := d.EncryptDeterministically(plain, []byte("ad"))
	if err != nil {
		t.Fatal(err)
	}
	for _, enc := range encodings {
		out := mustRun(t, plain, "daead", "encrypt", "--keyset", ks, "--encoding", enc, "--associated-data", "ad")
		if got := decodeAs(t, enc, out); !bytes.Equal(got, want) {
			t.Errorf("%s: circe's ciphertext differs from tink-go's", enc)
		}
		if again := mustRun(t, plain, "daead", "encrypt", "--keyset", ks, "--encoding", enc, "--associated-data", "ad"); again != out {
			t.Errorf("%s: encryption is not deterministic", enc)
		}
		if got := mustRun(t, encodeAs(enc, want), "daead", "decrypt", "--keyset", ks, "--encoding", enc, "--associated-data", "ad"); got != string(plain) {
			t.Errorf("%s: circe decrypts tink-go's ciphertext to %q", enc, got)
		}
	}
}

func TestDAEADDecryptFailures(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES256_SIV"))
	ct := []byte(mustRun(t, []byte("secret"), "daead", "encrypt", "--keyset", ks))
	tampered := bytes.Clone(ct)
	tampered[len(tampered)-1] ^= 1
	for name, c := range map[string]struct {
		in   []byte
		args []string
	}{
		"tampered": {tampered, nil},
		"wrong AD": {ct, []string{"--associated-data", "x"}},
	} {
		r := run(t, c.in, append([]string{"daead", "decrypt", "--keyset", ks}, c.args...)...)
		if r.code != 3 || r.stderr != "circe: decryption failed\n" || r.stdout != "" {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 3", name, r.code, r.stdout, r.stderr)
		}
	}
}

func TestDAEADWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	r := run(t, []byte("x"), "daead", "encrypt", "--keyset", ks)
	want := "circe: " + ks + " is not a deterministic AEAD keyset (type.googleapis.com/google.crypto.tink.AesGcmKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}
