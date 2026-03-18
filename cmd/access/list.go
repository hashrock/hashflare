package access

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Access Applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		page, err := client.ZeroTrust.Access.Applications.List(context.TODO(), zero_trust.AccessApplicationListParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to list access applications: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDOMAIN\tTYPE\tCREATED")
		for _, app := range page.Result {
			name := ""
			domain := ""
			appType := ""
			id := ""
			created := ""

			switch v := app.AsUnion().(type) {
			case zero_trust.AccessApplicationListResponseSelfHostedApplication:
				id = v.ID
				name = v.Name
				domain = v.Domain
				appType = "self_hosted"
				created = v.CreatedAt.Format("2006-01-02 15:04:05")
			case zero_trust.AccessApplicationListResponseSaaSApplication:
				id = v.ID
				name = v.Name
				appType = "saas"
				created = v.CreatedAt.Format("2006-01-02 15:04:05")
			case zero_trust.AccessApplicationListResponseBrowserSSHApplication:
				id = v.ID
				name = v.Name
				domain = v.Domain
				appType = "ssh"
				created = v.CreatedAt.Format("2006-01-02 15:04:05")
			case zero_trust.AccessApplicationListResponseBrowserVNCApplication:
				id = v.ID
				name = v.Name
				domain = v.Domain
				appType = "vnc"
				created = v.CreatedAt.Format("2006-01-02 15:04:05")
			case zero_trust.AccessApplicationListResponseBookmarkApplication:
				id = v.ID
				name = v.Name
				domain = v.Domain
				appType = "bookmark"
				created = v.CreatedAt.Format("2006-01-02 15:04:05")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", id, name, domain, appType, created)
		}
		w.Flush()
		return nil
	},
}
