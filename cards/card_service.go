package cards

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/armaho/flash_cards/config"
	"github.com/armaho/flash_cards/uuid"
)

var ErrCardNotFound = errors.New("cannot find card")

func getDataPath() string {
	cfg := config.LoadOrDie()
	return filepath.Join(cfg.DataPath, "cards")
}

func getDataPathById(id string) string {
	return filepath.Join(getDataPath(), id+".json")
}

func doesCardExist(id string) (bool, error) {
	path := getDataPathById(id)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking if card exists: %w", err)
	}

	return true, nil
}

func createCardFile(id string) (*os.File, error) {
	path := getDataPathById(id)

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return os.Create(path)
}

func SaveCard(c *Card) error {
	if c == nil {
		return errors.New("Cannot save nil card")
	}

	if c.Id == "" {
		c.Id = uuid.GenerateUUID()
	}

	file, err := createCardFile(c.Id)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func FetchCard(id string) (*Card, error) {
	path := getDataPathById(id)

	cardBytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrCardNotFound
	} else if err != nil {
		return nil, err
	}

	c := &Card{}
	err = json.Unmarshal(cardBytes, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func FetchAllCards() (cards []*Card, err error) {
	dir := getDataPath()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	cards = make([]*Card, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", filePath, err)
		}

		var card Card
		if err := json.Unmarshal(data, &card); err != nil {
			return nil, fmt.Errorf("parsing JSON from %s: %w", filePath, err)
		}

		cards = append(cards, &card)
	}

	return cards, nil
}

func DeleteCard(id string) error {
	path := getDataPathById(id)
	exists, err := doesCardExist(id)
	if !exists {
		return ErrCardNotFound
	}
	if err != nil {
		return fmt.Errorf("checking card file: %w", err)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting card file: %w", err)
	}

	return nil
}
