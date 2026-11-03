package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// tagCommand computes a MAC tag or a signature of the data: the P that
// factory makes runs op, and the result is written in --encoding.
type tagCommand[P any] struct {
	use, short string
	prim       primitive
	kind       keysetKind
	factory    func(*keyset.Handle) (P, error)
	op         func(p P, data []byte) ([]byte, error)
}

func (c tagCommand[P]) command() *cobra.Command {
	cmd := &cobra.Command{Use: c.use, Short: c.short}
	flags := addPrimitiveFlags(cmd, true)
	enc := addEncodingFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		h, err := flags.handle(cmd, c.prim, c.kind)
		if err != nil {
			return err
		}
		p, err := c.factory(h)
		if err != nil {
			return err
		}
		data, _, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		tag, err := c.op(p, data)
		if err != nil {
			return err
		}
		return flags.write(cmd, enc.encode(tag))
	}
	return cmd
}

// verifyCommand checks a MAC tag or signature, read in --encoding from the
// file its flag names, against the data. It prints nothing and exits 0
// when op accepts it, and fails with exit 3 when op does not.
type verifyCommand[P any] struct {
	use, short string
	prim       primitive
	kind       keysetKind
	flag       string // tag or signature
	failure    string // the operation that failed: "MAC verification"
	factory    func(*keyset.Handle) (P, error)
	op         func(p P, tag, data []byte) error
}

func (c verifyCommand[P]) command() *cobra.Command {
	cmd := &cobra.Command{Use: c.use, Short: c.short}
	flags := addPrimitiveFlags(cmd, false)
	enc := addEncodingFlag(cmd)
	var tagPath string
	cmd.Flags().StringVar(&tagPath, c.flag, "", fmt.Sprintf("read the %s from this file (required)", c.flag))
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if err := required(c.flag, tagPath); err != nil {
			return err
		}
		h, err := flags.handle(cmd, c.prim, c.kind)
		if err != nil {
			return err
		}
		p, err := c.factory(h)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(tagPath)
		if err != nil {
			return err
		}
		tag, err := enc.decode(b, tagPath)
		if err != nil {
			return err
		}
		data, _, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		if err := c.op(p, tag, data); err != nil {
			return failed(c.failure)
		}
		return nil
	}
	return cmd
}
