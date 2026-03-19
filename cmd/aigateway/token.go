package aigateway

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

const aiGatewayReadPermGroupID = "e8fed01d18df95aff5765ad66e6e4e78"

var tokenCreateName string

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage AI Gateway API tokens",
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a read-only API token for AI Gateway",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		// First, look up the AI Gateway Read permission group ID dynamically
		permGroupID, err := findAIGatewayReadPermGroup(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not look up permission group, using default: %v\n", err)
			permGroupID = aiGatewayReadPermGroupID
		}

		res, err := client.User.Tokens.New(context.TODO(), user.TokenNewParams{
			Name: cloudflare.F(tokenCreateName),
			Policies: cloudflare.F([]shared.TokenPolicyParam{
				{
					Effect: cloudflare.F(shared.TokenPolicyEffectAllow),
					PermissionGroups: cloudflare.F([]shared.TokenPolicyPermissionGroupParam{
						{ID: cloudflare.F(permGroupID)},
					}),
					Resources: cloudflare.F(map[string]shared.TokenPolicyResourcesUnionParam{
						fmt.Sprintf("com.cloudflare.api.account.%s", accountID): shared.UnionString("*"),
					}),
				},
			}),
		})
		if err != nil {
			return fmt.Errorf("failed to create token: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully created token: %s\n", tokenCreateName)
		fmt.Fprintf(os.Stderr, "Token ID: %s\n", res.ID)
		fmt.Println(res.Value)
		return nil
	},
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API tokens (filtered by AI Gateway name pattern)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()

		page, err := client.User.Tokens.List(context.TODO(), user.TokenListParams{})
		if err != nil {
			return fmt.Errorf("failed to list tokens: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tLAST USED\tEXPIRES")
		for _, t := range page.Result {
			lastUsed := "-"
			if !t.LastUsedOn.IsZero() {
				lastUsed = t.LastUsedOn.Format("2006-01-02 15:04")
			}
			expires := "-"
			if !t.ExpiresOn.IsZero() {
				expires = t.ExpiresOn.Format("2006-01-02")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.ID, t.Name, t.Status, lastUsed, expires,
			)
		}
		w.Flush()
		return nil
	},
}

func findAIGatewayReadPermGroup(client *cloudflare.Client) (string, error) {
	page, err := client.User.Tokens.PermissionGroups.List(context.TODO(), user.TokenPermissionGroupListParams{})
	if err != nil {
		return "", err
	}
	for _, pg := range page.Result {
		if pg.Name == "AI Gateway Read" {
			return pg.ID, nil
		}
	}
	return "", fmt.Errorf("AI Gateway Read permission group not found")
}

var tokenDeleteCmd = &cobra.Command{
	Use:   "delete <token-id>",
	Short: "Delete an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()

		_, err := client.User.Tokens.Delete(context.TODO(), args[0])
		if err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully deleted token: %s\n", args[0])
		return nil
	},
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenCreateName, "name", "aig-token", "Token name")

	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenDeleteCmd)
}
