package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const FileName = "hashflare.json"

type Config struct {
	AIGateway      string   `json:"ai_gateway,omitempty"`
	AccessApp      string   `json:"access_app,omitempty"`
	AccessPolicies []string `json:"access_policies,omitempty"`
}

func FindPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, FileName), nil
}

func Load() (*Config, error) {
	p, err := FindPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Save(c *Config) error {
	p, err := FindPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, append(data, '\n'), 0644)
}
