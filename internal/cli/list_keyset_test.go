package cli

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
)

func TestListKeyset(t *testing.T) {
	m := keyset.NewManagerFromHandle(newKeyset(t, "AES256_GCM"))
	id, err := m.Add(template(t, "CHACHA20_POLY1305_RAW"))
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Disable(id); err != nil {
		t.Fatal(err)
	}
	h, err := m.Handle()
	if err != nil {
		t.Fatal(err)
	}
	out := mustRun(t, jsonKeyset(t, h), "list-keyset")

	var got keysetInfo
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	want := h.KeysetInfo()
	if got.PrimaryKeyID != want.GetPrimaryKeyId() || len(got.KeyInfo) != 2 {
		t.Fatalf("got %+v", got)
	}
	for i, k := range want.GetKeyInfo() {
		g := got.KeyInfo[i]
		if g.KeyID != k.GetKeyId() || g.TypeURL != k.GetTypeUrl() ||
			g.Status != k.GetStatus().String() || g.OutputPrefixType != k.GetOutputPrefixType().String() {
			t.Errorf("key %d: got %+v, want %v", i, g, k)
		}
	}
	if got.KeyInfo[1].Status != "DISABLED" || got.KeyInfo[1].OutputPrefixType != "RAW" {
		t.Errorf("second key %+v, want DISABLED and RAW", got.KeyInfo[1])
	}

	// No key material: neither the base64 key value nor a field holding it.
	for _, k := range material(h).GetKey() {
		v := base64.StdEncoding.EncodeToString(k.GetKeyData().GetValue())
		if strings.Contains(out, v) || strings.Contains(out, "keyData") || strings.Contains(out, "value") {
			t.Error("list-keyset printed key material")
		}
	}
}

func TestListKeysetPublic(t *testing.T) {
	pub, err := newKeyset(t, "ECDSA_P256").Public()
	if err != nil {
		t.Fatal(err)
	}
	out := mustRun(t, jsonKeyset(t, pub), "list-keyset")
	if !strings.Contains(out, "EcdsaPublicKey") {
		t.Errorf("listing %s", out)
	}
}
