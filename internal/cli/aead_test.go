package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/tink-crypto/tink-go/v2/aead"

	"github.com/thatjames/circe-go/internal/keysetio"
)

func TestAEADRoundTripEveryEncoding(t *testing.T) {
	h := newKeyset(t, "AES256_GCM")
	ks := keysetFile(t, h)
	a, err := aead.New(h)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("plaintext \x00\xff")
	for _, enc := range encodings {
		out := mustRun(t, plain, "aead", "encrypt", "--keyset", ks, "--encoding", enc, "--associated-data", "ctx")
		ct := decodeAs(t, enc, out)
		if got, err := a.Decrypt(ct, []byte("ctx")); err != nil || !bytes.Equal(got, plain) {
			t.Errorf("%s: tink-go decrypts circe's ciphertext to %q, %v", enc, got, err)
		}
		if got := mustRun(t, []byte(out), "aead", "decrypt", "--keyset", ks, "--encoding", enc, "--associated-data", "ctx"); got != string(plain) {
			t.Errorf("%s: circe decrypts to %q", enc, got)
		}
	}
	// circe decrypts tink-go's ciphertext, with the associated data from a file.
	ct, err := a.Encrypt(plain, []byte("from a file"))
	if err != nil {
		t.Fatal(err)
	}
	adFile := dataFile(t, "ad", []byte("from a file"))
	if got := mustRun(t, ct, "aead", "decrypt", "--keyset", ks, "--associated-data-file", adFile); got != string(plain) {
		t.Errorf("circe decrypts tink-go's ciphertext to %q", got)
	}
}

func TestAEADDecryptFailures(t *testing.T) {
	h := newKeyset(t, "AES128_GCM")
	ks := keysetFile(t, h)
	ct := []byte(mustRun(t, []byte("secret"), "aead", "encrypt", "--keyset", ks, "--associated-data", "right"))
	tampered := bytes.Clone(ct)
	tampered[len(tampered)-1] ^= 1
	other := keysetFile(t, newKeyset(t, "AES128_GCM"))
	for name, args := range map[string][]string{
		"tampered":  {"--keyset", ks, "--associated-data", "right"},
		"wrong AD":  {"--keyset", ks, "--associated-data", "wrong"},
		"no AD":     {"--keyset", ks},
		"wrong key": {"--keyset", other, "--associated-data", "right"},
	} {
		in := ct
		if name == "tampered" {
			in = tampered
		}
		r := run(t, in, append([]string{"aead", "decrypt"}, args...)...)
		if r.code != 3 || r.stderr != "circe: decryption failed\n" || r.stdout != "" {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 3", name, r.code, r.stdout, r.stderr)
		}
	}
}

func TestAEADBadEncodingExits1(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	r := run(t, []byte("not hex"), "aead", "decrypt", "--keyset", ks, "--encoding", "hex")
	if r.code != 1 || r.stderr != "circe: stdin is not valid hex\n" {
		t.Errorf("exit %d, stderr %q; want 1 and the input named", r.code, r.stderr)
	}
}

func TestAEADWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES256_SIV"))
	r := run(t, []byte("x"), "aead", "decrypt", "--keyset", ks)
	want := "circe: " + ks + " is not an AEAD keyset (type.googleapis.com/google.crypto.tink.AesSivKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}

func TestAEADWithKMSEncryptedKeyset(t *testing.T) {
	f := useFakeKMS(t)
	h := newKeyset(t, "AES256_GCM")
	var buf bytes.Buffer
	if err := keysetio.WriteEncrypted(&buf, h, keysetio.JSON, f.aeads["fake://a"]); err != nil {
		t.Fatal(err)
	}
	ks := filepath.Join(t.TempDir(), "ks.enc.json")
	if err := os.WriteFile(ks, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	ct := mustRun(t, []byte("secret"), "aead", "encrypt", "--keyset", ks, "--master-key-uri", "fake://a", "--credential", "cred.json")
	if f.creds["fake://a"] != "cred.json" {
		t.Errorf("credential %q, want cred.json", f.creds["fake://a"])
	}
	a, err := aead.New(h)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := a.Decrypt([]byte(ct), nil); err != nil || string(got) != "secret" {
		t.Errorf("tink-go decrypts to %q, %v", got, err)
	}
	if r := run(t, []byte(ct), "aead", "decrypt", "--keyset", ks); r.code != 1 {
		t.Errorf("encrypted keyset without --master-key-uri: exit %d, want 1", r.code)
	}
	if r := run(t, []byte(ct), "aead", "decrypt", "--keyset", ks, "--master-key-uri", "fake://b"); r.code != 1 {
		t.Errorf("keyset with the wrong KMS key: exit %d, want 1", r.code)
	}
}
