package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func newDeleteKey() *cobra.Command {
	return newKeyIDCommand("delete-key", "Remove a key that is not the primary from the keyset",
		func(m *keyset.Manager, id uint32) error { return m.Delete(id) })
}
