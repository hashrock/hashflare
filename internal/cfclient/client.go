package cfclient

import (
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

func NewClient() *cloudflare.Client {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: CLOUDFLARE_API_TOKEN environment variable is required")
		os.Exit(1)
	}
	return cloudflare.NewClient(option.WithAPIToken(token))
}

func GetAccountID() string {
	id := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	if id == "" {
		fmt.Fprintln(os.Stderr, "Error: CLOUDFLARE_ACCOUNT_ID environment variable is required")
		os.Exit(1)
	}
	return id
}
