package aigateway

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "ai-gateway",
	Aliases: []string{"aig"},
	Short:   "Manage AI Gateways",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(tokenCmd)
}
