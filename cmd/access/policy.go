package access

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage Access Policies",
}

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List reusable Access Policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		page, err := client.ZeroTrust.Access.Policies.List(context.TODO(), zero_trust.AccessPolicyListParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to list policies: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDECISION\tAPP COUNT\tSESSION\tCREATED")
		for _, p := range page.Result {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
				p.ID,
				p.Name,
				p.Decision,
				p.AppCount,
				p.SessionDuration,
				p.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
		w.Flush()
		return nil
	},
}

var policyGetCmd = &cobra.Command{
	Use:   "get <policy-id>",
	Short: "Get details of a reusable Access Policy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		p, err := client.ZeroTrust.Access.Policies.Get(context.TODO(), args[0], zero_trust.AccessPolicyGetParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to get policy: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(p)
	},
}

func init() {
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyGetCmd)
}
