package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/daead"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newDAEAD() *cobra.Command {
	return newGroup("daead", "Encrypt and decrypt deterministically with a deterministic AEAD keyset",
		cipherCommand[tink.DeterministicAEAD]{
			use: "encrypt", short: "Encrypt --in deterministically with the keyset's primary key",
			prim: daeadPrimitive, extra: "associated-data",
			factory: daead.New, op: tink.DeterministicAEAD.EncryptDeterministically,
		}.command(),
		cipherCommand[tink.DeterministicAEAD]{
			use: "decrypt", short: "Decrypt --in with the keyset", decrypt: true,
			prim: daeadPrimitive, extra: "associated-data",
			factory: daead.New, op: tink.DeterministicAEAD.DecryptDeterministically,
		}.command(),
	)
}
