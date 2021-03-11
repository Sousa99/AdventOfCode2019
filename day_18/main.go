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

var Directions []Position = []Position{
	Position{1, 0},
	Position{-1, 0},
	Position{0, 1},
	Position{0, -1},
}

type Adventurer struct {
	// Current state
	position Position
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

func (adventurer *Adventurer) compute_distances() map[Vector]int {

	/* =============== AUXILIARY STRUCT =============== */
	type DistanceExplorer struct {
		position      Position
		last_position Position
		distance      int
	}

	make_copy_distance_explorer := func(original DistanceExplorer) DistanceExplorer {
		copy_position := Position{original.position.x, original.position.y}
		copy_last_position := Position{original.last_position.x, original.last_position.y}
		copy_distance := original.distance

		var copy_little DistanceExplorer = DistanceExplorer{copy_position, copy_last_position, copy_distance}
		return copy_little
	}

	/* =============== END AUXILIARY STRUCT =============== */

	var distances map[Vector]int = make(map[Vector]int)
	for key, key_position := range adventurer.keys_position {

		distances[Vector{key, key}] = 0
		visited_positions := make(map[Position]bool)

		var current_explorers []DistanceExplorer = []DistanceExplorer{
			DistanceExplorer{key_position, key_position, 0},
		}

		for len(current_explorers) > 0 {

			new_explorers := make([]DistanceExplorer, 0)
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
					new_explorer := make_copy_distance_explorer(current_explorer)
					new_explorer.last_position = new_explorer.position
					new_explorer.position = new_position
					new_explorer.distance = new_explorer.distance + 1
					new_explorers = append(new_explorers, new_explorer)

					if map_code >= FirstKeySymbol && map_code <= LastKeySymbol {
						// Found a key
						distances[Vector{key, map_code}] = new_explorer.distance
					}
				}
			}

			current_explorers = new_explorers
		}
	}

	return distances
}

func (adventurer *Adventurer) compute_restrictions() map[rune][]rune {

	/* =============== AUXILIARY STRUCT =============== */
	type RestrictionsExplorer struct {
		position      Position
		last_position Position
		restrictions  []rune
	}

	make_copy_restrictions_explorer := func(original RestrictionsExplorer) RestrictionsExplorer {
		copy_position := Position{original.position.x, original.position.y}
		copy_last_position := Position{original.last_position.x, original.last_position.y}
		copy_restrictions := make([]rune, 0, len(original.restrictions))
		copy_restrictions = append(copy_restrictions, original.restrictions...)

		var copy_little RestrictionsExplorer = RestrictionsExplorer{copy_position, copy_last_position, copy_restrictions}
		return copy_little
	}

	/* =============== END AUXILIARY STRUCT =============== */

	var restrictions map[rune][]rune = make(map[rune][]rune)
	visited_positions := make(map[Position]bool)
	var current_explorers []RestrictionsExplorer = []RestrictionsExplorer{
		RestrictionsExplorer{adventurer.position, adventurer.position, make([]rune, 0)},
	}

	for len(current_explorers) > 0 {

		new_explorers := make([]RestrictionsExplorer, 0)
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
				new_explorer := make_copy_restrictions_explorer(current_explorer)
				new_explorer.last_position = new_explorer.position
				new_explorer.position = new_position

				if map_code >= FirstKeySymbol && map_code <= LastKeySymbol {
					// Found a key
					restrictions[map_code] = make([]rune, 0, len(new_explorer.restrictions))
					restrictions[map_code] = append(restrictions[map_code], new_explorer.restrictions...)
				} else if map_code >= FistDoorSymbol && map_code <= LastDoorSymbol {
					// Found door
					new_explorer.restrictions = append(new_explorer.restrictions, map_code)
				}

				new_explorers = append(new_explorers, new_explorer)
			}
		}

		current_explorers = new_explorers
	}

	return restrictions
}

func (adventurer *Adventurer) compute_states(distances map[Vector]int, restrictions map[rune][]rune) MapGraph {

	/* =============== AUXILIARY STRUCT =============== */
	type StateExplorer struct {
		current_key rune
		keys        []rune
	}

	make_copy_state_explorer := func(original StateExplorer) StateExplorer {
		copy_current_key := original.current_key
		copy_keys := make([]rune, 0, len(original.keys))
		copy_keys = append(copy_keys, original.keys...)

		var copy_little StateExplorer = StateExplorer{copy_current_key, copy_keys}
		return copy_little
	}

	/* =============== END AUXILIARY STRUCT =============== */

	var map_graph MapGraph = make(MapGraph)

	starting_explorer := StateExplorer{'@', make([]rune, 0)}
	var current_explorers []StateExplorer = []StateExplorer{starting_explorer}

	for explorer_index := 0; explorer_index < len(current_explorers); explorer_index++ {

		current_explorer := current_explorers[explorer_index]
		state_code := convert_state_to_code(current_explorer.current_key, current_explorer.keys)
		_, map_set_for_state := map_graph[state_code]
		if map_set_for_state {
			continue
		}

		map_graph[state_code] = make([]Connection, 0)
		for going_to, restrictions := range restrictions {
			if going_to == current_explorer.current_key || slice_rune_contains(current_explorer.keys, going_to) {
				// No interest in going to itself or somewhere it has been
				continue
			}

			has_keys_needed := true
			for _, door_encountered := range restrictions {
				key_needed := convert_door_to_key(door_encountered)
				if !slice_rune_contains(current_explorer.keys, key_needed) && key_needed != current_explorer.current_key {
					has_keys_needed = false
				}
			}

			if !has_keys_needed {
				// It doesn't have all the keys needed
				continue
			}

			// Create new state
			new_state_explorer := make_copy_state_explorer(current_explorer)
			new_state_explorer.keys = append(new_state_explorer.keys, new_state_explorer.current_key)
			new_state_explorer.current_key = going_to
			current_explorers = append(current_explorers, new_state_explorer)

			// Add connection
			distance := distances[Vector{current_explorer.current_key, going_to}]
			new_state_code := convert_state_to_code(new_state_explorer.current_key, new_state_explorer.keys)
			new_connection := Connection{new_state_code, distance}
			map_graph[state_code] = append(map_graph[state_code], new_connection)
		}
	}

	return map_graph
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

// ----------------------- Djikstra Algorithm Start -----------------------

type StateCode = string

type Connection struct {
	state_to StateCode
	distance int
}

type MapGraph = map[StateCode][]Connection

func run_dijkstra(connections MapGraph) int {

	type Node struct {
		distance int
		previous string
	}

	// Initialize nodes in correct form
	var nodes map[string]Node = make(map[string]Node)
	for node_name, _ := range connections {
		// Add node
		new_node := Node{math.MaxInt32, ""}
		nodes[node_name] = new_node
	}

	// Initialize source with dist 0
	source_code := convert_state_to_code('@', make([]rune, 0))
	source_node, _ := nodes[source_code]
	source_node.distance = 0
	nodes[source_code] = source_node
	var active_vertexes []string = []string{source_code}

	// Run algorithm
	for len(active_vertexes) > 0 {

		// Sort and remove first element
		sort.Slice(active_vertexes, func(i, j int) bool { return nodes[active_vertexes[i]].distance < nodes[active_vertexes[j]].distance })
		node_name := active_vertexes[0]
		active_vertexes = active_vertexes[1:]

		if len(connections[node_name]) == 0 {
			// We are only interest in target
			return nodes[node_name].distance
		}

		for _, connection := range connections[node_name] {

			alt := nodes[node_name].distance + connection.distance
			if alt < nodes[connection.state_to].distance {
				nodes[connection.state_to] = Node{alt, node_name}

				if !slice_string_contains(active_vertexes, connection.state_to) {
					active_vertexes = append(active_vertexes, connection.state_to)
				}
			}
		}
	}

	return -1
}

// ----------------------- Djikstra Algorithm End -----------------------

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

func slice_rune_contains(slice []rune, elem rune) bool {
	for _, slice_elem := range slice {
		if slice_elem == elem {
			return true
		}
	}

	return false
}

func slice_string_contains(slice []string, elem string) bool {
	for _, slice_elem := range slice {
		if slice_elem == elem {
			return true
		}
	}

	return false
}

func convert_state_to_code(current_key rune, keys []rune) string {
	var final_code []rune = make([]rune, 0)
	for code := FirstKeySymbol; code <= LastKeySymbol; code++ {
		final_code = append(final_code, '0')
	}

	for _, key := range keys {
		if key != AdventurerSymbol {
			key_translated := key - FirstKeySymbol
			final_code[int(key_translated)] = '1'
		}
	}

	if current_key != AdventurerSymbol {
		key_translated := current_key - FirstKeySymbol
		final_code[int(key_translated)] = '2'
	}

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
	var adventurer Adventurer = Adventurer{position_0, position_0, position_0, make(map[Position]rune), make(map[rune]Position)}

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

	//adventurer.print_mapping(true)
	distances := adventurer.compute_distances()
	restrictions := adventurer.compute_restrictions()
	map_graph := adventurer.compute_states(distances, restrictions)
	min_distance := run_dijkstra(map_graph)
	fmt.Printf("The minimum distance is ' %d ' (part 1)\n", min_distance)
}
