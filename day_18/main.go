package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

type StatePositions = []Position

type Vector struct {
	from string
	to   string
}

var Directions []Position = []Position{
	Position{1, 0},
	Position{-1, 0},
	Position{0, 1},
	Position{0, -1},
}

type Adventurer struct {
	// Current state
	positions StatePositions
	// Mapping
	top_left      Position
	bottom_right  Position
	mapping       map[Position]rune
	keys_position map[string]Position
}

func (adventurer *Adventurer) add_mapping_position(position Position, mapping_code rune) {
	// Add mapping
	if mapping_code == AdventurerSymbol {
		adventurer.keys_position[strconv.Itoa(len(adventurer.positions))+string(mapping_code)] = position
		adventurer.positions = append(adventurer.positions, position)
		mapping_code = FreeSymbol
	}
	adventurer.mapping[position] = mapping_code

	if mapping_code >= FirstKeySymbol && mapping_code <= LastKeySymbol {
		adventurer.keys_position[string(mapping_code)] = position
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
						distances[Vector{key, string(map_code)}] = new_explorer.distance
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
	var current_explorers []RestrictionsExplorer = make([]RestrictionsExplorer, 0, len(adventurer.positions))
	for _, initial_position := range adventurer.positions {
		new_restriction_explorer := RestrictionsExplorer{initial_position, initial_position, make([]rune, 0)}
		current_explorers = append(current_explorers, new_restriction_explorer)
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
		current_keys []string
		keys         []string
	}

	make_copy_state_explorer := func(original StateExplorer) StateExplorer {
		copy_current_keys := make([]string, 0, len(original.current_keys))
		copy_current_keys = append(copy_current_keys, original.current_keys...)
		copy_keys := make([]string, 0, len(original.keys))
		copy_keys = append(copy_keys, original.keys...)

		var copy_little StateExplorer = StateExplorer{copy_current_keys, copy_keys}
		return copy_little
	}

	/* =============== END AUXILIARY STRUCT =============== */

	var map_graph MapGraph = make(MapGraph)

	start_current_keys := make([]string, 0, len(adventurer.positions))
	for index, _ := range adventurer.positions {
		start_current_keys = append(start_current_keys, strconv.Itoa(index)+string(AdventurerSymbol))
	}
	starting_explorer := StateExplorer{start_current_keys, make([]string, 0)}
	var current_explorers []StateExplorer = []StateExplorer{starting_explorer}

	for explorer_index := 0; explorer_index < len(current_explorers); explorer_index++ {

		current_explorer := current_explorers[explorer_index]
		state_code := convert_state_to_code(current_explorer.current_keys, current_explorer.keys)
		_, map_set_for_state := map_graph[state_code]
		if map_set_for_state {
			continue
		}

		map_graph[state_code] = make([]Connection, 0)
		for going_to, restrictions := range restrictions {
			going_to_string := string(going_to)
			if slice_string_contains(current_explorer.current_keys, going_to_string) || slice_string_contains(current_explorer.keys, going_to_string) {
				// No interest in going to itself or somewhere it has been
				continue
			}

			has_keys_needed := true
			for _, door_encountered := range restrictions {
				key_needed := string(convert_door_to_key(door_encountered))
				if !slice_string_contains(current_explorer.keys, key_needed) && !slice_string_contains(current_explorer.current_keys, key_needed) {
					has_keys_needed = false
				}
			}

			if !has_keys_needed {
				// It doesn't have all the keys needed
				continue
			}

			for adventurer_index := 0; adventurer_index < len(adventurer.positions); adventurer_index++ {
				distance, possible_choice := distances[Vector{current_explorer.current_keys[adventurer_index], going_to_string}]
				if !possible_choice {
					// If current adventurer doens't connect skip
					continue
				}

				// Create new state
				new_state_explorer := make_copy_state_explorer(current_explorer)
				new_state_explorer.keys = append(new_state_explorer.keys, new_state_explorer.current_keys[adventurer_index])
				new_state_explorer.current_keys[adventurer_index] = going_to_string
				current_explorers = append(current_explorers, new_state_explorer)

				// Add connection
				new_state_code := convert_state_to_code(new_state_explorer.current_keys, new_state_explorer.keys)
				new_connection := Connection{new_state_code, distance}
				map_graph[state_code] = append(map_graph[state_code], new_connection)
			}
		}
	}

	return map_graph
}

func (adventurer *Adventurer) get_default_source_code() string {
	start_current_keys := make([]string, 0, len(adventurer.positions))
	for index, _ := range adventurer.positions {
		start_current_keys = append(start_current_keys, strconv.Itoa(index)+string(AdventurerSymbol))
	}

	return convert_state_to_code(start_current_keys, make([]string, 0))
}

func (adventurer *Adventurer) print_mapping(print_adventurer bool) {
	var current_position Position = adventurer.top_left

	// Iterate until last line
	for current_position.y <= adventurer.bottom_right.y {
		code := adventurer.mapping[current_position]
		if print_adventurer && slice_position_contains(adventurer.positions, current_position) {
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

// ----------- HEAP -----------
type HeapOfNodes []HeapNode

func (h HeapOfNodes) Len() int           { return len(h) }
func (h HeapOfNodes) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h HeapOfNodes) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HeapOfNodes) Push(x interface{}) {
	*h = append(*h, x.(HeapNode))
}

func (h *HeapOfNodes) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type HeapNode struct {
	name     string
	distance int
}

// ----------- HEAP -----------

type StateCode = string

type Connection struct {
	state_to StateCode
	distance int
}

type MapGraph = map[StateCode][]Connection

func run_dijkstra(source_code string, connections MapGraph) int {

	type Node struct {
		distance int
		previous string
		visited  bool
	}

	node_heap := &HeapOfNodes{}
	heap.Init(node_heap)
	nodes := make(map[string]Node)

	// Initialize nodes in correct form
	for node_name, _ := range connections {

		if node_name == source_code {
			continue
		}

		// Add node
		heap.Push(node_heap, HeapNode{node_name, math.MaxInt32})
		nodes[node_name] = Node{math.MaxInt32, "", false}
	}

	// Initialize source with dist 0
	heap.Push(node_heap, HeapNode{source_code, 0})
	nodes[source_code] = Node{0, "", false}

	// Run algorithm
	for node_heap.Len() > 0 {

		// Get minimum
		var heap_node HeapNode = heap.Pop(node_heap).(HeapNode)
		var node Node = nodes[heap_node.name]
		if len(connections[heap_node.name]) == 0 {
			// We are only interest in target
			return node.distance
		}

		// Don't go for visited states
		if node.visited {
			continue
		}

		node.visited = true
		nodes[heap_node.name] = node
		for _, connection := range connections[heap_node.name] {

			alt := node.distance + connection.distance
			if alt < nodes[connection.state_to].distance {
				nodes[connection.state_to] = Node{alt, heap_node.name, false}
				heap.Push(node_heap, HeapNode{connection.state_to, alt})
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

func slice_string_contains(slice []string, elem string) bool {
	for _, slice_elem := range slice {
		if slice_elem == elem {
			return true
		}
	}

	return false
}

func slice_position_contains(slice []Position, elem Position) bool {
	for _, slice_elem := range slice {
		if slice_elem == elem {
			return true
		}
	}

	return false
}

func convert_state_to_code(current_keys []string, keys []string) string {
	var number_positions int = len(current_keys)
	var final_code []rune = make([]rune, 0)
	for index := 0; index < number_positions; index++ {
		final_code = append(final_code, '0')
	}

	final_code = append(final_code, ',')
	for code := FirstKeySymbol; code <= LastKeySymbol; code++ {
		final_code = append(final_code, '0')
	}

	for _, key := range keys {
		if !strings.ContainsRune(key, AdventurerSymbol) {
			key_translated := key[0] - FirstKeySymbol
			final_code[number_positions+1+int(key_translated)] = '1'
		}
	}

	for index, current_key := range current_keys {
		if !strings.ContainsRune(current_key, AdventurerSymbol) {
			key_translated := current_key[0] - FirstKeySymbol
			final_code[number_positions+1+int(key_translated)] = '2'
		} else {
			final_code[index] = '1'
		}
	}

	return string(final_code)
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file1, _ := os.Open("input_part1.txt")
	file2, _ := os.Open("input_part2.txt")
	// Defer closing of file
	defer file1.Close()
	defer file2.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	position_0 := Position{0, 0}
	var adventurer_part1 Adventurer = Adventurer{make([]Position, 0), position_0, position_0, make(map[Position]rune), make(map[string]Position)}
	var adventurer_part2 Adventurer = Adventurer{make([]Position, 0), position_0, position_0, make(map[Position]rune), make(map[string]Position)}

	// Create scanner over file for Part 1
	scanner := bufio.NewScanner(file1)
	line_index := 0
	for scanner.Scan() {

		var line string = scanner.Text()
		for column_index, characther := range line {
			new_position := Position{column_index, line_index}
			adventurer_part1.add_mapping_position(new_position, characther)
		}

		line_index = line_index + 1
	}

	// Create scanner over file for Part 1
	scanner = bufio.NewScanner(file2)
	line_index = 0
	for scanner.Scan() {

		var line string = scanner.Text()
		for column_index, characther := range line {
			new_position := Position{column_index, line_index}
			adventurer_part2.add_mapping_position(new_position, characther)
		}

		line_index = line_index + 1
	}

	// Part 1
	//adventurer.print_mapping(true)
	distances := adventurer_part1.compute_distances()
	restrictions := adventurer_part1.compute_restrictions()
	map_graph := adventurer_part1.compute_states(distances, restrictions)
	start_code := adventurer_part1.get_default_source_code()
	min_distance := run_dijkstra(start_code, map_graph)
	fmt.Printf("The minimum distance is ' %d ' (part 1)\n", min_distance)

	// Part 2
	//adventurer.print_mapping(true)
	distances = adventurer_part2.compute_distances()
	restrictions = adventurer_part2.compute_restrictions()
	map_graph = adventurer_part2.compute_states(distances, restrictions)
	start_code = adventurer_part2.get_default_source_code()
	min_distance = run_dijkstra(start_code, map_graph)
	fmt.Printf("The minimum distance is ' %d ' (part 2)\n", min_distance)
}
