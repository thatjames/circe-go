package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCreatePublicKeyset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-public-keyset",
		Short: "Write the public keyset of a private signature, hybrid or JWT signature keyset",
	}
	// The public keyset holds no secrets, so it is written in cleartext.
	in := addInFlags(cmd, addKMSFlags(cmd, ""))
	out := addOutFlags(cmd, nil)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		h, err := in.read(cmd)
		if err != nil {
			return err
		}
		pub, err := h.Public()
		if err != nil {
			return fmt.Errorf("the keyset has no private keys to make public (signature, hybrid or JWT signature): %w", err)
		}
		return out.write(cmd, pub)
	}
	return cmd
}
