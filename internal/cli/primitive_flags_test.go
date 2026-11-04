package cli

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
)

// encodings are the values of --encoding.
var encodings = []string{"raw", "base64", "hex"}

// keysetFile writes h as a cleartext JSON keyset to a new file.
func keysetFile(t *testing.T, h *keyset.Handle) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ks.json")
	if err := os.WriteFile(path, jsonKeyset(t, h), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// dataFile writes b to a new file.
func dataFile(t *testing.T, name string, b []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// decodeAs decodes circe's output in an --encoding with the standard library.
func decodeAs(t *testing.T, enc, out string) []byte {
	t.Helper()
	var (
		b   []byte
		err error
	)
	switch enc {
	case "base64":
		b, err = base64.StdEncoding.DecodeString(strings.TrimSuffix(out, "\n"))
	case "hex":
		b, err = hex.DecodeString(strings.TrimSuffix(out, "\n"))
	default:
		b = []byte(out)
	}
	if err != nil {
		t.Fatalf("output is not %s: %v", enc, err)
	}
	return b
}

// encodeAs encodes b in an --encoding with the standard library.
func encodeAs(enc string, b []byte) []byte {
	switch enc {
	case "base64":
		return []byte(base64.StdEncoding.EncodeToString(b))
	case "hex":
		return []byte(hex.EncodeToString(b))
	}
	return b
}

func TestPrimitiveGroupPrintsHelp(t *testing.T) {
	r := run(t, nil, "aead")
	if r.code != 0 || !strings.Contains(r.stdout, "encrypt") || !strings.Contains(r.stdout, "decrypt") {
		t.Errorf("circe aead: exit %d, stdout %q", r.code, r.stdout)
	}
}

func TestPrimitiveUsageErrors(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	for _, args := range [][]string{
		{"aead", "encrypt"},
		{"aead", "seal"},
		{"aead", "encrypt", "--keyset", ks, "stray"},
		{"aead", "encrypt", "--keyset", ks, "--associated-data", "a", "--associated-data-file", ks},
	} {
		r := run(t, []byte("data"), args...)
		if r.code != 2 || !strings.Contains(r.stderr, "Usage:") || r.stdout != "" {
			t.Errorf("circe %v: exit %d, stdout %q, stderr %q; want 2 and the usage", args, r.code, r.stdout, r.stderr)
		}
	}
}

func TestWrongPrimitiveKeyset(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "HMAC_SHA256_128BITTAG"))
	r := run(t, []byte("data"), "aead", "encrypt", "--keyset", ks)
	want := "circe: " + ks + " is not an AEAD keyset (type.googleapis.com/google.crypto.tink.HmacKey)\n"
	if r.code != 1 || r.stderr != want || r.stdout != "" {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}

func TestFailedDecryptionExits3(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	out := filepath.Join(t.TempDir(), "plain")
	for _, args := range [][]string{
		{"aead", "decrypt", "--keyset", ks},
		{"aead", "decrypt", "--keyset", ks, "--out", out},
	} {
		r := run(t, []byte("not a ciphertext, long enough to parse"), args...)
		if r.code != 3 || r.stderr != "circe: decryption failed\n" || r.stdout != "" {
			t.Errorf("circe %v: exit %d, stdout %q, stderr %q; want 3 and decryption failed", args, r.code, r.stdout, r.stderr)
		}
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Errorf("a failed decryption created %s: %v", out, err)
	}
}

func TestPrimitiveOutRules(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	out := filepath.Join(t.TempDir(), "ct")
	if got := mustRun(t, []byte("secret"), "aead", "encrypt", "--keyset", ks, "--out", out); got != "" {
		t.Errorf("stdout %q, want nothing", got)
	}
	st, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Errorf("mode %v, want 0600", st.Mode().Perm())
	}
	before, _ := os.ReadFile(out)
	r := run(t, []byte("secret"), "aead", "encrypt", "--keyset", ks, "--out", out)
	if r.code != 1 || !strings.Contains(r.stderr, "exists") {
		t.Errorf("existing --out: exit %d, stderr %q; want 1 and an exists error", r.code, r.stderr)
	}
	if after, _ := os.ReadFile(out); string(after) != string(before) {
		t.Error("existing --out was overwritten")
	}
}

func TestPrimitiveReadsInFile(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	in := dataFile(t, "plain", []byte("from a file"))
	ct := mustRun(t, nil, "aead", "encrypt", "--keyset", ks, "--in", in)
	if got := mustRun(t, []byte(ct), "aead", "decrypt", "--keyset", ks); got != "from a file" {
		t.Errorf("decrypted %q", got)
	}
	if r := run(t, nil, "aead", "encrypt", "--keyset", ks, "--in", filepath.Join(t.TempDir(), "missing")); r.code != 1 {
		t.Errorf("missing --in: exit %d, want 1", r.code)
	}
}
