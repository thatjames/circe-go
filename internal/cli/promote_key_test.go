package cli

import (
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
)

func TestPromoteKey(t *testing.T) {
	ks, _, other := twoKeys(t, keyset.Enabled)
	out := mustRun(t, ks, "promote-key", "--key-id", idArg(other))
	if got := parseJSON(t, []byte(out)).KeysetInfo().GetPrimaryKeyId(); got != other {
		t.Errorf("primary %d, want %d", got, other)
	}
}

func TestPromoteDisabledKey(t *testing.T) {
	ks, _, other := twoKeys(t, keyset.Disabled)
	if r := run(t, ks, "promote-key", "--key-id", idArg(other)); r.code != 1 || r.stdout != "" {
		t.Errorf("exit %d, stdout %q; want 1 and nothing written", r.code, r.stdout)
	}
}
