package cli

import (
	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/thatjames/circe-go/internal/templates"
)

func newAddKey() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add-key --key-template NAME",
		Short: "Add a new enabled, non-primary key from a Tink key template",
	}
	cmd.Flags().StringVar(&name, "key-template", "", "Tink key template name (see list-key-templates)")
	k := addKMSFlags(cmd, "")
	in := addInFlags(cmd, k)
	out := addOutFlags(cmd, k)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if err := required("key-template", name); err != nil {
			return err
		}
		kt, err := templates.Lookup(name)
		if err != nil {
			return err
		}
		h, err := in.read(cmd)
		if err != nil {
			return err
		}
		m := keyset.NewManagerFromHandle(h)
		if _, err := m.Add(kt); err != nil {
			return tinkError(err)
		}
		h, err = m.Handle()
		if err != nil {
			return tinkError(err)
		}
		return out.write(cmd, h)
	}
	return cmd
}
