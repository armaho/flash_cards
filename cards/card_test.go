package cards_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/armaho/flash_cards/cards"
	"github.com/armaho/flash_cards/config"
)

func checkCardsAreTheSame(c *cards.Card, expected *cards.Card, t *testing.T) {
	if (c == nil) && (expected == nil) {
		return
	}

	if (c == nil) || (expected == nil) {
		t.Errorf("expected card %#v, got %#v", expected, c)
	}

	cData, err := json.Marshal(c)
	if err != nil {
		t.Errorf("json.Marshal error: %s", err)
	}
	expectedCardData, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("json.Marshal error: %s", err)
	}
	if string(cData) != string(expectedCardData) {
		t.Errorf("expected card %#v, got %#v", expected, c)
	}
}

func resetCardData(t *testing.T) {
	cfg := config.LoadOrDie()
	path := filepath.Join(cfg.DataPath, "cards")

	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("Failed to reset card data: %v", err)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to recreate cards directory: %v", err)
	}
}

func makeCardIneligible(c *cards.Card) {
	now := time.Now()
	c.LastAsked = &now
	c.Interval = 24
}

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		panic("cannot create temporary direcory for config: " + err.Error())
	}
	configPath := filepath.Join(tempDir, "test_config.json")

	os.Setenv("CONFIG_PATH", configPath)

	cfg := config.GetDefaultConfig()
	cfg.InitialInterval = 1
	cfg.IncreaseFactor = 2
	cfg.DecreaseFactor = 0.5
	cfg.DataPath = filepath.Join(tempDir, "data")
	config.Save(&cfg)

	exitVal := m.Run()

	os.RemoveAll(tempDir)
	os.Exit(exitVal)
}

func TestEveryNewCardShouldBeEligibleTwice(t *testing.T) {
	c := cards.NewCard("question", "answer")

	if !c.IsEligibleToAsk() {
		t.Errorf("newly created card is not eligible to ask: %#v", c)
	}
	c.Upgrade()
	if !c.IsEligibleToAsk() {
		t.Errorf("newly created card is not eligible to ask after answering once: %#v", c)
	}
}

func TestFindCardToAskShouldReturnEligibleCards(t *testing.T) {
	cardList := []*cards.Card{
		cards.NewCard("q1", "a1"),
		cards.NewCard("q2", "a2"),
	}

	makeCardIneligible(cardList[0])
	c := cards.FindCardToAsk(cardList)
	if c == nil || !c.IsEligibleToAsk() {
		t.Fatalf("cannot find eligible card")
	}

	cardList = cardList[:1]
	c = cards.FindCardToAsk(cardList)
	if c != nil {
		t.Fatalf("found eligible card but there weren't any")
	}
}

func TestNewCardsShouldNotBeEligibleTheThirdTimeImmediately(t *testing.T) {
	c := cards.NewCard("question", "answer")

	c.Upgrade()
	c.Upgrade()

	if c.IsEligibleToAsk() {
		t.Errorf("newly created card has become eligible immediately for the third time: %#v", c)
	}
}

func TestSaveCardIdAssignment(t *testing.T) {
	c := cards.NewCard("question", "answer")

	err := cards.SaveCard(c)
	if err != nil {
		t.Errorf("error while saving card: %s", err)
	}

	if c.Id == "" {
		t.Errorf("saved card does not have an id")
	}
}

func TestSaveAndFetchCard(t *testing.T) {
	c := cards.NewCard("question", "answer")

	err := cards.SaveCard(c)
	if err != nil {
		t.Errorf("error while saving card: %s", err)
	}

	if c.Id == "" {
		t.Errorf("saved card does not have an id")
	}

	c2, err := cards.FetchCard(c.Id)
	if err != nil {
		t.Errorf("cannot read card: %s", err)
	}

	checkCardsAreTheSame(c2, c, t)
}

func TestFetchAllCards(t *testing.T) {
	resetCardData(t)

	cardList := []*cards.Card{
		cards.NewCard("q1", "a1"),
		cards.NewCard("q2", "a2"),
		cards.NewCard("q3", "a3"),
	}

	for _, c := range cardList {
		err := cards.SaveCard(c)
		if err != nil {
			t.Fatalf("error while saving card: %s", err)
		}
	}

	fetchedCardList, err := cards.FetchAllCards()
	if err != nil {
		t.Fatalf("cannot fetch cards: %v", err)
	}

	if len(fetchedCardList) != len(cardList) {
		t.Fatalf("not all cards are fetched. current count: %d, expected count: %d", len(fetchedCardList), len(cardList))
	}

	fetchedMap := make(map[string]*cards.Card)
	for _, card := range fetchedCardList {
		fetchedMap[card.Id] = card
	}
	for _, expectedCard := range cardList {
		fetchedCard, exists := fetchedMap[expectedCard.Id]
		if !exists {
			t.Errorf("card with Id %s not found in fetched results", expectedCard.Id)
			continue
		}
		checkCardsAreTheSame(fetchedCard, expectedCard, t)
	}
}

func TestCannotFetchCardThatDoesNotExist(t *testing.T) {
	id := "does-not-exist"
	_, err := cards.FetchCard(id)
	if err != cards.ErrCardNotFound {
		t.Fatalf("fetched card that does not exist")
	}
}

func TestDeleteCard(t *testing.T) {
	c := cards.NewCard("question", "answer")

	err := cards.SaveCard(c)
	if err != nil {
		t.Fatalf("error while saving card: %v", err)
	}

	if c.Id == "" {
		t.Fatalf("saved card does not have an id")
	}

	_, err = cards.FetchCard(c.Id)
	if err != nil {
		t.Fatalf("cannot read card after saving: %v", err)
	}

	err = cards.DeleteCard(c.Id)
	if err != nil {
		t.Fatalf("cannot delete card after saving: %v", err)
	}

	_, err = cards.FetchCard(c.Id)
	if err != cards.ErrCardNotFound {
		t.Fatalf("read deleted card")
	}
}
