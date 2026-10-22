package keysetio

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/tink"
)

// ReadEncrypted reads an EncryptedKeyset in format f from r and decrypts it
// with a, the AEAD of the KMS key uri.
func ReadEncrypted(r io.Reader, name string, f Format, a tink.AEAD, uri string) (*keyset.Handle, error) {
	b, err := ReadAll(r, name)
	if err != nil {
		return nil, err
	}
	h, err := keyset.Read(f.Reader(bytes.NewReader(b)), a)
	if err != nil {
		return nil, fmt.Errorf("reading %s as a %s keyset encrypted with %s: %w", name, f, uri, err)
	}
	return h, nil
}

// WriteEncrypted encrypts h with a and writes it as an EncryptedKeyset in
// format f to w.
func WriteEncrypted(w io.Writer, h *keyset.Handle, f Format, a tink.AEAD) error {
	if h == nil {
		return errors.New("no keyset to write")
	}
	if err := h.Write(f.Writer(w), a); err != nil {
		return fmt.Errorf("writing encrypted %s keyset: %w", f, err)
	}
	return EndJSON(w, f)
}
