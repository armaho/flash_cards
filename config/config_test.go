package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/armaho/flash_cards/config"
)

func checkFileExistance(path string, t *testing.T) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("File not found at %s", path)
	}
}

func checkFileContent(path, expected string, t *testing.T) {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected content: %s\nGot: %s", expected, string(data))
	}
}

func checkConfig(cfg, expected *config.Config, t *testing.T) {
	if (cfg == nil) && (expected == nil) {
		return
	}

	if (cfg == nil) || (expected == nil) ||
		(cfg.IncreaseFactor != expected.IncreaseFactor) ||
		(cfg.DecreaseFactor != expected.DecreaseFactor) ||
		(cfg.InitialInterval != expected.InitialInterval) {
		t.Errorf("expected config %#v, got %#v", expected, cfg)
	}
}

func TestSaveShouldReturnAnErrorForNilConfig(t *testing.T) {
	t.Setenv("CONFIG_PATH", filepath.Join(".", "config.json"))

	err := config.Save(nil)
	if err == nil {
		t.Error("Expected config.Save to return an error if the config is nil")
	}
}

func TestSaveShouldReturnAnErrorIfConfigPathHasNotBeenSet(t *testing.T) {
	t.Setenv("CONFIG_PATH", "")

	err := config.Save(&config.Config{
		DecreaseFactor:  0.5,
		IncreaseFactor:  2,
		InitialInterval: 24,
	})
	if err == nil {
		t.Error("Expected config.Save to return an error if CONFIG_PATH hasn't been set")
	}
}

func TestSaveShouldCreateANewFileIfConfigDoesNotExist(t *testing.T) {
	configPath := filepath.Join(t.TempDir(),
		"ThisDirDoesNotExist",
		"ThisFileDoesNotExist.json")

	t.Setenv("CONFIG_PATH", configPath)

	cfg := config.Config{
		DecreaseFactor:  0.5,
		IncreaseFactor:  2,
		InitialInterval: 25,
	}

	err := config.Save(&cfg)
	if err != nil {
		t.Errorf("config.Save throw an error: %s", err)
	}

	checkFileExistance(configPath, t)

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Errorf("json.Marshal error: %s", err)
	}
	checkFileContent(configPath, string(data), t)
}

func TestLoadShouldCreateANewFileIfConfigDoesNotExist(t *testing.T) {
	configPath := filepath.Join(t.TempDir(),
		"ThisDirDoesNotExist",
		"ThisFileDoesNotExist.json")

	t.Setenv("CONFIG_PATH", configPath)

	cfg, err := config.Load()
	if err != nil {
		t.Errorf("config.Load throw an error: %s", err)
		return
	}

	checkFileExistance(configPath, t)
	checkFileContent(configPath, `{"increase_factor":2,"decrease_factor":0.5,"initial_interval":24}`, t)
	checkConfig(cfg, &config.Config{
		DecreaseFactor:  0.5,
		IncreaseFactor:  2,
		InitialInterval: 24,
	}, t)
}

func TestLoadShouldReadCorrectly(t *testing.T) {
	configPath := filepath.Join(t.TempDir(),
		"ThisDirDoesNotExist",
		"ThisFileDoesNotExist.json")

	t.Setenv("CONFIG_PATH", configPath)

	expected_config := &config.Config{
		DecreaseFactor:  0.9,
		IncreaseFactor:  1.1,
		InitialInterval: 27,
	}

	err := config.Save(expected_config)
	if err != nil {
		t.Errorf("config.Save throw an error: %s", err)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		t.Errorf("config.Load throw an error: %s", err)
		return
	}

	checkConfig(cfg, expected_config, t)
}
