package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// cipherCommand is an encrypt or decrypt command over data held in memory:
// AEAD, deterministic AEAD and hybrid. P is the tink-go primitive that
// factory makes and op runs, with the bytes Tink binds to the ciphertext
// (associated data or context info) as its second argument. Encrypt writes
// the ciphertext in --encoding; decrypt reads it in --encoding, and a
// ciphertext that does not decrypt exits 3.
type cipherCommand[P any] struct {
	use, short string
	prim       primitive
	kind       keysetKind
	extra      string // the flag of the bound bytes: associated-data or context-info
	decrypt    bool
	factory    func(*keyset.Handle) (P, error)
	op         func(p P, data, extra []byte) ([]byte, error)
}

func (c cipherCommand[P]) command() *cobra.Command {
	cmd := &cobra.Command{Use: c.use, Short: c.short}
	flags := addPrimitiveFlags(cmd, true)
	enc := addEncodingFlag(cmd)
	extraFlags := addBytesFlags(cmd, c.extra)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		extra, err := extraFlags.get(cmd)
		if err != nil {
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
		data, name, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		if c.decrypt {
			if data, err = enc.decode(data, name); err != nil {
				return err
			}
		}
		out, err := c.op(p, data, extra)
		switch {
		case err != nil && c.decrypt:
			return failed("decryption")
		case err != nil:
			return err
		case !c.decrypt:
			out = enc.encode(out)
		}
		return flags.write(cmd, out)
	}
	return cmd
}
