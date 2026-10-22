// Command circe manages Tink keysets.
package main

import (
	"os"

	"github.com/thatjames/circe-go/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
