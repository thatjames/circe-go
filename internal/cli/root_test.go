package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"

	"github.com/thatjames/circe-go/internal/templates"
)

// result is what one circe run produced.
type result struct {
	code           int
	stdout, stderr string
}

// run runs circe in-process with stdin and args.
func run(t *testing.T, stdin []byte, args ...string) result {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run(args, bytes.NewReader(stdin), &out, &errb)
	return result{code, out.String(), errb.String()}
}

// mustRun runs circe and fails the test unless it exits 0.
func mustRun(t *testing.T, stdin []byte, args ...string) string {
	t.Helper()
	r := run(t, stdin, args...)
	if r.code != 0 {
		t.Fatalf("circe %s: exit %d: %s", strings.Join(args, " "), r.code, r.stderr)
	}
	return r.stdout
}

// template looks up a key template by name.
func template(t *testing.T, name string) *tinkpb.KeyTemplate {
	t.Helper()
	kt, err := templates.Lookup(name)
	if err != nil {
		t.Fatal(err)
	}
	return kt
}

// newKeyset makes a synthetic keyset from a template name.
func newKeyset(t *testing.T, name string) *keyset.Handle {
	t.Helper()
	h, err := keyset.NewHandle(template(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// jsonKeyset serializes h the way Tink writes cleartext JSON keysets.
func jsonKeyset(t *testing.T, h *keyset.Handle) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := insecurecleartextkeyset.Write(h, keyset.NewJSONWriter(&buf)); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// parseJSON reads a cleartext JSON keyset with tink-go alone.
func parseJSON(t *testing.T, b []byte) *keyset.Handle {
	t.Helper()
	h, err := insecurecleartextkeyset.Read(keyset.NewJSONReader(bytes.NewReader(b)))
	if err != nil {
		t.Fatalf("tink-go cannot read the keyset: %v", err)
	}
	return h
}

func material(h *keyset.Handle) *tinkpb.Keyset { return insecurecleartextkeyset.KeysetMaterial(h) }

func TestHelp(t *testing.T) {
	for _, args := range [][]string{nil, {"help"}} {
		r := run(t, nil, args...)
		if r.code != 0 || !strings.Contains(r.stdout, "create-keyset") || r.stderr != "" {
			t.Errorf("circe %v: exit %d, stdout %q, stderr %q", args, r.code, r.stdout, r.stderr)
		}
	}
}

func TestVersion(t *testing.T) {
	old := version
	defer func() { version = old }()
	version = "v1.2.3"
	if got := mustRun(t, nil, "--version"); got != "circe v1.2.3\n" {
		t.Errorf("--version = %q", got)
	}
	version = ""
	if got := mustRun(t, nil, "--version"); !strings.HasPrefix(got, "circe ") {
		t.Errorf("--version = %q", got)
	}
}

func TestUsageErrorsExit2(t *testing.T) {
	for _, args := range [][]string{
		{"no-such-command"},
		{"--no-such-flag"},
		{"list-key-templates", "stray"},
		{"list-keyset", "--in-format", "yaml"},
		{"create-keyset"},
	} {
		r := run(t, nil, args...)
		if r.code != 2 {
			t.Errorf("circe %v: exit %d, want 2", args, r.code)
		}
		if !strings.HasPrefix(r.stderr, "circe: ") || !strings.Contains(r.stderr, "Usage:") {
			t.Errorf("circe %v: stderr %q, want the error and the usage", args, r.stderr)
		}
		if r.stdout != "" {
			t.Errorf("circe %v: stdout %q, want nothing", args, r.stdout)
		}
	}
}

func TestErrorsExit1OnOneLine(t *testing.T) {
	r := run(t, []byte("not a keyset"), "list-keyset")
	if r.code != 1 {
		t.Fatalf("exit %d, want 1", r.code)
	}
	if !strings.HasPrefix(r.stderr, "circe: ") || strings.Count(r.stderr, "\n") != 1 {
		t.Errorf("stderr %q, want one circe: line", r.stderr)
	}
}
