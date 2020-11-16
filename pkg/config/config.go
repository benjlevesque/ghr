package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ConfigItem struct {
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
	Executable string `json:"executable,omitempty"`
}

type Config struct {
	Items []ConfigItem `json:"items,omitempty"`
}

func (c *Config) Get(repoName string) *ConfigItem {
	for _, item := range c.Items {
		if item.Name == repoName {
			return &item
		}
	}
	return nil
}

// Owner returns the owner part of the item's Name
func (i *ConfigItem) Owner() string {
	return strings.Split(i.Name, "/")[0]
}

// Repository returns the repo part of the item's Name
func (i *ConfigItem) Repository() string {
	return strings.Split(i.Name, "/")[1]
}

func Load() (*Config, error) {
	path := getConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		configDir := filepath.Dir(path)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return nil, err
		}
		err := ioutil.WriteFile(path, []byte(""), 0644)
		if err != nil {
			return nil, err
		}

	}
	var config Config

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config: %s", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func getConfigPath() string {
	path := viper.GetString("config")
	sudoUser := os.Getenv("SUDO_USER")
	var home string
	var err error
	if sudoUser == "" {
		home, err = os.UserHomeDir()
	} else {
		var u *user.User
		u, err = user.Lookup(sudoUser)
		if err == nil {
			home = u.HomeDir
		}
	}
	if err == nil {
		path = strings.ReplaceAll(path, "$HOME", home)
	}
	return path
}

func Save(config *Config) error {
	path := getConfigPath()
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func AddOrUpdate(item ConfigItem) error {
	config, err := Load()
	if err != nil {
		return err
	}
	return config.AddOrUpdate(item)
}

func (c *Config) AddOrUpdate(item ConfigItem) error {
	for i, currentItem := range c.Items {
		if currentItem.Name == item.Name {
			c.Items[i] = item
			return Save(c)
		}
	}
	c.Items = append(c.Items, item)
	return Save(c)
}
