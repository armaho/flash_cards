package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/armaho/flash_cards/cards"
)

func addNewCard() {
	var q, a string

	for q == "" {
		fmt.Print("question: ")
		fmt.Scan(&q)
		q = strings.TrimSpace(q)

		if q == "" {
			fmt.Println("enter the question")
		}
	}

	for a == "" {
		fmt.Print("answer: ")
		fmt.Scan(&a)
		a = strings.TrimSpace(a)

		if a == "" {
			fmt.Println("enter the answer")
		}
	}

	c := cards.NewCard(q, a)
	err := cards.SaveCard(c)
	if err != nil {
		fmt.Printf("cannot save card: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("successfuly saved. id: %s\n", c.Id)
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

func addVocabCard() {
	reader := bufio.NewReader(os.Stdin)
	var word, def, example string

	for word == "" {
		fmt.Print("word: ")
		word, _ = reader.ReadString('\n')
		word = strings.TrimSpace(word)

		if word == "" {
			fmt.Println("enter the word")
		}
	}

	for def == "" {
		fmt.Print("definition: ")
		def, _ = reader.ReadString('\n')
		def = strings.TrimSpace(def)

		if def == "" {
			fmt.Println("enter the definition")
		}
	}

	for example == "" {
		fmt.Print("example: ")
		example, _ = reader.ReadString('\n')
		example = strings.TrimSpace(example)

		if example == "" {
			fmt.Println("enter the example")
		}
	}

	fmt.Printf("word: %s\n", word)
	fmt.Printf("def: %s\n", def)
	fmt.Printf("example: %s\n", example)

	var q, a string

	q = fmt.Sprintf("What's the meaning of \"%s\"?", word)
	a = def
	c := cards.NewCard(q, a)
	err := cards.SaveCard(c)
	if err != nil {
		fmt.Printf("cannot save card: %s\n", err)
		os.Exit(1)
	}

	q = fmt.Sprintf("Use \"%s\" in a sentence", word)
	a = example
	c = cards.NewCard(q, a)
	err = cards.SaveCard(c)
	if err != nil {
		fmt.Printf("cannot save card: %s\n", err)
		os.Exit(1)
	}

	q = fmt.Sprintf("Say \"%s\" out loud", word)
	a = word
	c = cards.NewCard(q, a)
	err = cards.SaveCard(c)
	if err != nil {
		fmt.Printf("cannot save card: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("successfully saved all cards")
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
