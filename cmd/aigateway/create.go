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

var (
	createCollectLogs bool
	createCacheTTL    int64
)

var createCmd = &cobra.Command{
	Use:   "create <gateway-id>",
	Short: "Create a new AI Gateway",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		params := ai_gateway.AIGatewayNewParams{
			AccountID: cloudflare.F(accountID),
			ID:        cloudflare.F(args[0]),
		}

		if cmd.Flags().Changed("collect-logs") {
			params.CollectLogs = cloudflare.F(createCollectLogs)
		}
		if cmd.Flags().Changed("cache-ttl") {
			params.CacheTTL = cloudflare.F(createCacheTTL)
		}

		gw, err := client.AIGateway.New(context.TODO(), params)
		if err != nil {
			return fmt.Errorf("failed to create AI gateway: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully created AI Gateway: %s\n", args[0])
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(gw)
	},
}

func init() {
	createCmd.Flags().BoolVar(&createCollectLogs, "collect-logs", false, "Enable request/response logging")
	createCmd.Flags().Int64Var(&createCacheTTL, "cache-ttl", 0, "Cache TTL in seconds")
}
