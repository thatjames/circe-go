package cli

import (
	"bytes"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func TestConvertKeysetJSONBinaryJSON(t *testing.T) {
	h := newKeyset(t, "ECDSA_P256")
	bin := mustRun(t, jsonKeyset(t, h), "convert-keyset", "--out-format", "binary")

	fromBin, err := insecurecleartextkeyset.Read(keyset.NewBinaryReader(bytes.NewReader([]byte(bin))))
	if err != nil {
		t.Fatalf("tink-go cannot read the binary keyset: %v", err)
	}
	if !proto.Equal(material(fromBin), material(h)) {
		t.Error("JSON to binary changed the keyset")
	}

	js := mustRun(t, []byte(bin), "convert-keyset", "--in-format", "binary", "--out-format", "json")
	if !proto.Equal(material(parseJSON(t, []byte(js))), material(h)) {
		t.Error("binary to JSON changed the keyset")
	}
}

func TestConvertKeysetSameFormat(t *testing.T) {
	h := newKeyset(t, "AES256_GCM")
	js := mustRun(t, jsonKeyset(t, h), "convert-keyset")
	if !proto.Equal(material(parseJSON(t, []byte(js))), material(h)) {
		t.Error("JSON to JSON changed the keyset")
	}
}

func TestConvertKeysetWrongInFormat(t *testing.T) {
	r := run(t, jsonKeyset(t, newKeyset(t, "AES256_GCM")), "convert-keyset", "--in-format", "binary")
	if r.code != 1 || r.stdout != "" {
		t.Errorf("exit %d, stdout %q; want 1 and nothing written", r.code, r.stdout)
	}
}
