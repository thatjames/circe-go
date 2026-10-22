package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// tinkError drops tink-go's "keyset.Manager: " prefix from err.
func tinkError(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(strings.TrimPrefix(err.Error(), "keyset.Manager: "))
}

// newKeyIDCommand builds a command that reads a keyset, applies op to the
// key --key-id names, and writes the keyset.
func newKeyIDCommand(use, short string, op func(m *keyset.Manager, id uint32) error) *cobra.Command {
	var id uint32
	cmd := &cobra.Command{
		Use:   use + " --key-id ID",
		Short: short,
	}
	cmd.Flags().Uint32Var(&id, "key-id", 0, "ID of the key (see list-keyset)")
	k := addKMSFlags(cmd, "")
	in := addInFlags(cmd, k)
	out := addOutFlags(cmd, k)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if !cmd.Flags().Changed("key-id") {
			return usagef("--key-id is required")
		}
		h, err := in.read(cmd)
		if err != nil {
			return err
		}
		if !hasKey(h, id) {
			return fmt.Errorf("no key with ID %d in the keyset", id)
		}
		m := keyset.NewManagerFromHandle(h)
		if err := op(m, id); err != nil {
			return tinkError(err)
		}
		h, err = m.Handle()
		if err != nil {
			return tinkError(err)
		}
		return out.write(cmd, h)
	}
	return cmd
}

func hasKey(h *keyset.Handle, id uint32) bool {
	for _, k := range h.KeysetInfo().GetKeyInfo() {
		if k.GetKeyId() == id {
			return true
		}
	}
	return false
}
