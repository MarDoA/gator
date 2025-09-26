package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	Url  string `json:"db_url"`
	User string `json:"current_user_name"`
}

func (cfg *Config) SetUser(user string) error {
	cfg.User = user
	return write(*cfg)

}

func Read() (Config, error) {
	dir, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(dir)
	if err != nil {
		return Config{}, err
	}
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir + configFileName, nil
}

func write(cfg Config) error {
	dir, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}
