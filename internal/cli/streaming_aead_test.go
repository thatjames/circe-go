package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/tink-crypto/tink-go/v2/streamingaead"
)

// bigData is larger than several 4KB segments.
func bigData() []byte {
	return bytes.Repeat([]byte("0123456789abcdef"), 1500)
}

func TestStreamingAEADRoundTripAndInterop(t *testing.T) {
	h := newKeyset(t, "AES128_GCM_HKDF_4KB")
	ks := keysetFile(t, h)
	s, err := streamingaead.New(h)
	if err != nil {
		t.Fatal(err)
	}
	plain := bigData()
	dir := t.TempDir()
	in, ct, out := filepath.Join(dir, "plain"), filepath.Join(dir, "ct"), filepath.Join(dir, "out")
	if err := os.WriteFile(in, plain, 0o600); err != nil {
		t.Fatal(err)
	}
	mustRun(t, nil, "streaming-aead", "encrypt", "--keyset", ks, "--in", in, "--out", ct, "--associated-data", "ad")
	mustRun(t, nil, "streaming-aead", "decrypt", "--keyset", ks, "--in", ct, "--out", out, "--associated-data", "ad")
	if got, _ := os.ReadFile(out); !bytes.Equal(got, plain) {
		t.Errorf("round trip changed the data (%d bytes, want %d)", len(got), len(plain))
	}
	if st, err := os.Stat(out); err != nil || st.Mode().Perm() != 0o600 {
		t.Errorf("--out mode: %v, %v; want 0600", st.Mode().Perm(), err)
	}

	// tink-go decrypts circe's stream.
	b, _ := os.ReadFile(ct)
	r, err := s.NewDecryptingReader(bytes.NewReader(b), []byte("ad"))
	if err != nil {
		t.Fatal(err)
	}
	if got, err := io.ReadAll(r); err != nil || !bytes.Equal(got, plain) {
		t.Errorf("tink-go decrypts circe's stream: %d bytes, %v", len(got), err)
	}

	// circe decrypts tink-go's stream, on stdin and stdout.
	var buf bytes.Buffer
	w, err := s.NewEncryptingWriter(&buf, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(plain); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if got := mustRun(t, buf.Bytes(), "streaming-aead", "decrypt", "--keyset", ks); got != string(plain) {
		t.Errorf("circe decrypts tink-go's stream to %d bytes", len(got))
	}
}

func TestStreamingAEADFailedDecryptWritesNothing(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM_HKDF_4KB"))
	ct := []byte(mustRun(t, bigData(), "streaming-aead", "encrypt", "--keyset", ks))
	tampered := bytes.Clone(ct)
	tampered[len(tampered)-10] ^= 1
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	for name, c := range map[string]struct {
		in   []byte
		args []string
	}{
		"tampered": {tampered, nil},
		"wrong AD": {ct, []string{"--associated-data", "x"}},
	} {
		r := run(t, c.in, append([]string{"streaming-aead", "decrypt", "--keyset", ks, "--out", out}, c.args...)...)
		if r.code != 3 || r.stderr != "circe: decryption failed\n" {
			t.Errorf("%s: exit %d, stderr %q; want 3", name, r.code, r.stderr)
		}
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("a failed decrypt left %d files in the --out directory", len(entries))
	}
}

func TestStreamingAEADNeverOverwrites(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM_HKDF_4KB"))
	out := dataFile(t, "out", []byte("keep me"))
	r := run(t, []byte("x"), "streaming-aead", "encrypt", "--keyset", ks, "--out", out)
	if r.code != 1 || r.stderr != "circe: "+out+" exists; circe never overwrites a file\n" {
		t.Errorf("exit %d, stderr %q; want 1 and an exists error", r.code, r.stderr)
	}
	if b, _ := os.ReadFile(out); string(b) != "keep me" {
		t.Errorf("file changed to %q", b)
	}
	if entries, _ := os.ReadDir(filepath.Dir(out)); len(entries) != 1 {
		t.Errorf("%d files in the --out directory, want 1", len(entries))
	}
}

func TestStreamingAEADHasNoEncoding(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM_HKDF_4KB"))
	if r := run(t, []byte("x"), "streaming-aead", "encrypt", "--keyset", ks, "--encoding", "hex"); r.code != 2 {
		t.Errorf("--encoding: exit %d, want 2", r.code)
	}
}

func TestStreamingAEADWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	r := run(t, []byte("x"), "streaming-aead", "encrypt", "--keyset", ks)
	want := "circe: " + ks + " is not a streaming AEAD keyset (type.googleapis.com/google.crypto.tink.AesGcmKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}
