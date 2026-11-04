package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodingRoundTrip(t *testing.T) {
	data := []byte{0, 1, 0xfe, 0xff, 'x'}
	for _, e := range []encoding{rawEncoding, base64Encoding, hexEncoding} {
		out := e.encode(data)
		if e != rawEncoding && !bytes.HasSuffix(out, []byte("\n")) {
			t.Errorf("%s: output %q does not end in a newline", e, out)
		}
		got, err := e.decode(out, "stdin")
		if err != nil || !bytes.Equal(got, data) {
			t.Errorf("%s: decode(encode) = %x, %v; want %x", e, got, err, data)
		}
	}
	if got := string(hexEncoding.encode([]byte{0xab})); got != "ab\n" {
		t.Errorf("hex = %q, want lower case ab", got)
	}
	if got := string(base64Encoding.encode([]byte{0xff})); got != "/w==\n" {
		t.Errorf("base64 = %q, want standard padded /w==", got)
	}
}

func TestEncodingIgnoresSurroundingWhitespace(t *testing.T) {
	for e, in := range map[encoding]string{base64Encoding: " \n/w==\r\n\t", hexEncoding: "\n ff \n"} {
		if got, err := e.decode([]byte(in), "stdin"); err != nil || !bytes.Equal(got, []byte{0xff}) {
			t.Errorf("%s: decode(%q) = %x, %v", e, in, got, err)
		}
	}
}

func TestEncodingBadInputNamesInputNotData(t *testing.T) {
	for _, e := range []encoding{base64Encoding, hexEncoding} {
		_, err := e.decode([]byte("secretgarbage!"), "tag.txt")
		if err == nil {
			t.Fatalf("%s: no error", e)
		}
		msg := err.Error()
		if !strings.Contains(msg, "tag.txt") || !strings.Contains(msg, string(e)) || strings.Contains(msg, "secret") {
			t.Errorf("%s: error %q, want the input and encoding named and no data", e, msg)
		}
	}
}

func TestEncodingFlagRejectsUnknown(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "AES128_GCM"))
	if r := run(t, nil, "aead", "encrypt", "--keyset", ks, "--encoding", "base32"); r.code != 2 {
		t.Errorf("--encoding base32: exit %d, want 2", r.code)
	}
}
