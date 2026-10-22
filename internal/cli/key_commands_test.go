package cli

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

// twoKeys returns a JSON keyset with an enabled primary and a second key
// with the given status, and both IDs.
func twoKeys(t *testing.T, second keyset.KeyStatus) (ks []byte, primary, other uint32) {
	t.Helper()
	m := keyset.NewManagerFromHandle(newKeyset(t, "AES256_GCM"))
	id, err := m.Add(template(t, "AES128_GCM"))
	if err != nil {
		t.Fatal(err)
	}
	if second == keyset.Disabled {
		if err := m.Disable(id); err != nil {
			t.Fatal(err)
		}
	}
	h, err := m.Handle()
	if err != nil {
		t.Fatal(err)
	}
	return jsonKeyset(t, h), h.KeysetInfo().GetPrimaryKeyId(), id
}

// keyStatus returns the status of key id in a JSON keyset, or false when absent.
func keyStatus(t *testing.T, ks []byte, id uint32) (tinkpb.KeyStatusType, bool) {
	t.Helper()
	for _, k := range parseJSON(t, ks).KeysetInfo().GetKeyInfo() {
		if k.GetKeyId() == id {
			return k.GetStatus(), true
		}
	}
	return 0, false
}

func idArg(id uint32) string { return strconv.FormatUint(uint64(id), 10) }

var keyIDCommands = []string{"promote-key", "enable-key", "disable-key", "delete-key"}

func TestKeyIDRequired(t *testing.T) {
	ks, _, _ := twoKeys(t, keyset.Enabled)
	for _, c := range keyIDCommands {
		if r := run(t, ks, c); r.code != 2 || !strings.Contains(r.stderr, "--key-id") {
			t.Errorf("%s without --key-id: exit %d, stderr %q", c, r.code, r.stderr)
		}
	}
}

func TestUnknownKeyIDWritesNothing(t *testing.T) {
	ks, _, _ := twoKeys(t, keyset.Enabled)
	for _, c := range keyIDCommands {
		out := filepath.Join(t.TempDir(), "out.json")
		r := run(t, ks, c, "--key-id", "1", "--out", out)
		if r.code != 1 || !strings.Contains(r.stderr, "no key with ID 1") {
			t.Errorf("%s: exit %d, stderr %q", c, r.code, r.stderr)
		}
		if _, err := os.Stat(out); !os.IsNotExist(err) {
			t.Errorf("%s wrote %s", c, out)
		}
	}
}

func TestPrimaryCannotBeDisabledOrDeleted(t *testing.T) {
	ks, primary, _ := twoKeys(t, keyset.Enabled)
	for _, c := range []string{"disable-key", "delete-key"} {
		r := run(t, ks, c, "--key-id", idArg(primary))
		if r.code != 1 || !strings.Contains(r.stderr, "primary") || r.stdout != "" {
			t.Errorf("%s on the primary: exit %d, stdout %q, stderr %q", c, r.code, r.stdout, r.stderr)
		}
		if strings.Contains(r.stderr, "keyset.Manager") {
			t.Errorf("%s: stderr %q keeps tink-go's prefix", c, r.stderr)
		}
	}
}
