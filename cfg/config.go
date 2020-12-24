package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config for gf-metrics
type Config struct {
	Token               string   `json:"token"`
	Team                []string `json:"team"`
	RepoOwner           string   `json:"owner"`
	MetricsStartingDate string   `json:"timestamp"`
}

// FromFile creates a Config from a JSON file
func FromFile(path string) (*Config, error) {
	cfg := &Config{}
	jsonFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	bytesFile, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(bytesFile, cfg)

	return cfg, err
}

// FillEmpty fills any empty Config field and update any empty (not nil) value passed as a parameter
func (c *Config) FillEmpty(token, owner *string, team *[]string, startDate *string) {
	if c.Token == "" {
		c.Token = *token
	}

	if c.RepoOwner == "" {
		c.RepoOwner = *owner
	}

	if len(c.Team) == 0 {
		c.Team = *team
	}

	if len(c.MetricsStartingDate) == 0 {
		c.MetricsStartingDate = *startDate
	}

	*token = c.Token
	*owner = c.RepoOwner
	*team = c.Team
	*startDate = c.MetricsStartingDate
}
