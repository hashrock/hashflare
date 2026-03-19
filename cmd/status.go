package cmd

import (
	"fmt"

	"github.com/hashrock/hashflare/internal/project"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project-level resource links",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := project.Load()
		if err != nil {
			return err
		}

		p, _ := project.FindPath()
		fmt.Printf("Project config: %s\n\n", p)

		if c.AIGateway != "" {
			fmt.Printf("AI Gateway:      %s\n", c.AIGateway)
		} else {
			fmt.Println("AI Gateway:      (not linked)")
		}

		if c.AccessApp != "" {
			fmt.Printf("Access App:      %s\n", c.AccessApp)
		} else {
			fmt.Println("Access App:      (not linked)")
		}

		if len(c.AccessPolicies) > 0 {
			fmt.Println("Access Policies:")
			for _, id := range c.AccessPolicies {
				fmt.Printf("  - %s\n", id)
			}
		} else {
			fmt.Println("Access Policies: (none)")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
