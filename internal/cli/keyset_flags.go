package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/thatjames/circe-go/internal/keysetio"
)

// inFlags are the --in and --in-format flags of a command that reads a keyset.
type inFlags struct {
	path   string
	format keysetio.Format
}

func addInFlags(cmd *cobra.Command) *inFlags {
	f := &inFlags{format: keysetio.JSON}
	cmd.Flags().StringVar(&f.path, "in", "", "read the keyset from this file (default: stdin)")
	cmd.Flags().Var(&f.format, "in-format", "input keyset format: json or binary")
	return f
}

// read reads the input keyset.
func (f *inFlags) read(cmd *cobra.Command) (*keyset.Handle, error) {
	if f.path == "" {
		return keysetio.Read(cmd.InOrStdin(), "stdin", f.format)
	}
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return keysetio.Read(file, f.path, f.format)
}

// outFlags are the --out and --out-format flags of a command that writes a keyset.
type outFlags struct {
	path   string
	format keysetio.Format
}

func addOutFlags(cmd *cobra.Command) *outFlags {
	f := &outFlags{format: keysetio.JSON}
	cmd.Flags().StringVar(&f.path, "out", "", "write the keyset to this new file (default: stdout); never overwrites")
	cmd.Flags().Var(&f.format, "out-format", "output keyset format: json or binary")
	return f
}

// write writes h as the output keyset.
func (f *outFlags) write(cmd *cobra.Command, h *keyset.Handle) error {
	var buf bytes.Buffer
	if err := keysetio.Write(&buf, h, f.format); err != nil {
		return err
	}
	return f.emit(cmd, buf.Bytes())
}

// emit writes a finished output, to stdout or to a file that must not exist yet.
func (f *outFlags) emit(cmd *cobra.Command, b []byte) error {
	if f.path == "" {
		_, err := cmd.OutOrStdout().Write(b)
		return err
	}
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("%s exists; circe never overwrites a file", f.path)
	}
	if err != nil {
		return err
	}
	if _, err := file.Write(b); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}
