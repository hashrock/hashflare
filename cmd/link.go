package cmd

import (
	"fmt"

	"github.com/hashrock/hashflare/internal/project"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link Cloudflare resources to this project",
}

var linkAIGCmd = &cobra.Command{
	Use:   "aig <gateway-id>",
	Short: "Link an AI Gateway to this project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := project.Load()
		if err != nil {
			return err
		}
		c.AIGateway = args[0]
		if err := project.Save(c); err != nil {
			return err
		}
		fmt.Printf("Linked AI Gateway: %s\n", args[0])
		return nil
	},
}

var linkAppCmd = &cobra.Command{
	Use:   "app <app-id>",
	Short: "Link an Access Application to this project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := project.Load()
		if err != nil {
			return err
		}
		c.AccessApp = args[0]
		if err := project.Save(c); err != nil {
			return err
		}
		fmt.Printf("Linked Access App: %s\n", args[0])
		return nil
	},
}

var linkPolicyCmd = &cobra.Command{
	Use:   "policy <policy-id>",
	Short: "Link an Access Policy to this project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := project.Load()
		if err != nil {
			return err
		}
		// avoid duplicates
		for _, id := range c.AccessPolicies {
			if id == args[0] {
				fmt.Printf("Policy already linked: %s\n", args[0])
				return nil
			}
		}
		c.AccessPolicies = append(c.AccessPolicies, args[0])
		if err := project.Save(c); err != nil {
			return err
		}
		fmt.Printf("Linked Access Policy: %s\n", args[0])
		return nil
	},
}

func init() {
	linkCmd.AddCommand(linkAIGCmd)
	linkCmd.AddCommand(linkAppCmd)
	linkCmd.AddCommand(linkPolicyCmd)
	rootCmd.AddCommand(linkCmd)
}
