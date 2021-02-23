package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ----------------------- Reaction Struct Start -----------------------

type ReactionPart struct {
	quantity int
	product  string
}

type Reaction struct {
	products []ReactionPart
	result   ReactionPart
}

type System struct {
	low_level  string
	high_level string
	products   []string
	reactions  []Reaction
	// Final result
	ore_per_fuel int
	extras       map[string]int
}

func (system *System) add_reaction(reaction Reaction) {
	system.products = append(system.products, reaction.result.product)
	system.reactions = append(system.reactions, reaction)
}

func (system *System) find_quantity_per_high_level() int {
	var extras map[string]int = make(map[string]int)
	var quantities map[string]int = make(map[string]int)

	// Initialize extras and quantities
	for _, product := range system.products {
		extras[product] = 0
		quantities[product] = 0
	}

	// Initialize with recipe for fuel
	for _, reaction := range system.reactions {
		if reaction.result.product == system.high_level {
			for _, product := range reaction.products {
				quantities[product.product] = product.quantity
			}
		}
	}

	// Convert down materials
	var dumb_down bool = false
	for !dumb_down {

		for product_needed, quantity_needed := range quantities {
			// Skip if already low level or high_level
			if product_needed == system.high_level || product_needed == system.low_level || quantity_needed == 0 {
				continue
			}

			for _, reaction := range system.reactions {
				if reaction.result.product != product_needed {
					continue
				}

				number_times_reaction := int(math.Ceil(float64(quantity_needed) / float64(reaction.result.quantity)))
				extras[product_needed] = number_times_reaction*reaction.result.quantity - quantity_needed

				for _, sub_product := range reaction.products {

					sub_product_quantity := number_times_reaction * sub_product.quantity
					sub_total := sub_product_quantity - extras[sub_product.product]
					extras[sub_product.product] = 0

					if sub_total > 0 {
						quantities[sub_product.product] += sub_total
					} else {
						extras[sub_product.product] -= sub_total
					}
				}

				quantities[product_needed] = 0

				break
			}
		}

		// Check if dumb down enough
		dumb_down = true
		for product, quantity := range quantities {
			if product != system.low_level && quantity > 0 {
				dumb_down = false
			}
		}

	}

	// Return quantity
	final_quantity, is_set := quantities[system.low_level]
	system.extras = extras
	if is_set {
		system.ore_per_fuel = final_quantity
		return final_quantity
	} else {
		system.ore_per_fuel = 0
		return 0
	}
}

func (system *System) find_quantity_given_low_level(low_level_quantity int) int {
	var quantity_possible int = low_level_quantity / system.ore_per_fuel
	var quantity_added int = quantity_possible

	var current_extras map[string]int = make(map[string]int)
	for product, _ := range system.extras {
		current_extras[product] = 0
	}
	current_extras[system.low_level] = low_level_quantity % system.ore_per_fuel

	// Dumb down materials to ore
	var completed bool = false
	for !completed {
		// If nothing was transformed then by omission is completed
		completed = true

		// Initialize current_extras
		for product, quantity := range system.extras {
			current_extras[product] += quantity_added * quantity
		}

		// Iterate reaction see if any possible
		for _, reaction := range system.reactions {
			// If there are not enough extras
			extras_of_product := current_extras[reaction.result.product]
			if extras_of_product < reaction.result.quantity {
				continue
			}

			completed = false
			number_of_reactions_possible := extras_of_product / reaction.result.quantity
			current_extras[reaction.result.product] = extras_of_product % reaction.result.quantity

			for _, reagent := range reaction.products {
				current_extras[reagent.product] += number_of_reactions_possible * reagent.quantity
			}
		}

		// Verify if ore can be turned to fuel
		quantity_added = current_extras[system.low_level] / system.ore_per_fuel
		quantity_possible = quantity_possible + quantity_added
		current_extras[system.low_level] = current_extras[system.low_level] % system.ore_per_fuel
	}

	return quantity_possible + 1
}

// ----------------------- Reaction Struct Start -----------------------

func get_reaction_from_string(product_string string) ReactionPart {
	split := strings.Split(product_string, " ")

	quantity, _ := strconv.Atoi(split[0])
	product := split[1]

	return ReactionPart{quantity, product}
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var products []string = make([]string, 0)
	products = append(products, "ORE")
	var system System = System{"ORE", "FUEL", products, make([]Reaction, 0), -1, make(map[string]int)}

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		var split []string = strings.Split(line, " => ")

		products_line := split[0]
		result_line := split[1]

		var products []ReactionPart = make([]ReactionPart, 0)
		products_split := strings.Split(products_line, ", ")
		for _, product_string := range products_split {
			reaction_part := get_reaction_from_string(product_string)
			products = append(products, reaction_part)
		}

		result_reaction_part := get_reaction_from_string(result_line)
		var new_reaction Reaction = Reaction{products, result_reaction_part}
		system.add_reaction(new_reaction)
	}

	// Part 1
	quantity := system.find_quantity_per_high_level()
	fmt.Printf("It is needed ' %d ' of ' %s ' to produce one reaction of ' %s ' (part 1)\n", quantity, system.low_level, system.high_level)

	// Part 2
	ore_extrapolated := 1000000000000
	quantity = system.find_quantity_given_low_level(ore_extrapolated)
	fmt.Printf("We can have ' %d ' of ' %s ' if given ' %d ' of ' %s ' (part 2)\n", quantity, system.high_level, ore_extrapolated, system.low_level)
}
