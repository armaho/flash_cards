package cards

import (
	"math"
	"math/rand"
	"time"

	"github.com/armaho/flash_cards/config"
)

type Card struct {
	Id        string     `json:"id"`
	Question  string     `json:"question"`
	Answer    string     `json:"answer"`
	LastAsked *time.Time `json:"last_asked"`
	Interval  int        `json:"interval"`
}

func NewCard(question, answer string) *Card {
	return &Card{
		Id:        "",
		Question:  question,
		Answer:    answer,
		LastAsked: nil,
		Interval:  0,
	}
}

func (c *Card) IsEligibleToAsk() bool {
	if c.LastAsked == nil {
		return true
	}

	interval := time.Duration(c.Interval) * time.Hour
	nextTimeToAsk := c.LastAsked.Add(interval)
	now := time.Now()
	return !now.Before(nextTimeToAsk)
}

func FindCardToAsk(cards []*Card) *Card {
	eligibleCards := []*Card{}
	for _, c := range cards {
		if c.IsEligibleToAsk() {
			eligibleCards = append(eligibleCards, c)
		}
	}

	if len(eligibleCards) == 0 {
		return nil
	}

	randIdx := rand.Intn(len(eligibleCards))
	return eligibleCards[randIdx]
}

func (c *Card) Upgrade() {
	cfg := config.LoadOrDie()

	if c.LastAsked == nil {
		c.Interval = 0
	} else if c.Interval == 0 {
		c.Interval = cfg.InitialInterval
	} else {
		c.Interval = int(math.Floor(float64(c.Interval) * cfg.IncreaseFactor))
	}

	now := time.Now()
	c.LastAsked = &now
}

func (c *Card) Downgrade() {
	cfg := config.LoadOrDie()
	c.Interval = int(math.Floor(float64(c.Interval)) * cfg.DecreaseFactor)
}
