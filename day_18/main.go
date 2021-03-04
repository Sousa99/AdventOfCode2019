package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// ----------------------- Adventurer Struct Start -----------------------

const AdventurerSymbol = '@'
const WallSymbol = '#'
const FreeSymbol = '.'
const FistDoorSymbol = 'A'
const LastDoorSymbol = 'Z'
const FirstKeySymbol = 'a'
const LastKeySymbol = 'z'

type Position struct {
	x int
	y int
}

type Vector struct {
	from rune
	to   rune
}

type State struct {
	id          string
	current_key rune
	keys        []rune
	connections []string
}

type StateExplorer struct {
	current_state    int
	current_distance int
	current_path     []rune
}

func make_copy_state_explorer(original StateExplorer) StateExplorer {
	copy_current_state := original.current_state
	copy_distance := original.current_distance
	copy_current_path := make([]rune, 0, len(original.current_path))
	copy_current_path = append(copy_current_path, original.current_path...)

	var copy_state_explorer StateExplorer = StateExplorer{copy_current_state, copy_distance, copy_current_path}
	return copy_state_explorer
}

type LittleExplorer struct {
	position      Position
	last_position Position
	distance      int
	restrictions  []rune
}

func make_copy_little_explorer(original LittleExplorer) LittleExplorer {
	copy_position := Position{original.position.x, original.position.y}
	copy_last_position := Position{original.last_position.x, original.last_position.y}
	copy_distance := original.distance
	copy_restrictions := make([]rune, 0, len(original.restrictions))
	copy_restrictions = append(copy_restrictions, original.restrictions...)

	var copy_little LittleExplorer = LittleExplorer{copy_position, copy_last_position, copy_distance, copy_restrictions}
	return copy_little
}

var Directions []Position = []Position{
	Position{1, 0},
	Position{-1, 0},
	Position{0, 1},
	Position{0, -1},
}

type Adventurer struct {
	// Current state
	position     Position
	distances    map[Vector]int
	restrictions map[rune][]rune
	states       map[string]State
	// Mapping
	top_left      Position
	bottom_right  Position
	mapping       map[Position]rune
	keys_position map[rune]Position
}

func (adventurer *Adventurer) add_mapping_position(position Position, mapping_code rune) {
	// Add mapping
	if mapping_code == AdventurerSymbol {
		adventurer.keys_position[AdventurerSymbol] = position
		adventurer.position = position
		mapping_code = FreeSymbol
	}
	adventurer.mapping[position] = mapping_code

	if mapping_code >= FirstKeySymbol && mapping_code <= LastKeySymbol {
		adventurer.keys_position[mapping_code] = position
	}

	// Update Limits
	if position.x < adventurer.top_left.x {
		adventurer.top_left.x = position.x
	} else if position.x > adventurer.bottom_right.x {
		adventurer.bottom_right.x = position.x
	}
	if position.y < adventurer.top_left.y {
		adventurer.top_left.y = position.y
	} else if position.y > adventurer.bottom_right.y {
		adventurer.bottom_right.y = position.y
	}
}

func (adventurer *Adventurer) compute_distances() {
	for key, key_position := range adventurer.keys_position {

		adventurer.distances[Vector{key, key}] = 0
		visited_positions := make(map[Position]bool)

		var current_explorers []LittleExplorer = []LittleExplorer{
			LittleExplorer{key_position, key_position, 0, make([]rune, 0)},
		}

		for len(current_explorers) > 0 {

			new_explorers := make([]LittleExplorer, 0)
			for _, current_explorer := range current_explorers {
				for _, direction := range Directions {
					tmp_position_x := current_explorer.position.x + direction.x
					tmp_position_y := current_explorer.position.y + direction.y
					var new_position Position = Position{tmp_position_x, tmp_position_y}

					visited, visited_is_set := visited_positions[new_position]
					if current_explorer.last_position == new_position || (visited_is_set && visited) {
						// Little explorer should not go back
						continue
					}

					map_code, is_set := adventurer.mapping[new_position]
					if !is_set || map_code == WallSymbol {
						// Can't go through walls
						continue
					}

					visited_positions[new_position] = true
					// Update and create new little explorer
					new_explorer := make_copy_little_explorer(current_explorer)
					new_explorer.last_position = new_explorer.position
					new_explorer.position = new_position
					new_explorer.distance = new_explorer.distance + 1
					new_explorers = append(new_explorers, new_explorer)

					if map_code >= FirstKeySymbol && map_code <= LastKeySymbol {
						// Found a key
						adventurer.distances[Vector{key, map_code}] = new_explorer.distance
					}
				}
			}

			current_explorers = new_explorers
		}
	}
}

func (adventurer *Adventurer) compute_restrictions() {
	visited_positions := make(map[Position]bool)

	var current_explorers []LittleExplorer = []LittleExplorer{
		LittleExplorer{adventurer.position, adventurer.position, 0, make([]rune, 0)},
	}

	for len(current_explorers) > 0 {

		new_explorers := make([]LittleExplorer, 0)
		for _, current_explorer := range current_explorers {
			for _, direction := range Directions {
				tmp_position_x := current_explorer.position.x + direction.x
				tmp_position_y := current_explorer.position.y + direction.y
				var new_position Position = Position{tmp_position_x, tmp_position_y}

				visited, visited_is_set := visited_positions[new_position]
				if current_explorer.last_position == new_position || (visited_is_set && visited) {
					// Little explorer should not go back
					continue
				}

				map_code, is_set := adventurer.mapping[new_position]
				if !is_set || map_code == WallSymbol {
					// Can't go through walls
					continue
				}

				visited_positions[new_position] = true
				// Update and create new little explorer
				new_explorer := make_copy_little_explorer(current_explorer)
				new_explorer.last_position = new_explorer.position
				new_explorer.position = new_position

				if map_code >= FirstKeySymbol && map_code <= LastKeySymbol {
					// Found a key
					adventurer.restrictions[map_code] = make([]rune, 0, len(new_explorer.restrictions))
					adventurer.restrictions[map_code] = append(adventurer.restrictions[map_code], new_explorer.restrictions...)
				} else if map_code >= FistDoorSymbol && map_code <= LastDoorSymbol {
					// Found door
					new_explorer.restrictions = append(new_explorer.restrictions, map_code)
				}

				new_explorers = append(new_explorers, new_explorer)
			}
		}

		current_explorers = new_explorers
	}
}

func (adventurer *Adventurer) compute_states() {
	var possibilities map[string][]rune = make(map[string][]rune)

	keys := make([]rune, 0)
	starting_state_code := convert_state_keys_to_code(keys)
	var starting_state State = State{starting_state_code, AdventurerSymbol, keys, make([]string, 0)}
	adventurer.states[starting_state_code] = starting_state

	var processing []string = []string{starting_state_code}

	for len(processing) > 0 {
		fmt.Printf("\033[2K\rState Processment: %d", len(processing))

		// Retrieve state being changed
		current_state := adventurer.states[processing[0]]
		state_keys_code := convert_state_keys_to_code(current_state.keys)

		going_to_possibilities, is_set := possibilities[state_keys_code]
		if !is_set {
			// Find possibilities
			going_to_possibilities = make([]rune, 0)
			for going_to, restrictions := range adventurer.restrictions {
				if going_to == current_state.current_key || slice_contains(current_state.keys, going_to) {
					// No interest in going to itself or somewhere it has been
					continue
				}

				has_keys_needed := true
				for _, door_encountered := range restrictions {
					key_needed := convert_door_to_key(door_encountered)
					if !slice_contains(current_state.keys, key_needed) {
						has_keys_needed = false
					}
				}

				if !has_keys_needed {
					// It doesn't have all the keys needed
					continue
				}

				going_to_possibilities = append(going_to_possibilities, going_to)
				possibilities[state_keys_code] = going_to_possibilities
			}
		}

		for _, going_to := range going_to_possibilities {
			new_keys := make([]rune, 0, len(current_state.keys))
			new_keys = append(new_keys, current_state.keys...)
			if !slice_contains(new_keys, going_to) {
				new_keys = append(new_keys, going_to)
			}

			to_state_code := convert_state_to_code(going_to, new_keys)
			_, state_already_added := adventurer.states[to_state_code]
			if !state_already_added {
				new_state := State{to_state_code, going_to, new_keys, make([]string, 0)}
				adventurer.states[to_state_code] = new_state
				processing = append(processing, to_state_code)
			}

			current_state.connections = append(current_state.connections, to_state_code)
		}

		// Save back
		adventurer.states[processing[0]] = current_state
		if len(processing) > 1 {
			processing = processing[1:]
		} else {
			processing = make([]string, 0)
		}
	}

	fmt.Println()
}

func (adventurer *Adventurer) start_adventure() {
	var final_state string = ""

	var vertexes_dist map[string]int = make(map[string]int)
	var vertexes_prev map[string]string = make(map[string]string)
	var vertexes []string = make([]string, 0, len(adventurer.states))
	for key, _ := range adventurer.states {
		vertexes_dist[key] = math.MaxInt32
		vertexes_prev[key] = ""
		vertexes = append(vertexes, key)
	}

	keys := make([]rune, 0)
	starting_state_code := convert_state_keys_to_code(keys)
	vertexes_dist[starting_state_code] = 0

	for len(vertexes) > 0 {
		fmt.Printf("\033[2K\rVertexes Q: %d", len(vertexes))
		sort.SliceStable(vertexes, func(i, j int) bool {
			return vertexes_dist[vertexes[i]] < vertexes_dist[vertexes[j]]
		})

		processing_id := vertexes[0]
		vertexes = vertexes[1:]
		current_state := adventurer.states[processing_id]

		if len(current_state.connections) == 0 {
			// Final state
			final_state = processing_id
			break
		}

		for _, connection := range current_state.connections {
			real_connection := adventurer.states[connection]
			alt := vertexes_dist[processing_id] + adventurer.distances[Vector{current_state.current_key, real_connection.current_key}]

			// Djikstra
			if alt < vertexes_dist[connection] {
				vertexes_dist[connection] = alt
				vertexes_prev[connection] = processing_id
			}
		}
	}
	fmt.Println()

	fmt.Printf("Distance ' %d '\n", vertexes_dist[final_state])
}

func (adventurer *Adventurer) print_mapping(print_adventurer bool) {
	var current_position Position = adventurer.top_left

	// Iterate until last line
	for current_position.y <= adventurer.bottom_right.y {
		code := adventurer.mapping[current_position]
		if print_adventurer && adventurer.position == current_position {
			code = AdventurerSymbol
		}

		fmt.Printf("%c ", code)

		if current_position.x == adventurer.bottom_right.x {
			// Last pixel on line
			fmt.Println()
			current_position.x = adventurer.top_left.x
			current_position.y = current_position.y + 1
		} else {
			// Not the last pixel
			current_position.x = current_position.x + 1
		}
	}
}

// ----------------------- Adventurer Struct End -----------------------

func convert_slice_rune_to_string(runes []rune) []string {
	result := make([]string, 0)
	for _, code := range runes {
		result = append(result, string(code))
	}

	return result
}

func convert_door_to_key(door rune) rune {
	const differ = FistDoorSymbol - FirstKeySymbol
	var key rune = door - differ

	return key
}

func slice_contains(slice []rune, elem rune) bool {
	for _, slice_elem := range slice {
		if slice_elem == elem {
			return true
		}
	}

	return false
}

func convert_state_keys_to_code(keys []rune) string {
	var final_code []rune = make([]rune, 0)
	for code := FirstKeySymbol; code <= LastKeySymbol; code++ {
		final_code = append(final_code, '0')
	}

	for _, key := range keys {
		key_translated := key - FirstKeySymbol
		final_code[int(key_translated)] = '1'
	}

	return string(final_code)
}

func convert_state_to_code(current_key rune, keys []rune) string {
	var final_code []rune = make([]rune, 0)
	for code := FirstKeySymbol; code <= LastKeySymbol; code++ {
		final_code = append(final_code, '0')
	}

	for _, key := range keys {
		key_translated := key - FirstKeySymbol
		final_code[int(key_translated)] = '1'
	}

	key_translated := current_key - FirstKeySymbol
	final_code[int(key_translated)] = '2'

	return string(final_code)
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	position_0 := Position{0, 0}
	var adventurer Adventurer = Adventurer{position_0, make(map[Vector]int), make(map[rune][]rune), make(map[string]State), position_0, position_0, make(map[Position]rune), make(map[rune]Position)}

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	line_index := 0
	for scanner.Scan() {

		var line string = scanner.Text()
		for column_index, characther := range line {
			new_position := Position{column_index, line_index}
			adventurer.add_mapping_position(new_position, characther)
		}

		line_index = line_index + 1
	}

	// adventurer.print_mapping(true)
	adventurer.compute_distances()
	adventurer.compute_restrictions()
	adventurer.compute_states()
	adventurer.start_adventure()
	/*
		adventurer.start()
		fmt.Printf("Keys collected: %v\n", convert_slice_rune_to_string(adventurer.keys))
		fmt.Printf("Number of steps: ' %d '(part 1)\n", len(adventurer.history))~
	*/
}
