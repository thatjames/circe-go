package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/mac"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newMAC() *cobra.Command {
	return newGroup("mac", "Compute and verify MAC tags with a MAC keyset",
		tagCommand[tink.MAC]{
			use: "compute", short: "Write the MAC tag of --in, made with the keyset's primary key",
			prim: macPrimitive, factory: mac.New, op: tink.MAC.ComputeMAC,
		}.command(),
		verifyCommand[tink.MAC]{
			use: "verify", short: "Check a --tag against --in; exit 3 when it does not match",
			prim: macPrimitive, flag: "tag", failure: "MAC verification",
			factory: mac.New, op: tink.MAC.VerifyMAC,
		}.command(),
	)
}
