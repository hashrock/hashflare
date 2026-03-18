package aigateway

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/ai_gateway"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <gateway-id>",
	Short: "Delete an AI Gateway",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		_, err := client.AIGateway.Delete(context.TODO(), args[0], ai_gateway.AIGatewayDeleteParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to delete AI gateway: %w", err)
		}

		fmt.Printf("Successfully deleted AI Gateway: %s\n", args[0])
		return nil
	},
}
