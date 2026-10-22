package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thatjames/circe-go/internal/templates"
)

func newListKeyTemplates() *cobra.Command {
	return &cobra.Command{
		Use:   "list-key-templates",
		Short: "List the Tink key template names --key-template accepts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, n := range templates.Names() {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), n); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
