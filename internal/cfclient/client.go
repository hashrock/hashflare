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

func SettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hashflare", "setting.json"), nil
}

func LoadSettings() *Settings {
	p, err := SettingsPath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	return &s
}

func SaveSettings(s *Settings) error {
	p, err := SettingsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, append(data, '\n'), 0600)
}

func NewClient() *cloudflare.Client {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		if s := LoadSettings(); s != nil {
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
		if s := LoadSettings(); s != nil {
			id = s.AccountID
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "Error: CLOUDFLARE_ACCOUNT_ID is not set. Set the env var or configure ~/.hashflare/setting.json")
		os.Exit(1)
	}
	return id
}
