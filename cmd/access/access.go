package access

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "access",
	Short: "Manage Zero Trust Access",
}

var appCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"application"},
	Short:   "Manage Access Applications",
}

func init() {
	Cmd.AddCommand(appCmd)
	appCmd.AddCommand(listCmd)
	appCmd.AddCommand(getCmd)
	appCmd.AddCommand(createCmd)
	appCmd.AddCommand(updateCmd)
	appCmd.AddCommand(deleteCmd)
}
