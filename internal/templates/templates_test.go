package templates

import (
	"sort"
	"strings"
	"testing"

	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

// wantPrefix is the output prefix a name promises. Streaming AEAD and PRF
// keys have no prefix in Tink, whatever their name.
func wantPrefix(name, typeURL string) tinkpb.OutputPrefixType {
	if strings.Contains(typeURL, "Streaming") || strings.Contains(typeURL, "Prf") {
		return tinkpb.OutputPrefixType_RAW
	}
	for _, s := range []string{"_RAW", "_NO_PREFIX", "_WITHOUT_PREFIX"} {
		if strings.Contains(name, s) {
			return tinkpb.OutputPrefixType_RAW
		}
	}
	return tinkpb.OutputPrefixType_TINK
}

func TestEveryNameCreatesAKeyset(t *testing.T) {
	for _, name := range Names() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			kt, err := Lookup(name)
			if err != nil {
				t.Fatal(err)
			}
			h, err := keyset.NewHandle(kt)
			if err != nil {
				t.Fatalf("keyset.NewHandle: %v", err)
			}
			info := h.KeysetInfo().GetKeyInfo()
			if len(info) != 1 {
				t.Fatalf("got %d keys, want 1", len(info))
			}
			if got, want := info[0].GetOutputPrefixType(), wantPrefix(name, kt.GetTypeUrl()); got != want {
				t.Errorf("output prefix %v, want %v", got, want)
			}
		})
	}
}

func TestNamesAreSortedAndUpperCase(t *testing.T) {
	names := Names()
	if !sort.StringsAreSorted(names) {
		t.Error("Names() is not sorted")
	}
	for _, n := range names {
		if n != strings.ToUpper(n) {
			t.Errorf("%q is not upper case", n)
		}
	}
}

func TestUnknownName(t *testing.T) {
	for _, name := range []string{"NOPE", "aes256_gcm", ""} {
		_, err := Lookup(name)
		if err == nil || !strings.Contains(err.Error(), "circe list-key-templates") {
			t.Errorf("Lookup(%q) = %v, want an error suggesting list-key-templates", name, err)
		}
	}
}

func TestRawDoesNotChangeTheTinkTemplate(t *testing.T) {
	for _, name := range []string{"AES256_SIV", "CHACHA20_POLY1305", "AES256_CMAC"} {
		if _, err := Lookup(name + "_RAW"); err != nil {
			t.Fatal(err)
		}
		kt, err := Lookup(name)
		if err != nil {
			t.Fatal(err)
		}
		if kt.GetOutputPrefixType() != tinkpb.OutputPrefixType_TINK {
			t.Errorf("%s: output prefix %v after a _RAW lookup", name, kt.GetOutputPrefixType())
		}
	}
}
