package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/thatjames/circe-go/internal/keysetio"
)

// newGroup makes the command group of a primitive. Given no operation, it
// prints its help.
func newGroup(use, short string, ops ...*cobra.Command) *cobra.Command {
	g := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	for _, op := range ops {
		op.Args = noArgs
	}
	g.AddCommand(ops...)
	return g
}

// primitiveFlags are the flags every primitive command shares: the keyset it
// uses (--keyset, --keyset-format and the KMS flags, only ever read) and its
// data (--in, --out).
type primitiveFlags struct {
	keyset inFlags
	in     string
	out    outFlags
}

func addPrimitiveFlags(cmd *cobra.Command) *primitiveFlags {
	p := &primitiveFlags{keyset: inFlags{format: keysetio.JSON, kms: addKMSFlags(cmd, "")}}
	cmd.Flags().StringVar(&p.keyset.path, "keyset", "", "read the keyset from this file (required)")
	cmd.Flags().Var(&p.keyset.format, "keyset-format", "keyset format: json or binary")
	cmd.Flags().StringVar(&p.in, "in", "", "read the data from this file (default: stdin)")
	cmd.Flags().StringVar(&p.out.path, "out", "", "write the result to this new file (default: stdout); never overwrites")
	return p
}

// handle reads the keyset and checks that it is a keyset of prim, and of
// the kind the command needs.
func (p *primitiveFlags) handle(cmd *cobra.Command, prim primitive, kind keysetKind) (*keyset.Handle, error) {
	if err := required("keyset", p.keyset.path); err != nil {
		return nil, err
	}
	h, err := p.keyset.read(cmd)
	if err != nil {
		return nil, err
	}
	if err := checkKeyset(h, p.keyset.path, cmd.CommandPath(), prim, kind); err != nil {
		return nil, err
	}
	return h, nil
}

// readData reads all of the data, from --in or stdin. It also returns the
// input's name, for errors.
func (p *primitiveFlags) readData(cmd *cobra.Command) ([]byte, string, error) {
	if p.in == "" {
		b, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, "stdin", fmt.Errorf("reading stdin: %w", err)
		}
		return b, "stdin", nil
	}
	b, err := os.ReadFile(p.in)
	return b, p.in, err
}

// write writes a finished result to --out or stdout.
func (p *primitiveFlags) write(cmd *cobra.Command, b []byte) error {
	return p.out.emit(cmd, b)
}

// bytesFlags give one input as text (--associated-data) or as a file
// (--associated-data-file). At most one may be given; the default is empty.
type bytesFlags struct {
	name, text, file string
}

func addBytesFlags(cmd *cobra.Command, name string) *bytesFlags {
	f := &bytesFlags{name: name}
	what := strings.ReplaceAll(name, "-", " ")
	cmd.Flags().StringVar(&f.text, name, "", what+" (default: empty)")
	cmd.Flags().StringVar(&f.file, name+"-file", "", "read the "+what+" from this file")
	return f
}

// get returns the bytes given; both flags at once is a usage error.
func (f *bytesFlags) get(cmd *cobra.Command) ([]byte, error) {
	text, file := cmd.Flags().Changed(f.name), cmd.Flags().Changed(f.name+"-file")
	switch {
	case text && file:
		return nil, usagef("--%s and --%s-file are mutually exclusive", f.name, f.name)
	case file:
		return os.ReadFile(f.file)
	}
	return []byte(f.text), nil
}
