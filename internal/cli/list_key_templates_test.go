package cli

import (
	"slices"
	"strings"
	"testing"

	"github.com/thatjames/circe-go/internal/templates"
)

func TestListKeyTemplates(t *testing.T) {
	out := mustRun(t, nil, "list-key-templates")
	got := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if !slices.Equal(got, templates.Names()) {
		t.Errorf("listed %d names, want the %d sorted template names", len(got), len(templates.Names()))
	}
	for _, n := range []string{"AES256_GCM", "ECDSA_P256", "ED25519"} {
		if !slices.Contains(got, n) {
			t.Errorf("%s missing", n)
		}
	}
}
