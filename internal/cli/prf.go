package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/prf"
	"github.com/tink-crypto/tink-go/v2/prf/aescmacprf"
	"github.com/tink-crypto/tink-go/v2/prf/hkdfprf"
	"github.com/tink-crypto/tink-go/v2/prf/hmacprf"
)

func newPRF() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Write --output-length bytes of the keyset's primary PRF over --in",
	}
	flags := addPrimitiveFlags(cmd, true)
	enc := addEncodingFlag(cmd)
	var length uint32
	cmd.Flags().Uint32Var(&length, "output-length", 0, "number of bytes to write, at least 1 (required)")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if !cmd.Flags().Changed("output-length") {
			return usagef("--output-length is required")
		}
		if length == 0 {
			return usagef("--output-length must be at least 1")
		}
		h, err := flags.handle(cmd, prfPrimitive, anyKeyset)
		if err != nil {
			return err
		}
		if err := checkPRFLength(h, length); err != nil {
			return err
		}
		set, err := prf.NewPRFSet(h)
		if err != nil {
			return err
		}
		data, _, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		out, err := set.ComputePrimaryPRF(data, length)
		if err != nil {
			return err
		}
		return flags.write(cmd, enc.encode(out))
	}
	return newGroup("prf", "Compute pseudorandom output with a PRF keyset", cmd)
}

// hashSizes are the output sizes of the hashes PRF keys name.
var hashSizes = map[string]uint32{"SHA1": 20, "SHA224": 28, "SHA256": 32, "SHA384": 48, "SHA512": 64}

// checkPRFLength fails when the primary PRF of h cannot give length bytes.
// tink-go checks too, but its HKDF error does not name the limit, and it
// allocates the output first.
func checkPRFLength(h *keyset.Handle, length uint32) error {
	e, err := h.Primary()
	if err != nil {
		return err
	}
	var (
		name  string
		limit uint32
	)
	switch p := e.Key().Parameters().(type) {
	case *hmacprf.Parameters:
		name, limit = "HMAC-PRF "+p.HashType().String(), hashSizes[p.HashType().String()]
	case *hkdfprf.Parameters:
		name, limit = "HKDF-PRF "+p.HashType().String(), 255*hashSizes[p.HashType().String()]
	case *aescmacprf.Parameters:
		name, limit = "AES-CMAC-PRF", 16
	default:
		return nil
	}
	if length > limit {
		return fmt.Errorf("--output-length %d is more than the %d bytes the primary PRF (%s) gives", length, limit, name)
	}
	return nil
}
