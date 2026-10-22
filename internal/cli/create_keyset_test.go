package cli

import (
	"strings"
	"testing"

	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

func TestCreateKeyset(t *testing.T) {
	for _, name := range []string{"AES256_GCM", "AES256_SIV", "HMAC_SHA256_256BITTAG", "ED25519", "JWT_ES256",
		"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM"} {
		t.Run(name, func(t *testing.T) {
			h := parseJSON(t, []byte(mustRun(t, nil, "create-keyset", "--key-template", name)))
			info := h.KeysetInfo()
			if len(info.GetKeyInfo()) != 1 {
				t.Fatalf("%d keys, want 1", len(info.GetKeyInfo()))
			}
			k := info.GetKeyInfo()[0]
			if k.GetKeyId() != info.GetPrimaryKeyId() || k.GetStatus() != tinkpb.KeyStatusType_ENABLED {
				t.Errorf("key %d status %v, primary %d; want one enabled primary", k.GetKeyId(), k.GetStatus(), info.GetPrimaryKeyId())
			}
		})
	}
}

func TestCreateKeysetBinary(t *testing.T) {
	bin := mustRun(t, nil, "create-keyset", "--key-template", "AES128_GCM", "--out-format", "binary")
	js := mustRun(t, []byte(bin), "convert-keyset", "--in-format", "binary")
	parseJSON(t, []byte(js))
}

func TestCreateKeysetUnknownTemplate(t *testing.T) {
	r := run(t, nil, "create-keyset", "--key-template", "NOPE")
	if r.code != 1 || !strings.Contains(r.stderr, "list-key-templates") || r.stdout != "" {
		t.Errorf("exit %d, stdout %q, stderr %q", r.code, r.stdout, r.stderr)
	}
}
