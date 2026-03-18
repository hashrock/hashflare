package access

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <app-id>",
	Short: "Delete an Access Application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		_, err := client.ZeroTrust.Access.Applications.Delete(context.TODO(), args[0], zero_trust.AccessApplicationDeleteParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to delete access application: %w", err)
		}

		fmt.Printf("Successfully deleted Access Application: %s\n", args[0])
		return nil
	},
}
