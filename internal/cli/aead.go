package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newAEAD() *cobra.Command {
	return newGroup("aead", "Encrypt and decrypt with an AEAD keyset",
		cipherCommand[tink.AEAD]{
			use: "encrypt", short: "Encrypt --in with the keyset's primary key",
			prim: aeadPrimitive, extra: "associated-data",
			factory: aead.New, op: tink.AEAD.Encrypt,
		}.command(),
		cipherCommand[tink.AEAD]{
			use: "decrypt", short: "Decrypt --in with the keyset", decrypt: true,
			prim: aeadPrimitive, extra: "associated-data",
			factory: aead.New, op: tink.AEAD.Decrypt,
		}.command(),
	)
}
