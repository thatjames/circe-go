// Package cli is the circe command line: the cobra root and one file per
// subcommand.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is set at link time by release builds:
// -X github.com/thatjames/circe-go/internal/cli.version=<tag>
var version string

func versionString() string {
	if version != "" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return "devel"
}

// usageError is an error in how circe was called; it exits 2 with the usage.
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }

func usagef(format string, a ...any) error { return usageError{fmt.Errorf(format, a...)} }

// noArgs rejects positional arguments as a usage error.
func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		if !cmd.HasParent() {
			return usagef("unknown command %q", args[0])
		}
		return usagef("unexpected argument %q", args[0])
	}
	return nil
}

// required fails with a usage error when a flag that must be given is empty.
func required(name, value string) error {
	if value == "" {
		return usagef("--%s is required", name)
	}
	return nil
}

func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           "circe",
		Short:         "circe manages Tink keysets",
		Long:          "circe creates, manages, inspects and converts Tink keysets, on tink-go.",
		Version:       versionString(),
		Args:          noArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("circe {{.Version}}\n")
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageError{err} })
	root.AddCommand(
		newCreateKeyset(),
		newAddKey(),
		newPromoteKey(),
		newEnableKey(),
		newDisableKey(),
		newDeleteKey(),
		newCreatePublicKeyset(),
		newListKeyset(),
		newListKeyTemplates(),
		newConvertKeyset(),
	)
	for _, c := range root.Commands() {
		if c.Args == nil {
			c.Args = noArgs
		}
	}
	return root
}

// Run runs circe with args and returns its exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root := newRoot()
	root.SetArgs(args)
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	cmd, err := root.ExecuteC()
	if err == nil {
		return 0
	}
	fmt.Fprintf(stderr, "circe: %s\n", err)
	var ue usageError
	if errors.As(err, &ue) {
		fmt.Fprint(stderr, cmd.UsageString())
		return 2
	}
	return 1
}

// Execute runs circe on the process's arguments and standard streams.
func Execute() int {
	return Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}
