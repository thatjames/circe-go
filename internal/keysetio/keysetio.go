// Package keysetio reads and writes Tink keysets in Tink's JSON and binary
// formats.
package keysetio

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// Format is a keyset serialization: "json" or "binary". It is a pflag.Value.
type Format string

const (
	JSON   Format = "json"
	Binary Format = "binary"
)

func (f *Format) String() string { return string(*f) }

func (f *Format) Set(s string) error {
	switch Format(s) {
	case JSON, Binary:
		*f = Format(s)
		return nil
	}
	return fmt.Errorf("format %q: want json or binary", s)
}

func (f *Format) Type() string { return "json|binary" }

// Reader returns the tink-go keyset reader for f over r.
func (f Format) Reader(r io.Reader) keyset.Reader {
	if f == Binary {
		return keyset.NewBinaryReader(r)
	}
	return keyset.NewJSONReader(r)
}

// Writer returns the tink-go keyset writer for f over w.
func (f Format) Writer(w io.Writer) keyset.Writer {
	if f == Binary {
		return keyset.NewBinaryWriter(w)
	}
	return keyset.NewJSONWriter(w)
}

// ReadAll reads all of r, failing on empty input. name is the input's
// path, or "stdin", for errors.
func ReadAll(r io.Reader, name string) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", name, err)
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("reading %s: empty input", name)
	}
	return b, nil
}

// Read reads a cleartext keyset in format f from r.
func Read(r io.Reader, name string, f Format) (*keyset.Handle, error) {
	b, err := ReadAll(r, name)
	if err != nil {
		return nil, err
	}
	h, err := insecurecleartextkeyset.Read(f.Reader(bytes.NewReader(b)))
	if err != nil {
		return nil, fmt.Errorf("reading %s as a %s keyset: %w", name, f, err)
	}
	return h, nil
}

// Write writes h in cleartext in format f to w.
func Write(w io.Writer, h *keyset.Handle, f Format) error {
	if h == nil {
		return errors.New("no keyset to write")
	}
	if err := insecurecleartextkeyset.Write(h, f.Writer(w)); err != nil {
		return fmt.Errorf("writing %s keyset: %w", f, err)
	}
	return EndJSON(w, f)
}

// EndJSON ends a JSON keyset with a newline; tink-go writes none.
func EndJSON(w io.Writer, f Format) error {
	if f != JSON {
		return nil
	}
	_, err := io.WriteString(w, "\n")
	return err
}
