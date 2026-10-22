package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/thatjames/circe-go/internal/templates"
)

func newCreateKeyset() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create-keyset --key-template NAME",
		Short: "Create a keyset with one primary key from a Tink key template",
	}
	cmd.Flags().StringVar(&name, "key-template", "", "Tink key template name (see list-key-templates)")
	out := addOutFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if err := required("key-template", name); err != nil {
			return err
		}
		kt, err := templates.Lookup(name)
		if err != nil {
			return err
		}
		h, err := keyset.NewHandle(kt)
		if err != nil {
			return err
		}
		return out.write(cmd, h)
	}
	return cmd
}
