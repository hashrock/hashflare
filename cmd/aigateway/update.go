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
	updateCollectLogs bool
	updateCacheTTL    int64
	updateRateLimit   int64
	updateRateInterval int64
)

var updateCmd = &cobra.Command{
	Use:   "update <gateway-id>",
	Short: "Update an AI Gateway",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		params := ai_gateway.AIGatewayUpdateParams{
			AccountID: cloudflare.F(accountID),
		}

		if cmd.Flags().Changed("collect-logs") {
			params.CollectLogs = cloudflare.F(updateCollectLogs)
		}
		if cmd.Flags().Changed("cache-ttl") {
			params.CacheTTL = cloudflare.F(updateCacheTTL)
		}
		if cmd.Flags().Changed("rate-limit") {
			params.RateLimitingLimit = cloudflare.F(updateRateLimit)
		}
		if cmd.Flags().Changed("rate-interval") {
			params.RateLimitingInterval = cloudflare.F(updateRateInterval)
		}

		gw, err := client.AIGateway.Update(context.TODO(), args[0], params)
		if err != nil {
			return fmt.Errorf("failed to update AI gateway: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully updated AI Gateway: %s\n", args[0])
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(gw)
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateCollectLogs, "collect-logs", false, "Enable request/response logging")
	updateCmd.Flags().Int64Var(&updateCacheTTL, "cache-ttl", 0, "Cache TTL in seconds")
	updateCmd.Flags().Int64Var(&updateRateLimit, "rate-limit", 0, "Rate limiting max requests")
	updateCmd.Flags().Int64Var(&updateRateInterval, "rate-interval", 0, "Rate limiting interval in seconds")
}
