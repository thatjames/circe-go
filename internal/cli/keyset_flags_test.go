package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestOutWritesANewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ks.json")
	if out := mustRun(t, nil, "create-keyset", "--key-template", "AES256_GCM", "--out", path); out != "" {
		t.Errorf("stdout %q, want nothing", out)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parseJSON(t, b)
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Errorf("mode %v, want 0600", st.Mode().Perm())
	}
}

func TestOutNeverOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ks.json")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := run(t, nil, "create-keyset", "--key-template", "AES256_GCM", "--out", path)
	if r.code != 1 || !strings.Contains(r.stderr, "exists") {
		t.Errorf("exit %d, stderr %q; want 1 and an exists error", r.code, r.stderr)
	}
	if b, _ := os.ReadFile(path); string(b) != "keep me" {
		t.Errorf("file changed to %q", b)
	}
}

func TestFailingCommandWritesNothing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ks.json")
	r := run(t, []byte("not a keyset"), "convert-keyset", "--out", path)
	if r.code != 1 {
		t.Fatalf("exit %d, want 1", r.code)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("%s was created: %v", path, err)
	}
}

func TestInReadsAFile(t *testing.T) {
	h := newKeyset(t, "AES128_GCM")
	path := filepath.Join(t.TempDir(), "ks.json")
	if err := os.WriteFile(path, jsonKeyset(t, h), 0o600); err != nil {
		t.Fatal(err)
	}
	got := parseJSON(t, []byte(mustRun(t, nil, "convert-keyset", "--in", path)))
	if !proto.Equal(material(got), material(h)) {
		t.Error("keyset changed")
	}
}

func TestInErrorNamesThePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := run(t, nil, "list-keyset", "--in", path)
	if r.code != 1 || !strings.Contains(r.stderr, path) || !strings.Contains(r.stderr, "json") {
		t.Errorf("exit %d, stderr %q; want 1 naming %s and json", r.code, r.stderr, path)
	}
	r = run(t, nil, "list-keyset", "--in", filepath.Join(t.TempDir(), "missing.json"))
	if r.code != 1 {
		t.Errorf("missing --in file: exit %d, want 1", r.code)
	}
}
