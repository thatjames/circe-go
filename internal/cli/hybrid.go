package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/hybrid"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newHybrid() *cobra.Command {
	return newGroup("hybrid", "Encrypt to a public hybrid keyset and decrypt with its private keyset",
		cipherCommand[tink.HybridEncrypt]{
			use: "encrypt", short: "Encrypt --in to the public keyset's primary key",
			prim: hybridPrimitive, kind: publicKeyset, extra: "context-info",
			factory: hybrid.NewHybridEncrypt, op: tink.HybridEncrypt.Encrypt,
		}.command(),
		cipherCommand[tink.HybridDecrypt]{
			use: "decrypt", short: "Decrypt --in with the private keyset", decrypt: true,
			prim: hybridPrimitive, kind: privateKeyset, extra: "context-info",
			factory: hybrid.NewHybridDecrypt, op: tink.HybridDecrypt.Decrypt,
		}.command(),
	)
}
