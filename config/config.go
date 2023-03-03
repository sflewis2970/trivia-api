package config

import (
	"log"
	"os"
)

// Config variable keys
const (
	// ENV System ENV setting
	ENV string = "ENV"

	// HOST system info
	HOST string = "HOST"
	PORT string = "PORT"

	// Redis server settings
	REDIS_TLS_URL string = "REDIS_TLS_URL"
	REDIS_URL     string = "REDIS_URL"
	REDIS_PORT    string = "REDIS_PORT"
)

// PRODUCTION Config variable values
const (
	PRODUCTION string = "PROD"
)

type CfgData struct {
	Env         string `json:"env"`
	Host        string `json:"hostname"`
	Port        string `json:"hostport"`
	RedisTLSURL string `json:"redistlsurl"`
	RedisURL    string `json:"redisurl"`
	RedisPort   string `json:"redisport"`
}

type Config struct {
	cfgData *CfgData
}

var config *Config

// Unexported type functions
func (c *Config) loadConfigEnv() {
	// Loading config environment variables
	log.Print("loading config environment variables...")

	// Load host config data
	c.cfgData.Env = os.Getenv(ENV)
	c.cfgData.Host = os.Getenv(HOST)
	c.cfgData.Port = os.Getenv(PORT)

	// Load redis config data
	c.cfgData.RedisTLSURL = os.Getenv(REDIS_TLS_URL)
	c.cfgData.RedisURL = os.Getenv(REDIS_URL)
	c.cfgData.RedisPort = os.Getenv(REDIS_PORT)
}

func (c *Config) LoadCfgData() *CfgData {
	log.Print("Using config environment to load config")

	c.loadConfigEnv()

	return c.cfgData
}

func NewConfig() *Config {
	if config == nil {
		log.Print("creating config object")

		// Initialize config
		config = new(Config)

		// Initialize config data
		config.cfgData = new(CfgData)
	} else {
		log.Print("returning config object")
	}

	return config
}
