package cli

import "github.com/spf13/cobra"

func newConvertKeyset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-keyset",
		Short: "Rewrite a keyset between Tink's JSON and binary formats",
	}
	in := addInFlags(cmd)
	out := addOutFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		h, err := in.read(cmd)
		if err != nil {
			return err
		}
		return out.write(cmd, h)
	}
	return cmd
}
