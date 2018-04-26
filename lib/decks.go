package decks

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

// Hand organizes the modules' output
type Hand struct {
	Theme  string
	Medium string
	Noun   string
}

// TestDecks is just to make sure my modules work
func TestDecks() string {
	return "1 2 3"
}

// drawFromDeck draws a card (a line) from the specified deck
func drawFromDeck(dName string) (string, error) {
	dBytes, err := ioutil.ReadFile(dName)
	if err != nil {
		return "z", err
	}
	dString := string(dBytes)

	dString = strings.Replace(dString, "\r", "", -1)
	dLines := strings.Split(dString, "\n")
	rInt := rand.Intn(len(dLines))
	return dLines[rInt], nil
}

// DrawHand returns cards from all 4 decks
func DrawHand() Hand {
	rand.Seed(time.Now().UnixNano())
	theme, err := drawFromDeck("decks/themes.txt")
	medium, err := drawFromDeck("decks/media.txt")
	noun, err := drawFromDeck("decks/nouns.txt")
	if err != nil {
		log.Fatal(err)
	}
	hand := Hand{
		Theme:  theme,
		Medium: medium,
		Noun:   noun,
	}
	return hand
}
