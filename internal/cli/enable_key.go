package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func newEnableKey() *cobra.Command {
	return newKeyIDCommand("enable-key", "Enable a disabled key",
		func(m *keyset.Manager, id uint32) error { return m.Enable(id) })
}
