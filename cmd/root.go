package cmd

import (
	"fmt"
	"os"

	"github.com/hashrock/hashflare/cmd/access"
	"github.com/hashrock/hashflare/cmd/aigateway"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashflare",
	Short: "Cloudflare management CLI",
	Long:  "A wrangler-like CLI tool for managing Cloudflare AI Gateway and Zero Trust Access Applications.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(aigateway.Cmd)
	rootCmd.AddCommand(access.Cmd)
}
