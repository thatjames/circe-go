package cli

import (
	"testing"

	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

func TestAddKey(t *testing.T) {
	h := newKeyset(t, "AES256_GCM")
	primary := h.KeysetInfo().GetPrimaryKeyId()
	out := mustRun(t, jsonKeyset(t, h), "add-key", "--key-template", "AES128_GCM")
	info := parseJSON(t, []byte(out)).KeysetInfo()
	if len(info.GetKeyInfo()) != 2 {
		t.Fatalf("%d keys, want 2", len(info.GetKeyInfo()))
	}
	if info.GetPrimaryKeyId() != primary {
		t.Errorf("primary changed to %d", info.GetPrimaryKeyId())
	}
	k := info.GetKeyInfo()[1]
	if k.GetKeyId() == primary || k.GetStatus() != tinkpb.KeyStatusType_ENABLED ||
		k.GetTypeUrl() != "type.googleapis.com/google.crypto.tink.AesGcmKey" {
		t.Errorf("new key %v, want an enabled non-primary AES-GCM key", k)
	}
}

func TestAddKeyErrors(t *testing.T) {
	ks := jsonKeyset(t, newKeyset(t, "AES256_GCM"))
	if r := run(t, ks, "add-key"); r.code != 2 {
		t.Errorf("without --key-template: exit %d, want 2", r.code)
	}
	if r := run(t, ks, "add-key", "--key-template", "NOPE"); r.code != 1 || r.stdout != "" {
		t.Errorf("unknown template: exit %d, stdout %q", r.code, r.stdout)
	}
}
