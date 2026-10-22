package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func newPromoteKey() *cobra.Command {
	return newKeyIDCommand("promote-key", "Make an enabled key the keyset's primary",
		func(m *keyset.Manager, id uint32) error { return m.SetPrimary(id) })
}
