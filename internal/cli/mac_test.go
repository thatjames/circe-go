package cli

import (
	"bytes"
	"testing"

	"github.com/tink-crypto/tink-go/v2/mac"
)

func TestMACComputeAndVerifyEveryEncoding(t *testing.T) {
	h := newKeyset(t, "HMAC_SHA256_256BITTAG")
	ks := keysetFile(t, h)
	m, err := mac.New(h)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("authenticate me")
	want, err := m.ComputeMAC(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, enc := range encodings {
		out := mustRun(t, data, "mac", "compute", "--keyset", ks, "--encoding", enc)
		if got := decodeAs(t, enc, out); !bytes.Equal(got, want) {
			t.Errorf("%s: circe's tag differs from tink-go's", enc)
		}
		tag := dataFile(t, "tag", []byte(out))
		if r := run(t, data, "mac", "verify", "--keyset", ks, "--tag", tag, "--encoding", enc); r.code != 0 || r.stdout != "" || r.stderr != "" {
			t.Errorf("%s: verify exit %d, stdout %q, stderr %q; want 0 and nothing", enc, r.code, r.stdout, r.stderr)
		}
	}
}

func TestMACVerifyFailures(t *testing.T) {
	h := newKeyset(t, "HMAC_SHA256_128BITTAG")
	ks := keysetFile(t, h)
	tag := []byte(mustRun(t, []byte("data"), "mac", "compute", "--keyset", ks))
	bad := bytes.Clone(tag)
	bad[len(bad)-1] ^= 1
	other := keysetFile(t, newKeyset(t, "HMAC_SHA256_128BITTAG"))
	for name, c := range map[string]struct {
		ks, data string
		tag      []byte
	}{
		"other data": {ks, "other", tag},
		"bad tag":    {ks, "data", bad},
		"other key":  {other, "data", tag},
	} {
		r := run(t, []byte(c.data), "mac", "verify", "--keyset", c.ks, "--tag", dataFile(t, "tag", c.tag))
		if r.code != 3 || r.stderr != "circe: MAC verification failed\n" || r.stdout != "" {
			t.Errorf("%s: exit %d, stderr %q; want 3 and MAC verification failed", name, r.code, r.stderr)
		}
	}
}

func TestMACVerifyUsage(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "HMAC_SHA256_128BITTAG"))
	for _, args := range [][]string{
		{"mac", "verify", "--keyset", ks},
		{"mac", "verify", "--keyset", ks, "--tag", ks, "--out", "x"},
	} {
		if r := run(t, []byte("data"), args...); r.code != 2 {
			t.Errorf("circe %v: exit %d, want 2", args, r.code)
		}
	}
	r := run(t, []byte("data"), "mac", "verify", "--keyset", ks, "--tag", dataFile(t, "tag", []byte("zz")), "--encoding", "hex")
	if r.code != 1 {
		t.Errorf("tag that is not hex: exit %d, want 1", r.code)
	}
}

func TestMACWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "HMAC_SHA256_PRF"))
	r := run(t, []byte("x"), "mac", "compute", "--keyset", ks)
	want := "circe: " + ks + " is not a MAC keyset (type.googleapis.com/google.crypto.tink.HmacPrfKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}
