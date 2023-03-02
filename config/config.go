package config

import (
	"encoding/json"
	"github.com/sflewis2970/trivia-api/common"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

const BASE_DIR_NAME string = "trivia-service"
const CFG_FILE_NAME = "./config/config.json"
const UPDATE_CONFIG_DATA string = "update"

// Config variable keys
const (
	// ENV System ENV setting
	ENV string = "ENV"

	// HOST system info
	HOST string = "HOST"
	PORT string = "PORT"

	// REDIS_TLS_URL Redis Constants
	REDIS_TLS_URL  string = "REDIS_TLS_URL"
	REDIS_URL      string = "REDIS_URL"
	REDIS_PORT     string = "REDIS_PORT"
	REDIS_PASSWORD string = "REDIS_PASSWORD"
)

// PRODUCTION Config variable values
const (
	PRODUCTION string = "PROD"
)

type Redis struct {
	TLS_URL  string `json:"tls_url"`
	URL      string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type Message struct {
	CongratsMsg string `json:"congrats"`
	TryAgainMsg string `json:"tryagain"`
}

type CfgData struct {
	Env      string `json:"env"`
	Host     string `json:"hostname"`
	Port     string `json:"hostport"`
	Messages Message
	Redis    Redis
}

type config struct {
	cfgData *CfgData
}

var cfg *config

// Unexported type functions
func (c *config) findBaseDir(currentDir string, targetDir string) int {
	level := 0
	dirs := strings.Split(currentDir, "\\")

	dirsSize := len(dirs)
	for idx := dirsSize - 1; idx >= 0; idx-- {
		if dirs[idx] == targetDir {
			break
		} else {
			level++
		}
	}

	return level
}

func (c *config) readConfigFile() error {
	// Get working directory
	wd, getErr := common.GetWorkingDir()
	if getErr != nil {
		log.Print("Error getting working directory")
		return getErr
	}

	// Find path
	levels := c.findBaseDir(wd, BASE_DIR_NAME)
	for levels > 0 {
		chErr := os.Chdir("..")
		if chErr != nil {
			log.Print("Error changing dir: ", chErr)
		}

		// Update levels
		levels--
	}

	data, readErr := ioutil.ReadFile(CFG_FILE_NAME)
	if readErr != nil {
		return readErr
	}

	unmarshalErr := json.Unmarshal(data, c.cfgData)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

func (c *config) getConfigEnv() error {
	// Loading config environment variables
	log.Print("loading config environment variables...")

	// Update config data
	// Base config settings
	c.cfgData.Env = os.Getenv(ENV)
	c.cfgData.Host = os.Getenv(HOST)
	c.cfgData.Port = os.Getenv(PORT)

	// Get response messages
	c.cfgData.Messages.CongratsMsg = "Congratulations! That is correct"
	c.cfgData.Messages.TryAgainMsg = "Nice Try! Better luck on the next answer"

	// Go-redis settings
	log.Print("Setting go-redis environment variables...")
	c.cfgData.Redis.TLS_URL = os.Getenv(REDIS_TLS_URL)
	c.cfgData.Redis.URL = os.Getenv(REDIS_URL)
	c.cfgData.Redis.Port = os.Getenv(REDIS_PORT)

	if c.cfgData.Env == PRODUCTION {
		log.Print("Loading prod settings...")
		redisURL, parseErr := url.Parse(c.cfgData.Redis.URL)
		if parseErr != nil {
			log.Print("Error parsing url: ", parseErr)
			return parseErr
		}

		// Update URL and Port after parsing
		delimiter := ":"
		if strings.Contains(redisURL.Host, delimiter) {
			urlSlice := strings.Split(redisURL.Host, delimiter)
			c.cfgData.Redis.URL = urlSlice[0]
			c.cfgData.Redis.Port = urlSlice[1]
		} else {
			c.cfgData.Redis.URL = redisURL.Host
			c.cfgData.Redis.Port = ":" + redisURL.Port()
		}

		log.Print("Redis URL: ", c.cfgData.Redis.URL)
		log.Print("Redis Port: ", c.cfgData.Redis.Port)

		// redis Password
		c.cfgData.Redis.Password, _ = redisURL.User.Password()
	} else {
		c.cfgData.Redis.URL = os.Getenv(REDIS_URL)
		c.cfgData.Redis.Port = os.Getenv(REDIS_PORT)
		c.cfgData.Redis.Password = os.Getenv(REDIS_PASSWORD)
	}

	return nil
}

func (c *config) GetData(args ...string) (*CfgData, error) {
	if len(args) > 0 {
		if args[0] == UPDATE_CONFIG_DATA {
			useCfgFile := os.Getenv("USECONFIGFILE")
			if len(useCfgFile) > 0 {
				log.Print("Using config file to load config")

				readErr := cfg.readConfigFile()
				if readErr != nil {
					log.Print("Error reading config file: ", readErr)
					return nil, readErr
				}
			} else {
				log.Print("Using config environment to load config")

				getErr := cfg.getConfigEnv()
				if getErr != nil {
					log.Print("Error getting config environment data: ", getErr)
					return nil, getErr
				}
			}
		}
	}

	return c.cfgData, nil
}

func Get() *config {
	if cfg == nil {
		log.Print("creating config object")

		// Initialize config
		cfg = new(config)

		// Initialize config data
		cfg.cfgData = new(CfgData)
	} else {
		log.Print("returning config object")
	}

	return cfg
}
