package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ----------------------- Deck Struct Start -----------------------

func new_deck(to int) Deck {
	var cards []Card = make([]Card, 0)
	for card_number := 0; card_number < to; card_number++ {
		cards = append(cards, card_number)
	}

	return Deck{cards}
}

type Card = int
type Deck struct {
	cards []Card
}

func (deck *Deck) deal_into_new_stack() {
	for i, j := 0, len(deck.cards)-1; i < j; i, j = i+1, j-1 {
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	}
}

func (deck *Deck) cut_n(n int) {

	// Deal with negative values
	if n < 0 {
		n = n + len(deck.cards)
	}

	var cut []Card = deck.cards[:n]
	deck.cards = deck.cards[n:]
	deck.cards = append(deck.cards, cut...)
}

func (deck *Deck) deal_with_increment_n(n int) {
	var new_deck []Card = make([]int, len(deck.cards))

	var current_index int = 0
	for _, card := range deck.cards {

		new_deck[current_index] = card
		// Update current_index
		current_index = current_index + n
		current_index = current_index % len(deck.cards)
	}

	deck.cards = new_deck
}

func (deck *Deck) find_index_of(number int) int {
	for index, card := range deck.cards {
		if card == number {
			return index
		}
	}

	return -1
}

// ----------------------- Deck Struct End -----------------------~

// ----------------------- Focused Deck Struct Start -----------------------

type ReverseDeck struct {
	size     int
	position int
}

func (deck *ReverseDeck) deal_into_new_stack() {
	deck.position = deck.size - deck.position - 1
}

func (deck *ReverseDeck) cut_n(n int) {

	// Deal with negative values
	if n < 0 {
		n = n + deck.size
	}

	// Update element
	if deck.position >= deck.size-n {
		// On cut
		deck.position = deck.position - (deck.size - n)
	} else {
		// Not on cut
		deck.position = deck.position + n
	}
}

func (deck *ReverseDeck) deal_with_increment_n(n int) {

	// Some problem here
	if deck.position == 0 {
		return
	}

	// Find euclidean inverse of n
	moduled_n := n % deck.size
	inverse := 1

	for true {
		if (moduled_n*inverse)%deck.size == 1 {
			break
		}

		inverse = inverse + 1
	}

	deck.position = ((n * inverse) % (deck.position * n)) / n

}

func (deck *ReverseDeck) find_index_stored() int {
	return deck.position
}

// ----------------------- Focused Deck Struct End -----------------------

func run_on_deck(deck *Deck, line string) {
	var DEAL_INTO_NEW_STACK string = "deal into new stack"
	var CUT_N string = "cut "
	var DEAL_WITH_INCREMENT string = "deal with increment "

	if strings.HasPrefix(line, DEAL_INTO_NEW_STACK) {
		// Deal into new stack
		deck.deal_into_new_stack()

	} else if strings.HasPrefix(line, CUT_N) {
		// Cut n
		line = strings.TrimPrefix(line, CUT_N)
		n, _ := strconv.Atoi(line)
		deck.cut_n(n)

	} else if strings.HasPrefix(line, DEAL_WITH_INCREMENT) {
		// Cut n
		line = strings.TrimPrefix(line, DEAL_WITH_INCREMENT)
		n, _ := strconv.Atoi(line)
		deck.deal_with_increment_n(n)

	} else {
		fmt.Printf("Code ' %s ' not recognized!\n", line)
		os.Exit(1)
	}
}

func run_on_reverse_deck(deck *ReverseDeck, line string) {
	var DEAL_INTO_NEW_STACK string = "deal into new stack"
	var CUT_N string = "cut "
	var DEAL_WITH_INCREMENT string = "deal with increment "

	if strings.HasPrefix(line, DEAL_INTO_NEW_STACK) {
		// Deal into new stack
		deck.deal_into_new_stack()

	} else if strings.HasPrefix(line, CUT_N) {
		// Cut n
		line = strings.TrimPrefix(line, CUT_N)
		n, _ := strconv.Atoi(line)
		deck.cut_n(n)

	} else if strings.HasPrefix(line, DEAL_WITH_INCREMENT) {
		// Cut n
		line = strings.TrimPrefix(line, DEAL_WITH_INCREMENT)
		n, _ := strconv.Atoi(line)
		deck.deal_with_increment_n(n)

	} else {
		fmt.Printf("Code ' %s ' not recognized!\n", line)
		os.Exit(1)
	}
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var code_lines []string = make([]string, 0)
	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		code_lines = append(code_lines, line)
	}

	// Part 1
	var DECK_SIZE_TO int = 10007
	var CARD_TO_FIND int = 2019

	var deck Deck = new_deck(DECK_SIZE_TO)
	for _, line := range code_lines {
		run_on_deck(&deck, line)
	}
	index := deck.find_index_of(CARD_TO_FIND)
	fmt.Printf("Card ' %d ' found at ' %d ' (part 1)\n", CARD_TO_FIND, index)

	// Part 2
	DECK_SIZE_TO = 10007
	var POSITION_TO_FIND int = 2020
	var REPETITIONS int = 1

	var reverse_deck ReverseDeck = ReverseDeck{DECK_SIZE_TO, POSITION_TO_FIND}
	for rep := 0; rep < REPETITIONS; rep++ {
		for index_line := len(code_lines) - 1; index_line >= 0; index_line-- {
			run_on_reverse_deck(&reverse_deck, code_lines[index_line])
		}
	}

	var card_found int = reverse_deck.find_index_stored()
	fmt.Printf("Card ' %d ' found at ' %d ' (part 2)\n", card_found, POSITION_TO_FIND)
}
