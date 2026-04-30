package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/armaho/flash_cards/cards"
)

func readNonEmptyInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	var input string

	for input == "" {
		fmt.Print(prompt)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("input cannot be empty")
		}
	}

	return input
}

func addNewCard() {
	q := readNonEmptyInput("question: ")
	a := readNonEmptyInput("answer: ")

	c := cards.NewCard(q, a)
	err := cards.SaveCard(c)
	if err != nil {
		fmt.Printf("cannot save card: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("successfully saved. id: %s\n", c.Id)
}

func addVocabCard() {
	word := readNonEmptyInput("word: ")
	def := readNonEmptyInput("definition: ")
	example := readNonEmptyInput("example: ")

	cardList := []*cards.Card{
		cards.NewCard(fmt.Sprintf("What's the meaning of \"%s\"?", word), def),
		cards.NewCard(fmt.Sprintf("Use \"%s\" in a sentence", word), example),
		cards.NewCard(fmt.Sprintf("Say \"%s\" out loud", word), word),
	}

	for _, c := range cardList {
		if err := cards.SaveCard(c); err != nil {
			fmt.Printf("cannot save card: %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("successfully saved all cards")
}

func askCard() {
	cardList, err := cards.FetchAllCards()
	if err != nil {
		fmt.Printf("cannot fetch cards: %s\n", err)
		os.Exit(1)
	}

	for {
		c := cards.FindCardToAsk(cardList)
		if c == nil {
			fmt.Println("you have answered every card!\ntry again in another time.")
			return
		}

		fmt.Printf("%s?\n\n", c.Question)

		fmt.Println("press enter when you answered")
		fmt.Scanln()

		fmt.Printf("answer: %s\n\n", c.Answer)

		var ans string
		for ans != "y" && ans != "n" {
			fmt.Printf("have you answered correctly? [y/n]")
			fmt.Scan(&ans)

			ans = strings.TrimSpace(ans)
			ans = strings.ToLower(ans)

			if ans != "y" && ans != "n" {
				fmt.Println("invalid option")
			}
		}

		switch ans {
		case "y":
			c.Upgrade()
		case "n":
			c.Downgrade()
		}

		err = cards.SaveCard(c)
		if err != nil {
			fmt.Printf("cannot save card: %s", err)
			os.Exit(1)
		}

		fmt.Println()
	}
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		fmt.Println("usage: cards [command]")
		fmt.Println("commands:")
		fmt.Println("\tadd [eng]   add a new card (with eng, if you're adding a vocab card)")
		fmt.Println("\task         asking cards on a loop")

		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) >= 3 && os.Args[2] == "eng" {
			addVocabCard()
			break
		}
		addNewCard()
	case "ask":
		askCard()
	}
}
