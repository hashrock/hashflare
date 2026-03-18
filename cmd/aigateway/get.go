package aigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/ai_gateway"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <gateway-id>",
	Short: "Get details of an AI Gateway",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		gw, err := client.AIGateway.Get(context.TODO(), args[0], ai_gateway.AIGatewayGetParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to get AI gateway: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(gw)
	},
}
