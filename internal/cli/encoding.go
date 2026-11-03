package cli

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

// encoding is how the bytes Tink makes (ciphertexts, tags, signatures, PRF
// output) are written and read: raw, base64 (standard, padded) or hex
// (lower case). It is a pflag.Value.
type encoding string

const (
	rawEncoding    encoding = "raw"
	base64Encoding encoding = "base64"
	hexEncoding    encoding = "hex"
)

func (e *encoding) String() string { return string(*e) }

// Set parses s as an encoding; it is the pflag.Value setter.
func (e *encoding) Set(s string) error {
	switch encoding(s) {
	case rawEncoding, base64Encoding, hexEncoding:
		*e = encoding(s)
		return nil
	}
	return fmt.Errorf("encoding %q: want raw, base64 or hex", s)
}

// Type names the accepted values in flag help.
func (e *encoding) Type() string { return "raw|base64|hex" }

func addEncodingFlag(cmd *cobra.Command) *encoding {
	e := rawEncoding
	cmd.Flags().Var(&e, "encoding", "encoding of ciphertexts, tags, signatures and PRF output: raw, base64 or hex")
	return &e
}

// encode encodes b for output. Encoded output ends in a newline.
func (e encoding) encode(b []byte) []byte {
	switch e {
	case base64Encoding:
		return []byte(base64.StdEncoding.EncodeToString(b) + "\n")
	case hexEncoding:
		return []byte(hex.EncodeToString(b) + "\n")
	}
	return b
}

// decode decodes b, read from name (a path or "stdin"). Surrounding
// whitespace is ignored. The error names the input and the encoding, never
// the data.
func (e encoding) decode(b []byte, name string) ([]byte, error) {
	var (
		out []byte
		err error
	)
	switch e {
	case base64Encoding:
		out, err = base64.StdEncoding.DecodeString(string(bytes.TrimSpace(b)))
	case hexEncoding:
		out, err = hex.DecodeString(string(bytes.TrimSpace(b)))
	default:
		return b, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s is not valid %s", name, e)
	}
	return out, nil
}
