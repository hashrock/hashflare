package access

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var appPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage policies for an Access Application",
}

var appPolicyListCmd = &cobra.Command{
	Use:   "list <app-id>",
	Short: "List policies for an Access Application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		page, err := client.ZeroTrust.Access.Applications.Policies.List(
			context.TODO(), args[0],
			zero_trust.AccessApplicationPolicyListParams{
				AccountID: cloudflare.F(accountID),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to list app policies: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDECISION\tPRECEDENCE")
		for _, p := range page.Result {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				p.ID,
				p.Name,
				p.Decision,
				p.Precedence,
			)
		}
		w.Flush()
		return nil
	},
}

var appPolicyGetCmd = &cobra.Command{
	Use:   "get <app-id> <policy-id>",
	Short: "Get details of a policy for an Access Application",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		p, err := client.ZeroTrust.Access.Applications.Policies.Get(
			context.TODO(), args[0], args[1],
			zero_trust.AccessApplicationPolicyGetParams{
				AccountID: cloudflare.F(accountID),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to get app policy: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(p)
	},
}

var (
	appPolicyAddName     string
	appPolicyAddDecision string
	appPolicyAddEmails   []string
)

var appPolicyAddCmd = &cobra.Command{
	Use:   "add <app-id>",
	Short: "Add a policy to an Access Application",
	Long: `Add a policy to an Access Application.

Example:
  hashflare access app policy add <app-id> --name "Allow team" --decision allow --email user@example.com --email admin@example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		includeRules := make([]map[string]interface{}, 0, len(appPolicyAddEmails))
		for _, email := range appPolicyAddEmails {
			includeRules = append(includeRules, map[string]interface{}{
				"email": map[string]string{"email": email},
			})
		}

		p, err := client.ZeroTrust.Access.Applications.Policies.New(
			context.TODO(), args[0],
			zero_trust.AccessApplicationPolicyNewParams{
				AccountID: cloudflare.F(accountID),
			},
			option.WithJSONSet("name", appPolicyAddName),
			option.WithJSONSet("decision", appPolicyAddDecision),
			option.WithJSONSet("include", includeRules),
		)
		if err != nil {
			return fmt.Errorf("failed to add policy: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully added policy: %s\n", appPolicyAddName)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(p)
	},
}

var appPolicyDeleteCmd = &cobra.Command{
	Use:   "delete <app-id> <policy-id>",
	Short: "Delete a policy from an Access Application",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		_, err := client.ZeroTrust.Access.Applications.Policies.Delete(
			context.TODO(), args[0], args[1],
			zero_trust.AccessApplicationPolicyDeleteParams{
				AccountID: cloudflare.F(accountID),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to delete policy: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully deleted policy: %s\n", args[1])
		return nil
	},
}

func init() {
	appPolicyAddCmd.Flags().StringVar(&appPolicyAddName, "name", "", "Policy name (required)")
	appPolicyAddCmd.Flags().StringVar(&appPolicyAddDecision, "decision", "allow", "Policy decision (allow, deny, bypass, non_identity)")
	appPolicyAddCmd.Flags().StringSliceVar(&appPolicyAddEmails, "email", nil, "Email addresses to include (repeatable)")
	appPolicyAddCmd.MarkFlagRequired("name")
	appPolicyAddCmd.MarkFlagRequired("email")

	appPolicyCmd.AddCommand(appPolicyListCmd)
	appPolicyCmd.AddCommand(appPolicyGetCmd)
	appPolicyCmd.AddCommand(appPolicyAddCmd)
	appPolicyCmd.AddCommand(appPolicyDeleteCmd)
}
