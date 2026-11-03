package cli

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/tink-crypto/tink-go/v2/prf"
)

func TestPRFComputeEqualsTinkGo(t *testing.T) {
	for _, c := range []struct {
		template string
		length   uint32
	}{
		{"HMAC_SHA256_PRF", 32},
		{"HMAC_SHA512_PRF", 16},
		{"HKDF_SHA256", 100},
		{"AES_CMAC_PRF", 16},
	} {
		h := newKeyset(t, c.template)
		set, err := prf.NewPRFSet(h)
		if err != nil {
			t.Fatal(err)
		}
		want, err := set.ComputePrimaryPRF([]byte("user-42"), c.length)
		if err != nil {
			t.Fatal(err)
		}
		ks := keysetFile(t, h)
		for _, enc := range encodings {
			out := mustRun(t, []byte("user-42"), "prf", "compute", "--keyset", ks, "--output-length", strconv.Itoa(int(c.length)), "--encoding", enc)
			if got := decodeAs(t, enc, out); !bytes.Equal(got, want) {
				t.Errorf("%s, %s: circe's output differs from tink-go's", c.template, enc)
			}
		}
	}
}

func TestPRFLengthLimits(t *testing.T) {
	for _, c := range []struct {
		template, length, want string
	}{
		{"HMAC_SHA256_PRF", "33", "circe: --output-length 33 is more than the 32 bytes the primary PRF (HMAC-PRF SHA256) gives\n"},
		{"HMAC_SHA512_PRF", "65", "circe: --output-length 65 is more than the 64 bytes the primary PRF (HMAC-PRF SHA512) gives\n"},
		{"HKDF_SHA256", "8161", "circe: --output-length 8161 is more than the 8160 bytes the primary PRF (HKDF-PRF SHA256) gives\n"},
		{"AES_CMAC_PRF", "17", "circe: --output-length 17 is more than the 16 bytes the primary PRF (AES-CMAC-PRF) gives\n"},
		{"HKDF_SHA256", "4000000000", "circe: --output-length 4000000000 is more than the 8160 bytes the primary PRF (HKDF-PRF SHA256) gives\n"},
	} {
		ks := keysetFile(t, newKeyset(t, c.template))
		r := run(t, []byte("x"), "prf", "compute", "--keyset", ks, "--output-length", c.length)
		if r.code != 1 || r.stderr != c.want {
			t.Errorf("%s %s: exit %d, stderr %q; want 1 and %q", c.template, c.length, r.code, r.stderr, c.want)
		}
	}
	// The limits themselves are allowed.
	ks := keysetFile(t, newKeyset(t, "HKDF_SHA256"))
	if out := mustRun(t, []byte("x"), "prf", "compute", "--keyset", ks, "--output-length", "8160"); len(out) != 8160 {
		t.Errorf("HKDF at its limit wrote %d bytes", len(out))
	}
}

func TestPRFOutputLengthUsage(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "HMAC_SHA256_PRF"))
	for _, args := range [][]string{
		{"prf", "compute", "--keyset", ks},
		{"prf", "compute", "--keyset", ks, "--output-length", "0"},
		{"prf", "compute", "--keyset", ks, "--output-length", "-1"},
	} {
		if r := run(t, []byte("x"), args...); r.code != 2 {
			t.Errorf("circe %v: exit %d, want 2", args, r.code)
		}
	}
}

func TestPRFWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "HMAC_SHA256_128BITTAG"))
	r := run(t, []byte("x"), "prf", "compute", "--keyset", ks, "--output-length", "16")
	want := "circe: " + ks + " is not a PRF keyset (type.googleapis.com/google.crypto.tink.HmacKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}
