package cfclient

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

type Settings struct {
	APIToken  string `json:"api_token"`
	AccountID string `json:"account_id"`
}

func loadSettings() *Settings {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(home, ".hashflare", "setting.json"))
	if err != nil {
		return nil
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	return &s
}

func NewClient() *cloudflare.Client {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		if s := loadSettings(); s != nil {
			token = s.APIToken
		}
	}
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: CLOUDFLARE_API_TOKEN is not set. Set the env var or configure ~/.hashflare/setting.json")
		os.Exit(1)
	}
	return cloudflare.NewClient(option.WithAPIToken(token))
}

func GetAccountID() string {
	id := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	if id == "" {
		if s := loadSettings(); s != nil {
			id = s.AccountID
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "Error: CLOUDFLARE_ACCOUNT_ID is not set. Set the env var or configure ~/.hashflare/setting.json")
		os.Exit(1)
	}
	return id
}
