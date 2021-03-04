package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const configDirName = ".ddboard"
const configFileName = "config.json"

type Gitlab struct {
	Url      string   `json:"url"`
	Token    string   `json:"token"`
	UserId   int      `json:"user_id"`
	Projects []string `json:"projects"`
}

type Data struct {
	*Gitlab
}

type Config struct {
	ConfigFile string
	ConfigDir  string
	isLoaded   bool
	*Data
}

func New(configDir string) *Config {
	sep := string(os.PathSeparator)
	configFileDir := configDir + sep + configDirName
	configFileName := configFileDir + sep + configFileName

	conf := &Config{
		ConfigFile: configFileName,
		ConfigDir:  configFileDir,
		isLoaded:   false,
		Data: &Data{
			Gitlab: &Gitlab{
				Token:    "",
				UserId:   0,
				Projects: []string{},
			},
		},
	}

	return conf
}

func (c *Config) load() error {
	if !c.isLoaded {
		_ = c.createConfigDirectory()
		configFile, err := os.Open(c.ConfigFile)

		_, err = os.Stat(c.ConfigFile)
		if os.IsNotExist(err) {
			_, err = os.Create(c.ConfigFile)
			if err != nil {
				return fmt.Errorf("could not create config file: %w", err)
			}
		}

		if err != nil {
			return fmt.Errorf("could not open config file: %w", err)
		}
		defer configFile.Close()

		byteValue, err := ioutil.ReadAll(configFile)
		if err != nil {
			return fmt.Errorf("could not read from config file: %w", err)
		}

		var data Data
		err = json.Unmarshal(byteValue, &data)
		if err != nil {
			return fmt.Errorf("could not decode json from config file: %w", err)
		}

		c.Data = &data
		c.isLoaded = true
	}

	return nil
}

func (c *Config) IsReady() bool {
	err := c.load()
	if err != nil {
		return false
	}

	// todo - better ready validation
	if c.Data.Gitlab.Token == "" {
		return false
	}

	return true
}

func (c *Config) Write() error {
	data, err := json.Marshal(c.Data)
	_ = os.Truncate(c.ConfigFile, 0)
	configFile, err := os.OpenFile(c.ConfigFile, os.O_WRONLY, 0700) // todo permissions
	if err != nil {
		return fmt.Errorf("could not open config file for write: %w", err)
	}

	_, err = configFile.Write(data)
	if err != nil {
		return fmt.Errorf("could not write changes into config file: %w", err)
	}

	return nil
}

func (c *Config) createConfigDirectory() error {
	return os.MkdirAll(c.ConfigDir, 0700) // todo permissions
}
