package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/tink"

	"github.com/thatjames/circe-go/internal/kms"
)

// kmsAEAD gets the AEAD of a KMS key URI; tests replace it with fakes.
var kmsAEAD = kms.AEAD

// kmsFlags are --master-key-uri and --credential, or with prefix "new-"
// the output side of convert-keyset.
type kmsFlags struct {
	uri, credential string

	aead tink.AEAD
	err  error
	done bool
}

func addKMSFlags(cmd *cobra.Command, prefix string) *kmsFlags {
	k := &kmsFlags{}
	cmd.Flags().StringVar(&k.uri, prefix+"master-key-uri", "", "KMS key that encrypts the keyset: gcp-kms://, aws-kms:// or hcvault://")
	cmd.Flags().StringVar(&k.credential, prefix+"credential", "", "credentials file for that KMS (default: the backend's default credentials)")
	return k
}

// get returns the KMS AEAD, or nil for a cleartext keyset. The AEAD is made
// once, so the input and output of a command share it.
func (k *kmsFlags) get(cmd *cobra.Command) (tink.AEAD, error) {
	if k == nil || k.uri == "" {
		return nil, nil
	}
	if !k.done {
		k.aead, k.err = kmsAEAD(cmd.Context(), k.uri, k.credential)
		k.done = true
	}
	return k.aead, k.err
}
