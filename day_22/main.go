package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
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

type BigInt = big.Int

type LinFunc struct {
	multiply BigInt
	sum      BigInt
}

func (function *LinFunc) run(value BigInt) BigInt {
	value.Mul(&function.multiply, &value)
	value.Add(&function.sum, &value)
	return value
}

func agg_functions(f LinFunc, g LinFunc, mod BigInt) LinFunc {
	// Should be read as: f(g(x))
	// Let f(x)=k*x+m and g(x)=j*x+n, then h(x) = f(g(x)) = Ax+B = k*(j*x+n)+m = k*j*x + k*n + m => A=k*j, B=k*n+m
	var new_multiply, new_add BigInt
	new_multiply.Mul(&f.multiply, &g.multiply)
	new_multiply.Mod(&new_multiply, &mod)

	new_add.Mul(&f.multiply, &g.sum)
	new_add.Add(&f.sum, &new_add)
	new_add.Mod(&new_add, &mod)

	return LinFunc{new_multiply, new_add}
}

type ReverseDeck struct {
	size             BigInt
	current_function LinFunc
}

func (deck *ReverseDeck) deal_into_new_stack() {
	// f (x) = - x + ( nCards - 1 )
	var new_multiply, new_add BigInt
	new_multiply = *big.NewInt(-1)
	new_add.Sub(&deck.size, big.NewInt(1))

	var new_function LinFunc = LinFunc{new_multiply, new_add}
	deck.current_function = agg_functions(new_function, deck.current_function, deck.size)
}

func (deck *ReverseDeck) cut_n(n BigInt) {
	// f (x) = x + n mod nCards
	var new_multiply, new_add BigInt
	new_multiply = *big.NewInt(1)
	new_add.Mod(&n, &deck.size)

	var new_function LinFunc = LinFunc{new_multiply, new_add}
	deck.current_function = agg_functions(new_function, deck.current_function, deck.size)
}

func (deck *ReverseDeck) deal_with_increment_n(n BigInt) {
	// Being mod_inverse the mod inverse n of nCards
	// f (x) = (z mod (nCards)) x + 0
	var new_multiply BigInt
	new_multiply.ModInverse(&n, &deck.size)
	new_multiply.Mod(&new_multiply, &deck.size)
	var new_function LinFunc = LinFunc{new_multiply, *big.NewInt(0)}
	deck.current_function = agg_functions(new_function, deck.current_function, deck.size)
}

func (deck *ReverseDeck) find_card(position BigInt) BigInt {
	var value BigInt = deck.current_function.run(position)
	value.Mod(&value, &deck.size)
	if value.Cmp(big.NewInt(0)) == -1 {
		value.Add(&value, &deck.size)
	}

	return value
}

func (deck *ReverseDeck) do_reps(reps int) {
	var base_two_to int = int(math.Floor(math.Log2(float64(reps))))

	var functions map[int]LinFunc = map[int]LinFunc{0: deck.current_function}
	for base := 1; base <= base_two_to; base++ {
		previous_function := functions[base-1]
		new_function := agg_functions(previous_function, previous_function, deck.size)
		functions[base] = new_function
	}

	var final_function LinFunc = LinFunc{*big.NewInt(1), *big.NewInt(0)}
	value := reps
	for value != 0 {

		base_two_to_do := int(math.Floor(math.Log2(float64(value))))
		final_function = agg_functions(final_function, functions[base_two_to_do], deck.size)
		value = value - int(math.Pow(2.0, float64(base_two_to_do)))
	}

	deck.current_function = final_function
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
		deck.cut_n(*big.NewInt(int64(n)))

	} else if strings.HasPrefix(line, DEAL_WITH_INCREMENT) {
		// Cut n
		line = strings.TrimPrefix(line, DEAL_WITH_INCREMENT)
		n, _ := strconv.Atoi(line)
		deck.deal_with_increment_n(*big.NewInt(int64(n)))

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
	DECK_SIZE_TO = 119315717514047
	var POSITION_TO_FIND BigInt = *big.NewInt(int64(2020))
	var REPETITIONS int = 101741582076661

	var original_function LinFunc = LinFunc{*big.NewInt(1), *big.NewInt(0)}
	var reverse_deck ReverseDeck = ReverseDeck{*big.NewInt(int64(DECK_SIZE_TO)), original_function}
	for index_line := len(code_lines) - 1; index_line >= 0; index_line-- {
		run_on_reverse_deck(&reverse_deck, code_lines[index_line])
	}

	reverse_deck.do_reps(REPETITIONS)
	var card_found BigInt = reverse_deck.find_card(POSITION_TO_FIND)
	fmt.Printf("Card ' %s ' found at ' %s ' (part 2)\n", card_found.Text(10), POSITION_TO_FIND.Text(10))
}
