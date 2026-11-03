package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/signature"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newSignature() *cobra.Command {
	return newGroup("signature", "Sign with a private signature keyset and verify with its public keyset",
		tagCommand[tink.Signer]{
			use: "sign", short: "Sign --in with the private keyset's primary key",
			prim: signaturePrimitive, kind: privateKeyset,
			factory: signature.NewSigner, op: tink.Signer.Sign,
		}.command(),
		verifyCommand[tink.Verifier]{
			use: "verify", short: "Check a --signature against --in with the public keyset; exit 3 when it does not verify",
			prim: signaturePrimitive, kind: publicKeyset, flag: "signature", failure: "signature verification",
			factory: signature.NewVerifier, op: tink.Verifier.Verify,
		}.command(),
	)
}
