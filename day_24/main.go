package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
)

type Position struct {
	x int
	y int
}

type Object = string

var CODE_TO_OBJECT map[rune]Object = map[rune]Object{
	'.': "Free",
	'#': "Infested",
	'?': "Recurisve",
}
var OBJECT_TO_CODE map[Object]rune = map[Object]rune{
	"Free":      '.',
	"Infested":  '#',
	"Recursive": '?',
}

// ----------------------- BioSystem Struct Start -----------------------

type Tile = map[Position]Object

func new_BioSystem(lines []string, top_left Position) BioSystem {
	var mapping Tile = make(Tile)

	for line_index, line := range lines {
		for position_index, element := range line {

			// Add to mapping
			new_position := Position{top_left.x + position_index, top_left.y + line_index}
			object := CODE_TO_OBJECT[element]
			mapping[new_position] = object
		}
	}

	var bottom_right Position = Position{top_left.x + len(lines[0]) - 1, top_left.y + len(lines) - 1}
	return BioSystem{mapping, make([]map[Position]string, 0), top_left, bottom_right}
}

type BioSystem struct {
	actual  Tile
	history []Tile
	// Limits
	top_left     Position
	bottom_right Position
}

func (system *BioSystem) count_adjacent(state Object, position Position) int {
	var CHECK_POSITIONS []Position = []Position{
		Position{-1, 0},
		Position{1, 0},
		Position{0, -1},
		Position{0, 1},
	}

	var count int = 0
	for _, check_position := range CHECK_POSITIONS {

		// Creating position to check
		tmp_position_x := position.x + check_position.x
		tmp_position_y := position.y + check_position.y
		tmp_position := Position{tmp_position_x, tmp_position_y}

		// Check position
		object, position_is_set := system.actual[tmp_position]
		if position_is_set && object == state {
			count = count + 1
		}
	}

	return count
}

func (system *BioSystem) run_iteration() {
	// Add current mapping to history
	system.history = append(system.history, system.actual)

	// Create new mapping
	var new_mapping Tile = make(Tile)
	current_position := system.top_left
	for current_position.y <= system.bottom_right.x {

		// Implementation of rules
		var new_object Object
		current_object := system.actual[current_position]
		count_infested := system.count_adjacent("Infested", current_position)
		if current_object == "Infested" && count_infested != 1 {
			new_object = "Free"
		} else if current_object == "Free" && (count_infested == 1 || count_infested == 2) {
			new_object = "Infested"
		} else {
			new_object = current_object
		}

		// Store back object
		new_mapping[current_position] = new_object

		// Update current position
		if current_position.x == system.bottom_right.x {
			current_position.x = system.top_left.x
			current_position.y = current_position.y + 1

		} else {
			current_position.x = current_position.x + 1
		}
	}

	// Update current mapping
	system.actual = new_mapping
}

func (system *BioSystem) check_if_previous_state() bool {
	for _, other_state := range system.history {
		equals := reflect.DeepEqual(system.actual, other_state)
		if equals {
			return true
		}
	}

	return false
}

func (system *BioSystem) run_until_rep_state() {
	for !system.check_if_previous_state() {
		system.run_iteration()
	}
}

func (system *BioSystem) calcualte_biodiversity_rating() int {
	var count int = 0

	current_position := system.top_left
	for current_position.y <= system.bottom_right.x {

		element := system.actual[current_position]
		current_index := current_position.y*(system.bottom_right.x-system.top_left.x+1) + current_position.x
		if element == "Infested" {
			value := int(math.Pow(2.0, float64(current_index)))
			count = count + value
		}

		// Update current position
		if current_position.x == system.bottom_right.x {
			current_position.x = system.top_left.x
			current_position.y = current_position.y + 1

		} else {
			current_position.x = current_position.x + 1
		}
	}

	return count
}

func (system *BioSystem) print_current() {
	current_position := system.top_left
	for current_position.y <= system.bottom_right.x {

		element := system.actual[current_position]
		code := OBJECT_TO_CODE[element]
		fmt.Printf("%c ", code)

		// Update current position
		if current_position.x == system.bottom_right.x {
			current_position.x = system.top_left.x
			current_position.y = current_position.y + 1
			fmt.Println()

		} else {
			current_position.x = current_position.x + 1
		}
	}
}

// ----------------------- BioSystem Struct End -----------------------

// ----------------------- RecursiveBioSystem Struct Start -----------------------

type RecursiveTile map[int]map[Position]Object

func create_blank_tile(size int) map[Position]Object {
	mid_point := (size - 1) / 2
	// Creating a blank tiles
	blank_tile := make(map[Position]Object)
	for temp_y := 0; temp_y < size; temp_y++ {
		for temp_x := 0; temp_x < size; temp_x++ {
			tmp_position := Position{temp_x, temp_y}

			if temp_y == mid_point && temp_x == mid_point {
				blank_tile[tmp_position] = "Recursive"
				continue
			}

			blank_tile[tmp_position] = "Free"
		}
	}

	return blank_tile
}

func new_RecursiveBioSystem(lines []string) RecursiveBioSystem {
	var size int = len(lines)
	var mapping RecursiveTile = make(RecursiveTile)
	mid_point := (size - 1) / 2

	// Add level 0
	mapping[0] = create_blank_tile(size)
	for line_index, line := range lines {
		for position_index, element := range line {
			new_position := Position{position_index, line_index}

			if line_index == mid_point && position_index == mid_point {
				mapping[0][new_position] = "Recursive"
				continue
			}

			// Add to mapping
			object := CODE_TO_OBJECT[element]
			mapping[0][new_position] = object
		}
	}

	mapping[1] = create_blank_tile(size)
	mapping[-1] = create_blank_tile(size)

	return RecursiveBioSystem{0, mapping, -1, 1, size}
}

type RecursiveBioSystem struct {
	time   int
	actual RecursiveTile
	// Limits
	min_level int
	max_level int
	size      int
}

func (system *RecursiveBioSystem) count_adjacent(state Object, level int, position Position) int {
	var mid_point int = (system.size - 1) / 2
	var CHECK_POSITIONS []Position = []Position{
		Position{-1, 0},
		Position{1, 0},
		Position{0, -1},
		Position{0, 1},
	}

	var CHECK_POSITIONS_TO_INSIDE_LEVEL map[Position][2]Position = map[Position][2]Position{
		Position{-1, 0}: [2]Position{Position{system.size - 1, 0}, Position{0, 1}},
		Position{1, 0}:  [2]Position{Position{0, 0}, Position{0, 1}},
		Position{0, -1}: [2]Position{Position{0, system.size - 1}, Position{1, 0}},
		Position{0, 1}:  [2]Position{Position{0, 0}, Position{1, 0}},
	}

	var count int = 0
	for _, check_position := range CHECK_POSITIONS {

		// Creating position to check
		tmp_position_x := position.x + check_position.x
		tmp_position_y := position.y + check_position.y
		tmp_position := Position{tmp_position_x, tmp_position_y}

		// Check position
		object, position_is_set := system.actual[level][tmp_position]
		if position_is_set && object != "Recursive" && object == state {
			// Normal scenario
			count = count + 1

		} else if position_is_set && object == "Recursive" {
			// Go one level inside => level + 1
			sub_level_info := CHECK_POSITIONS_TO_INSIDE_LEVEL[check_position]
			sub_level_current_pos, sub_level_inc := sub_level_info[0], sub_level_info[1]

			for sub_level_current_pos.x < system.size && sub_level_current_pos.y < system.size {

				sub_object, sub_position_is_set := system.actual[level+1][sub_level_current_pos]
				if sub_position_is_set && sub_object == state {
					count = count + 1
				}

				// Update current_pos
				new_sub_level_current_pos_x := sub_level_current_pos.x + sub_level_inc.x
				new_sub_level_current_pos_y := sub_level_current_pos.y + sub_level_inc.y
				sub_level_current_pos = Position{new_sub_level_current_pos_x, new_sub_level_current_pos_y}
			}

		} else if !position_is_set {
			// Go one level outside => level - 1
			fetch_position_x := mid_point + check_position.x
			fetch_position_y := mid_point + check_position.y
			fetch_position := Position{fetch_position_x, fetch_position_y}

			sob_object, _ := system.actual[level-1][fetch_position]
			if sob_object == state {
				count = count + 1
			}
		}
	}

	return count
}

func (system *RecursiveBioSystem) run_iteration() {
	// Increment time
	system.time = system.time + 1

	// Create new mapping
	var new_mapping RecursiveTile = make(RecursiveTile)
	new_low_level, new_high_level := system.min_level, system.max_level

	for level := system.min_level; level <= system.max_level; level++ {

		current_position := Position{0, 0}
		new_mapping[level] = make(map[Position]string)
		at_least_one_turned_infected := false

		for current_position.y < system.size {

			// Implementation of rules
			var new_object Object
			current_object := system.actual[level][current_position]
			count_infested := system.count_adjacent("Infested", level, current_position)
			if current_object == "Infested" && count_infested != 1 {
				new_object = "Free"
			} else if current_object == "Free" && (count_infested == 1 || count_infested == 2) {
				at_least_one_turned_infected = true
				new_object = "Infested"
			} else {
				new_object = current_object
			}

			// Store back object
			new_mapping[level][current_position] = new_object

			// Update current position
			if current_position.x == system.size-1 {
				current_position.x = 0
				current_position.y = current_position.y + 1

			} else {
				current_position.x = current_position.x + 1
			}
		}

		// Update and create new levels
		if level == system.min_level && at_least_one_turned_infected {
			new_mapping[system.min_level-1] = create_blank_tile(system.size)
			new_low_level = system.min_level - 1
		}
		if level == system.max_level && at_least_one_turned_infected {
			new_mapping[system.max_level+1] = create_blank_tile(system.size)
			new_high_level = system.max_level + 1
		}
	}

	// Update current mapping
	system.actual = new_mapping
	system.min_level, system.max_level = new_low_level, new_high_level
}

func (system *RecursiveBioSystem) count_number_of(state string) int {
	var count int = 0

	for level := system.min_level; level <= system.max_level; level++ {

		current_position := Position{0, 0}
		for current_position.y < system.size {

			element := system.actual[level][current_position]
			if element == state {
				count = count + 1
			}

			// Update current position
			if current_position.x == system.size-1 {
				current_position.x = 0
				current_position.y = current_position.y + 1

			} else {
				current_position.x = current_position.x + 1
			}
		}
	}

	return count
}

func (system *RecursiveBioSystem) print_current() {

	for level := system.min_level; level <= system.max_level; level++ {

		fmt.Printf("Level %d:\n", level)
		current_position := Position{0, 0}
		for current_position.y < system.size {

			element := system.actual[level][current_position]
			code := OBJECT_TO_CODE[element]
			fmt.Printf("%c ", code)

			// Update current position
			if current_position.x == system.size-1 {
				current_position.x = 0
				current_position.y = current_position.y + 1
				fmt.Println()

			} else {
				current_position.x = current_position.x + 1
			}
		}

		fmt.Println("----------------")
	}
}

// ----------------------- RecursiveBioSystem Struct End -----------------------

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var lines []string = make([]string, 0)

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		lines = append(lines, line)
	}

	// Creating BioSystem (Part 1)
	var TOP_LEFT Position = Position{0, 0}

	var system BioSystem = new_BioSystem(lines, TOP_LEFT)
	system.run_until_rep_state()
	biodiversity_rating := system.calcualte_biodiversity_rating()
	fmt.Printf("The biodiversity rating is ' %d ' (part 1)\n", biodiversity_rating)

	// Creating RecursiveBioSystem (Part 2)
	var NUMBER_ITERATIONS int = 200
	var CELL_TYPE_TO_COUNT string = "Infested"

	var recursive_system RecursiveBioSystem = new_RecursiveBioSystem(lines)
	for index := 0; index < NUMBER_ITERATIONS; index++ {
		recursive_system.run_iteration()
	}
	number_infested := recursive_system.count_number_of("Infested")
	fmt.Printf("The number of ' %s ' cells after ' %d ' minutes is ' %d ' (part 2)\n", CELL_TYPE_TO_COUNT, NUMBER_ITERATIONS, number_infested)
}
