package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func buildURL(scheme, address string) string {
	return fmt.Sprintf("%s://%s", scheme, address)
}

// Config describes the applications settings
type Config struct {
	Feed struct {
		Address string `envconfig:"FEED_ADDR" default:"localhost:8000" desc:"Address of the feed"`
		Scheme  string `envconfig:"FEED_SCHEME" default:"http" desc:"Scheme of the communication for the feed api"`
	}

	Store struct {
		Address string `envconfig:"STORE_ADDR" default:"localhost:8001" desc:"Address of the store"`
		Scheme  string `envconfig:"STORE_SCHEME" default:"http" desc:"Scheme of the communication for the store api"`
	}
}

// NewConfig is for creating a new instance of Config
func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("importer", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetStoreURL returns correct URL from scheme and address variables
func (c Config) GetStoreURL() string {
	return buildURL(c.Store.Scheme, c.Store.Address)
}

// GetFeedURL returns correct URL from scheme and address variables
func (c Config) GetFeedURL() string {
	return buildURL(c.Feed.Scheme, c.Feed.Address)
}
