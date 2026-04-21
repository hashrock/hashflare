package aigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/ai_gateway"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/hashrock/hashflare/internal/project"
	"github.com/spf13/cobra"
)

var (
	createCollectLogs            bool
	createCacheTTL               int64
	createCacheInvalidateOnUpdate bool
	createRateLimit              int64
	createRateInterval           int64
	createRateTechnique          string
	createAuthentication         bool
)

var createCmd = &cobra.Command{
	Use:   "create <gateway-id>",
	Short: "Create a new AI Gateway",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		technique := ai_gateway.AIGatewayNewParamsRateLimitingTechnique(createRateTechnique)
		if !technique.IsKnown() {
			return fmt.Errorf("invalid --rate-technique %q: must be one of fixed, sliding", createRateTechnique)
		}

		params := ai_gateway.AIGatewayNewParams{
			AccountID:               cloudflare.F(accountID),
			ID:                      cloudflare.F(args[0]),
			CacheInvalidateOnUpdate: cloudflare.F(createCacheInvalidateOnUpdate),
			CacheTTL:                cloudflare.F(createCacheTTL),
			CollectLogs:             cloudflare.F(createCollectLogs),
			RateLimitingInterval:    cloudflare.F(createRateInterval),
			RateLimitingLimit:       cloudflare.F(createRateLimit),
			RateLimitingTechnique:   cloudflare.F(technique),
		}

		if cmd.Flags().Changed("authentication") {
			params.Authentication = cloudflare.F(createAuthentication)
		}

		gw, err := client.AIGateway.New(context.TODO(), params)
		if err != nil {
			return fmt.Errorf("failed to create AI gateway: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully created AI Gateway: %s\n", args[0])

		if c, err := project.Load(); err == nil {
			c.AIGateway = args[0]
			if err := project.Save(c); err == nil {
				fmt.Fprintf(os.Stderr, "Linked to project: %s\n", project.FileName)
			}
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(gw)
	},
}

func init() {
	createCmd.Flags().BoolVar(&createCollectLogs, "collect-logs", false, "Enable request/response logging")
	createCmd.Flags().Int64Var(&createCacheTTL, "cache-ttl", 0, "Cache TTL in seconds (0 = disabled)")
	createCmd.Flags().BoolVar(&createCacheInvalidateOnUpdate, "cache-invalidate-on-update", false, "Invalidate cache on update")
	createCmd.Flags().Int64Var(&createRateLimit, "rate-limit", 0, "Rate limiting max requests (0 = disabled)")
	createCmd.Flags().Int64Var(&createRateInterval, "rate-interval", 0, "Rate limiting interval in seconds (0 = disabled)")
	createCmd.Flags().StringVar(&createRateTechnique, "rate-technique", "fixed", "Rate limiting technique: fixed | sliding")
	createCmd.Flags().BoolVar(&createAuthentication, "authentication", false, "Require authentication for gateway requests")
}
