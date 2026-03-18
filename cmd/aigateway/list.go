package aigateway

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/ai_gateway"
	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all AI Gateways",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cfclient.NewClient()
		accountID := cfclient.GetAccountID()

		page, err := client.AIGateway.List(context.TODO(), ai_gateway.AIGatewayListParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("failed to list AI gateways: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCOLLECT LOGS\tCACHE TTL\tCREATED")
		for _, gw := range page.Result {
			fmt.Fprintf(w, "%s\t%v\t%v\t%s\n",
				gw.ID,
				gw.CollectLogs,
				gw.CacheTTL,
				gw.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
		w.Flush()
		return nil
	},
}
