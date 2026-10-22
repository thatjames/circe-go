package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/tink"
)

// fakeKMS replaces the KMS with local AEADs for the URIs fake://a and
// fake://b, records the credential each URI was asked with, and fails
// for fake://down.
type fakeKMS struct {
	aeads map[string]tink.AEAD
	creds map[string]string
}

func useFakeKMS(t *testing.T) *fakeKMS {
	t.Helper()
	f := &fakeKMS{aeads: map[string]tink.AEAD{}, creds: map[string]string{}}
	for _, uri := range []string{"fake://a", "fake://b"} {
		a, err := aead.New(newKeyset(t, "AES256_GCM"))
		if err != nil {
			t.Fatal(err)
		}
		f.aeads[uri] = a
	}
	old := kmsAEAD
	kmsAEAD = func(_ context.Context, uri, credential string) (tink.AEAD, error) {
		f.creds[uri] = credential
		if a, ok := f.aeads[uri]; ok {
			return a, nil
		}
		return nil, errors.New("Fake KMS key " + uri + ": permission denied")
	}
	t.Cleanup(func() { kmsAEAD = old })
	return f
}

// decrypt reads an encrypted JSON keyset with tink-go alone.
func (f *fakeKMS) decrypt(t *testing.T, uri string, b []byte) *keyset.Handle {
	t.Helper()
	h, err := keyset.Read(keyset.NewJSONReader(bytes.NewReader(b)), f.aeads[uri])
	if err != nil {
		t.Fatalf("keyset.Read with %s: %v", uri, err)
	}
	return h
}

func TestEncryptedThroughEveryCommand(t *testing.T) {
	f := useFakeKMS(t)
	a := []string{"--master-key-uri", "fake://a"}
	step := func(in []byte, args ...string) []byte {
		t.Helper()
		out := []byte(mustRun(t, in, append(args, a...)...))
		if bytes.Contains(out, []byte("keyData")) {
			t.Fatalf("%s wrote cleartext key data", args[0])
		}
		return out
	}

	ks := step(nil, "create-keyset", "--key-template", "ED25519")
	first := f.decrypt(t, "fake://a", ks).KeysetInfo().GetPrimaryKeyId()

	ks = step(ks, "add-key", "--key-template", "ED25519")
	info := f.decrypt(t, "fake://a", ks).KeysetInfo()
	second := info.GetKeyInfo()[1].GetKeyId()

	ks = step(ks, "promote-key", "--key-id", idArg(second))
	ks = step(ks, "disable-key", "--key-id", idArg(first))
	ks = step(ks, "enable-key", "--key-id", idArg(first))
	ks = step(ks, "disable-key", "--key-id", idArg(first))
	ks = step(ks, "delete-key", "--key-id", idArg(first))
	h := f.decrypt(t, "fake://a", ks)
	if h.KeysetInfo().GetPrimaryKeyId() != second || len(h.KeysetInfo().GetKeyInfo()) != 1 {
		t.Fatalf("keyset after the lifecycle: %v", h.KeysetInfo())
	}

	if out := mustRun(t, ks, append([]string{"list-keyset"}, a...)...); !strings.Contains(out, idArg(second)) {
		t.Errorf("list-keyset: %s", out)
	}

	// The public keyset is written in cleartext.
	pub := parseJSON(t, []byte(mustRun(t, ks, append([]string{"create-public-keyset"}, a...)...)))
	if pub.KeysetInfo().GetPrimaryKeyId() != second {
		t.Error("public keyset lost the primary")
	}
}

func TestConvertBetweenKMSKeys(t *testing.T) {
	f := useFakeKMS(t)
	h := newKeyset(t, "AES256_GCM")
	clear := jsonKeyset(t, h)

	// Cleartext to fake://a, a to b, b to cleartext.
	encA := []byte(mustRun(t, clear, "convert-keyset", "--new-master-key-uri", "fake://a"))
	encB := []byte(mustRun(t, encA, "convert-keyset", "--master-key-uri", "fake://a", "--new-master-key-uri", "fake://b",
		"--credential", "cred-a", "--new-credential", "cred-b"))
	if f.creds["fake://a"] != "cred-a" || f.creds["fake://b"] != "cred-b" {
		t.Errorf("credentials %v, want one per side", f.creds)
	}
	back := []byte(mustRun(t, encB, "convert-keyset", "--master-key-uri", "fake://b"))

	if !proto.Equal(material(f.decrypt(t, "fake://b", encB)), material(h)) {
		t.Error("a to b changed the keyset")
	}
	if !proto.Equal(material(parseJSON(t, back)), material(h)) {
		t.Error("b to cleartext changed the keyset")
	}

	// Binary encrypted output too.
	bin := mustRun(t, clear, "convert-keyset", "--new-master-key-uri", "fake://a", "--out-format", "binary")
	if _, err := keyset.Read(keyset.NewBinaryReader(strings.NewReader(bin)), f.aeads["fake://a"]); err != nil {
		t.Errorf("binary encrypted keyset: %v", err)
	}
}

func TestDecryptFailureWritesNothing(t *testing.T) {
	useFakeKMS(t)
	encA := []byte(mustRun(t, nil, "create-keyset", "--key-template", "AES256_GCM", "--master-key-uri", "fake://a"))
	in := filepath.Join(t.TempDir(), "in.json")
	if err := os.WriteFile(in, encA, 0o600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "out.json")
	r := run(t, nil, "convert-keyset", "--in", in, "--master-key-uri", "fake://b", "--out", out)
	if r.code != 1 || !strings.Contains(r.stderr, in) || !strings.Contains(r.stderr, "fake://b") {
		t.Errorf("wrong key: exit %d, stderr %q", r.code, r.stderr)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("output written after a failed decrypt")
	}

	// An encrypted keyset read as cleartext, and the reverse, fail too.
	if r := run(t, encA, "list-keyset"); r.code != 1 {
		t.Errorf("encrypted input without --master-key-uri: exit %d", r.code)
	}
	if r := run(t, jsonKeyset(t, newKeyset(t, "AES256_GCM")), "list-keyset", "--master-key-uri", "fake://a"); r.code != 1 {
		t.Errorf("cleartext input with --master-key-uri: exit %d", r.code)
	}
}

func TestKMSErrorExits1(t *testing.T) {
	useFakeKMS(t)
	r := run(t, nil, "create-keyset", "--key-template", "AES256_GCM", "--master-key-uri", "fake://down")
	if r.code != 1 || !strings.Contains(r.stderr, "fake://down") || r.stdout != "" {
		t.Errorf("exit %d, stdout %q, stderr %q", r.code, r.stdout, r.stderr)
	}
}
