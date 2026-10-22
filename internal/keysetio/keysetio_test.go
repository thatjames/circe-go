package keysetio

import (
	"bytes"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/signature"
)

// synthetic returns a keyset with two keys, the second disabled.
func synthetic(t *testing.T) *keyset.Handle {
	t.Helper()
	m := keyset.NewManager()
	id, err := m.Add(aead.AES256GCMKeyTemplate())
	if err != nil {
		t.Fatal(err)
	}
	if err := m.SetPrimary(id); err != nil {
		t.Fatal(err)
	}
	id2, err := m.Add(aead.ChaCha20Poly1305KeyTemplate())
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Disable(id2); err != nil {
		t.Fatal(err)
	}
	h, err := m.Handle()
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestRoundTrip(t *testing.T) {
	h := synthetic(t)
	for _, f := range []Format{JSON, Binary} {
		var buf bytes.Buffer
		if err := Write(&buf, h, f); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		got, err := Read(&buf, "test", f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		want := insecurecleartextkeyset.KeysetMaterial(h)
		if !proto.Equal(insecurecleartextkeyset.KeysetMaterial(got), want) {
			t.Errorf("%s: keyset changed in a round trip", f)
		}
	}
}

func TestReadsWhatTinkWrites(t *testing.T) {
	h := synthetic(t)
	var buf bytes.Buffer
	if err := insecurecleartextkeyset.Write(h, keyset.NewJSONWriter(&buf)); err != nil {
		t.Fatal(err)
	}
	got, err := Read(&buf, "test", JSON)
	if err != nil {
		t.Fatal(err)
	}
	if got.KeysetInfo().GetPrimaryKeyId() != h.KeysetInfo().GetPrimaryKeyId() {
		t.Error("primary key changed")
	}
}

func TestReadsPublicKeysets(t *testing.T) {
	priv, err := keyset.NewHandle(signature.ED25519KeyTemplate())
	if err != nil {
		t.Fatal(err)
	}
	pub, err := priv.Public()
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Write(&buf, pub, JSON); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(&buf, "test", JSON); err != nil {
		t.Fatal(err)
	}
}

func TestReadErrorsNameInputAndFormat(t *testing.T) {
	for _, tc := range []struct {
		in   string
		f    Format
		want []string
	}{
		{"not a keyset", JSON, []string{"in.json", "json"}},
		{"\x01\x02garbage", Binary, []string{"in.json", "binary"}},
		{"", JSON, []string{"in.json", "empty"}},
	} {
		_, err := Read(strings.NewReader(tc.in), "in.json", tc.f)
		if err == nil {
			t.Errorf("%q as %s: no error", tc.in, tc.f)
			continue
		}
		for _, w := range tc.want {
			if !strings.Contains(err.Error(), w) {
				t.Errorf("%q as %s: error %q does not mention %q", tc.in, tc.f, err, w)
			}
		}
	}
}

func TestFormatSet(t *testing.T) {
	var f Format
	for _, ok := range []string{"json", "binary"} {
		if err := f.Set(ok); err != nil || string(f) != ok {
			t.Errorf("Set(%q) = %v, f = %q", ok, err, f)
		}
	}
	for _, bad := range []string{"yaml", "JSON", ""} {
		if err := f.Set(bad); err == nil {
			t.Errorf("Set(%q): no error", bad)
		}
	}
}
