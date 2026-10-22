package cli

import (
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

func TestDisableKey(t *testing.T) {
	ks, primary, other := twoKeys(t, keyset.Enabled)
	out := []byte(mustRun(t, ks, "disable-key", "--key-id", idArg(other)))
	if st, _ := keyStatus(t, out, other); st != tinkpb.KeyStatusType_DISABLED {
		t.Errorf("status %v, want DISABLED", st)
	}
	if st, _ := keyStatus(t, out, primary); st != tinkpb.KeyStatusType_ENABLED {
		t.Errorf("primary status %v, want ENABLED", st)
	}
	// Disabling a disabled key is fine.
	mustRun(t, out, "disable-key", "--key-id", idArg(other))
}
