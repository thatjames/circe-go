package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func newDisableKey() *cobra.Command {
	return newKeyIDCommand("disable-key", "Disable a key that is not the primary",
		func(m *keyset.Manager, id uint32) error { return m.Disable(id) })
}
