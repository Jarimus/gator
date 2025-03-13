package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

// Read the .gatorconfig.json file and returns it in Config struct format.
func Read() (Config, error) {

	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// Check if the config file exists. If not, create an empty config file
	_, err = os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		newConfig := Config{DbUrl: "postgres://example"}
		newFile, err := json.Marshal(newConfig)
		if err != nil {
			return Config{}, err
		}
		os.WriteFile(filepath, newFile, 0644)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshaling config file: %s", err)
	}
	return config, nil
}

// Sets the current user in the config struct and writes to config file
func (c *Config) SetUser(currentUser string) error {
	c.CurrentUser = currentUser
	err := writeConfig(c)
	if err != nil {
		return err
	}
	return nil
}

// Gets the file path to the config file.
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filepath := homePath + "/" + configFileName

	return filepath, nil
}

// Writes the config struct into the config file
func writeConfig(c *Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
