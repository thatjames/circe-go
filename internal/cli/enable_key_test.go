package cli

import (
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

func TestEnableKey(t *testing.T) {
	ks, _, other := twoKeys(t, keyset.Disabled)
	out := []byte(mustRun(t, ks, "enable-key", "--key-id", idArg(other)))
	if st, _ := keyStatus(t, out, other); st != tinkpb.KeyStatusType_ENABLED {
		t.Errorf("status %v, want ENABLED", st)
	}
	// Enabling an enabled key is fine.
	mustRun(t, out, "enable-key", "--key-id", idArg(other))
}
