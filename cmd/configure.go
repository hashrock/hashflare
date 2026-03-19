package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hashrock/hashflare/internal/cfclient"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure API token and account ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		existing := cfclient.LoadSettings()
		if existing == nil {
			existing = &cfclient.Settings{}
		}

		reader := bufio.NewReader(os.Stdin)

		token, err := prompt(reader, "API Token", existing.APIToken)
		if err != nil {
			return err
		}
		accountID, err := prompt(reader, "Account ID", existing.AccountID)
		if err != nil {
			return err
		}

		s := &cfclient.Settings{
			APIToken:  token,
			AccountID: accountID,
		}
		if err := cfclient.SaveSettings(s); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		p, _ := cfclient.SettingsPath()
		fmt.Printf("Settings saved to %s\n", p)
		return nil
	},
}

func prompt(reader *bufio.Reader, label, current string) (string, error) {
	if current != "" {
		masked := mask(current)
		fmt.Printf("%s [%s]: ", label, masked)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return current, nil
	}
	return line, nil
}

func mask(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
