package access

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var (
	updateName   string
	updateDomain string
)

var updateCmd = &cobra.Command{
	Use:   "update <app-id>",
	Short: "Update an Access Application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		body := zero_trust.AccessApplicationUpdateParamsBodySelfHostedApplication{
			Type: cloudflare.F(zero_trust.ApplicationTypeSelfHosted),
		}
		if cmd.Flags().Changed("name") {
			body.Name = cloudflare.F(updateName)
		}
		if cmd.Flags().Changed("domain") {
			body.Domain = cloudflare.F(updateDomain)
		}

		params := zero_trust.AccessApplicationUpdateParams{
			AccountID: cloudflare.F(accountID),
			Body:      body,
		}

		app, err := client.ZeroTrust.Access.Applications.Update(context.TODO(), args[0], params)
		if err != nil {
			return fmt.Errorf("failed to update access application: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully updated Access Application: %s\n", args[0])
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(app)
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateName, "name", "", "Application name")
	updateCmd.Flags().StringVar(&updateDomain, "domain", "", "Application domain")
}
