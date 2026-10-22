package keysetio

import (
	"bytes"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/tink"
)

// fakeKMS is an AEAD from a fresh local keyset, standing in for a KMS key.
func fakeKMS(t *testing.T) tink.AEAD {
	t.Helper()
	h, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		t.Fatal(err)
	}
	a, err := aead.New(h)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestEncryptedRoundTrip(t *testing.T) {
	h := synthetic(t)
	kek := fakeKMS(t)
	for _, f := range []Format{JSON, Binary} {
		var buf bytes.Buffer
		if err := WriteEncrypted(&buf, h, f, kek); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		// tink-go alone reads what circe writes.
		got, err := keyset.Read(f.Reader(bytes.NewReader(buf.Bytes())), kek)
		if err != nil {
			t.Fatalf("%s: keyset.Read: %v", f, err)
		}
		if !proto.Equal(insecurecleartextkeyset.KeysetMaterial(got), insecurecleartextkeyset.KeysetMaterial(h)) {
			t.Errorf("%s: keyset changed", f)
		}
		got, err = ReadEncrypted(&buf, "test", f, kek, "fake://kek")
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if !proto.Equal(insecurecleartextkeyset.KeysetMaterial(got), insecurecleartextkeyset.KeysetMaterial(h)) {
			t.Errorf("%s: keyset changed", f)
		}
	}
}

func TestEncryptedHoldsNoCleartext(t *testing.T) {
	h := synthetic(t)
	var buf bytes.Buffer
	if err := WriteEncrypted(&buf, h, JSON, fakeKMS(t)); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "keyData") {
		t.Error("encrypted keyset holds cleartext key data")
	}
}

func TestDecryptErrorNamesInputAndURI(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteEncrypted(&buf, synthetic(t), JSON, fakeKMS(t)); err != nil {
		t.Fatal(err)
	}
	_, err := ReadEncrypted(&buf, "in.json", JSON, fakeKMS(t), "fake://other")
	if err == nil || !strings.Contains(err.Error(), "in.json") || !strings.Contains(err.Error(), "fake://other") {
		t.Errorf("error %v, want one naming in.json and fake://other", err)
	}
}
