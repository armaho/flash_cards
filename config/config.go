package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	IncreaseFactor  float64 `json:"increase_factor"`
	DecreaseFactor  float64 `json:"decrease_factor"`
	InitialInterval int     `json:"initial_interval"`
}

func GetDefaultConfig() Config {
	return Config{
		IncreaseFactor:  2,
		DecreaseFactor:  0.5,
		InitialInterval: 24,
	}
}

func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, errors.New("CONFIG_PATH is not provided")
	}

	cfg := Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		cfg = GetDefaultConfig()

		err = Save(&cfg)
		if err != nil {
			return nil, err
		}

		return &cfg, nil
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

func Save(cfg *Config) error {
	if cfg == nil {
		return errors.New("Cannot save null as config")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return errors.New("CONFIG_PATH is not provided")
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
