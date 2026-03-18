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
	createName            string
	createDomain          string
	createSessionDuration string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Access Application",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		params := zero_trust.AccessApplicationNewParams{
			AccountID: cloudflare.F(accountID),
			Body: zero_trust.AccessApplicationNewParamsBodySelfHostedApplication{
				Name:   cloudflare.F(createName),
				Domain: cloudflare.F(createDomain),
				Type:   cloudflare.F(zero_trust.ApplicationTypeSelfHosted),
			},
		}

		app, err := client.ZeroTrust.Access.Applications.New(context.TODO(), params)
		if err != nil {
			return fmt.Errorf("failed to create access application: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully created Access Application: %s\n", createName)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(app)
	},
}

func init() {
	createCmd.Flags().StringVar(&createName, "name", "", "Application name (required)")
	createCmd.Flags().StringVar(&createDomain, "domain", "", "Application domain (required)")
	createCmd.Flags().StringVar(&createSessionDuration, "session-duration", "", "Session duration (e.g., 24h)")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("domain")
}
