package cli

import (
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
)

func TestDeleteKey(t *testing.T) {
	ks, primary, other := twoKeys(t, keyset.Disabled)
	out := []byte(mustRun(t, ks, "delete-key", "--key-id", idArg(other)))
	if _, ok := keyStatus(t, out, other); ok {
		t.Error("deleted key still in the keyset")
	}
	if _, ok := keyStatus(t, out, primary); !ok {
		t.Error("primary gone")
	}
}
