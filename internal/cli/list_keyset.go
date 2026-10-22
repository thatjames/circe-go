package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// keysetInfo is Tink's KeysetInfo with its proto JSON field names.
type keysetInfo struct {
	PrimaryKeyID uint32    `json:"primaryKeyId"`
	KeyInfo      []keyInfo `json:"keyInfo"`
}

type keyInfo struct {
	TypeURL          string `json:"typeUrl"`
	Status           string `json:"status"`
	KeyID            uint32 `json:"keyId"`
	OutputPrefixType string `json:"outputPrefixType"`
}

func newListKeyset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-keyset",
		Short: "Show a keyset's key IDs, types and statuses, without key material",
	}
	in := addInFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		h, err := in.read(cmd)
		if err != nil {
			return err
		}
		ki := h.KeysetInfo()
		info := keysetInfo{PrimaryKeyID: ki.GetPrimaryKeyId(), KeyInfo: []keyInfo{}}
		for _, k := range ki.GetKeyInfo() {
			info.KeyInfo = append(info.KeyInfo, keyInfo{
				TypeURL:          k.GetTypeUrl(),
				Status:           k.GetStatus().String(),
				KeyID:            k.GetKeyId(),
				OutputPrefixType: k.GetOutputPrefixType().String(),
			})
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}
	return cmd
}
