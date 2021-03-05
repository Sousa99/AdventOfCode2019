package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Position struct {
	x int
	y int
}

// ----------------------- Labyrinth Struct Start -----------------------

type Portal struct {
	name        string
	position_to Position
	portal_type string
}

type Labyrinth struct {
	top_left          Position
	bottom_right      Position
	start_portal_name string
	start_portal      Position
	end_portal_name   string
	end_portal        Position
	mapping           map[Position]string
	portals           map[Position]Portal
}

var RUNETOCODE map[rune]string = map[rune]string{
	'#': "Wall",
	'.': "FreeSpace",
}

var CODETORUNE map[string]rune = map[string]rune{
	"Wall":      '#',
	"FreeSpace": '.',
	"Nothing":   ' ',
	"Portal":    'O',
}

var PORTALTYPETOACTION map[string]string = map[string]string{
	"inner": "increment",
	"outer": "decrement",
}

func (labyrinth *Labyrinth) import_mapping(file_lines []string, main_points MainPoints) {

	// Inport map for real
	for index_line := main_points.top_left.y; index_line <= main_points.bottom_right.y; index_line++ {

		line := []rune(file_lines[index_line])
		for index_row := main_points.top_left.x; index_row <= main_points.bottom_right.x; index_row++ {

			// Check if middle of donut
			if index_line >= main_points.top_left_inner.y && index_line <= main_points.bottom_right_inner.y &&
				index_row == main_points.top_left_inner.x {
				index_row = main_points.bottom_right_inner.x + 1
			}

			position := Position{index_row - main_points.top_left.x, index_line - main_points.top_left.y}
			characther := line[index_row]
			object_to_store := RUNETOCODE[characther]

			labyrinth.mapping[position] = object_to_store

			// Update bottom right
			if position.x > labyrinth.bottom_right.x {
				labyrinth.bottom_right.x = position.x
			}
			if position.y > labyrinth.bottom_right.y {
				labyrinth.bottom_right.y = position.y
			}
		}
	}

	// Import portals
	type PortalInfo struct {
		position          Position
		portal_type       string
		inverse_portal_to Position
	}

	// Retrieve portals
	var portals map[string]PortalInfo = make(map[string]PortalInfo)
	// Horizontal Portals
	for index_line := main_points.top_left.y; index_line <= main_points.bottom_right.y; index_line++ {

		current_portal_name := ""
		line := []rune(file_lines[index_line])
		for index_row := 0; index_row < len(line); index_row++ {

			characther := line[index_row]
			if characther >= MIN_PORTAL_ID && characther <= MAX_PORTAL_ID {
				current_portal_name = current_portal_name + string(characther)

				if len(current_portal_name) == 2 {

					// Current portal formulated
					var position_portal, inverse_portal_to Position = POSITION_INVALID, POSITION_INVALID
					var portal_type string
					var current_position Position = Position{index_row - main_points.top_left.x, index_line - main_points.top_left.y}
					if index_row == main_points.top_left.x-1 {
						// Left most side
						position_portal = Position{current_position.x, current_position.y}
						inverse_portal_to = Position{current_position.x + 1, current_position.y}
						portal_type = "outer"

					} else if index_row == main_points.top_left_inner.x+1 {
						// Left inner side
						position_portal = Position{current_position.x - 1, current_position.y}
						inverse_portal_to = Position{current_position.x - 2, current_position.y}
						portal_type = "inner"

					} else if index_row == main_points.bottom_right_inner.x {
						// Right inner side
						position_portal = Position{current_position.x, current_position.y}
						inverse_portal_to = Position{current_position.x + 1, current_position.y}
						portal_type = "inner"

					} else if index_row == main_points.bottom_right.x+2 {
						// Right outer side
						position_portal = Position{current_position.x - 1, current_position.y}
						inverse_portal_to = Position{current_position.x - 2, current_position.y}
						portal_type = "outer"

					}

					// Change name for pair
					_, pair_set := portals[current_portal_name]
					if pair_set {
						current_portal_name = current_portal_name + "_pair"
					}

					// Save portal
					new_portal := PortalInfo{position_portal, portal_type, inverse_portal_to}
					portals[current_portal_name] = new_portal
					if current_portal_name == labyrinth.start_portal_name || current_portal_name == labyrinth.end_portal_name {
						labyrinth.mapping[inverse_portal_to] = "Portal"
					} else {
						labyrinth.mapping[position_portal] = "Portal"
					}
					current_portal_name = ""
				}

			} else {
				current_portal_name = ""
			}
		}
	}
	// Vertical Portals
	for index_row := main_points.top_left.x; index_row <= main_points.bottom_right.y; index_row++ {

		current_portal_name := ""
		for index_line := 0; index_line < len(file_lines); index_line++ {

			characther := []rune(file_lines[index_line])[index_row]
			if characther >= MIN_PORTAL_ID && characther <= MAX_PORTAL_ID {
				current_portal_name = current_portal_name + string(characther)

				if len(current_portal_name) == 2 {

					// Current portal formulated
					var position_portal, inverse_portal_to Position = POSITION_INVALID, POSITION_INVALID
					var portal_type string
					var current_position Position = Position{index_row - main_points.top_left.x, index_line - main_points.top_left.y}
					if index_line == main_points.top_left.y-1 {
						// Upper most side
						position_portal = Position{current_position.x, current_position.y}
						inverse_portal_to = Position{current_position.x, current_position.y + 1}
						portal_type = "outer"

					} else if index_line == main_points.top_left_inner.y+1 {
						// Upper inner side
						position_portal = Position{current_position.x, current_position.y - 1}
						inverse_portal_to = Position{current_position.x, current_position.y - 2}
						portal_type = "inner"

					} else if index_line == main_points.bottom_right_inner.y {
						// Lower inner side
						position_portal = Position{current_position.x, current_position.y}
						inverse_portal_to = Position{current_position.x, current_position.y + 1}
						portal_type = "inner"

					} else if index_line == main_points.bottom_right.y+2 {
						// Lower outer side
						position_portal = Position{current_position.x, current_position.y - 1}
						inverse_portal_to = Position{current_position.x, current_position.y - 2}
						portal_type = "outer"

					}

					// Change name for pair
					_, pair_set := portals[current_portal_name]
					if pair_set {
						current_portal_name = current_portal_name + "_pair"
					}

					// Save portal
					new_portal := PortalInfo{position_portal, portal_type, inverse_portal_to}
					portals[current_portal_name] = new_portal
					if current_portal_name == labyrinth.start_portal_name || current_portal_name == labyrinth.end_portal_name {
						labyrinth.mapping[inverse_portal_to] = "Portal"
					} else {
						labyrinth.mapping[position_portal] = "Portal"
					}
					current_portal_name = ""
				}

			} else {
				current_portal_name = ""
			}
		}
	}

	// Save back portals in labyrinth
	for portal_name, first_portal_info := range portals {
		if strings.Contains(portal_name, "_pair") {
			// Lets deal with the first only
			continue
		}

		if portal_name == labyrinth.start_portal_name {
			labyrinth.start_portal = first_portal_info.inverse_portal_to
			continue
		} else if portal_name == labyrinth.end_portal_name {
			labyrinth.end_portal = first_portal_info.inverse_portal_to
			continue
		}

		first_portal_name := portal_name + "_0"
		other_portal_name := portal_name + "_1"
		other_portal_info := portals[portal_name+"_pair"]

		labyrinth.portals[first_portal_info.position] = Portal{first_portal_name, other_portal_info.inverse_portal_to, first_portal_info.portal_type}
		labyrinth.portals[other_portal_info.position] = Portal{other_portal_name, first_portal_info.inverse_portal_to, other_portal_info.portal_type}
	}
}

func (labyrinth *Labyrinth) smallest_path(recursion bool) int {

	type GraphNode struct {
		name           string
		position_enter Position
		position_leave Position
	}

	type Explorer struct {
		current_position Position
		distance         int
	}

	var graph_nodes []GraphNode = []GraphNode{
		GraphNode{labyrinth.start_portal_name, POSITION_INVALID, labyrinth.start_portal},
		GraphNode{labyrinth.end_portal_name, labyrinth.end_portal, POSITION_INVALID},
	}

	// Add graph_nodes
	for position_enter, portal := range labyrinth.portals {
		new_node := GraphNode{portal.name, position_enter, portal.position_to}
		graph_nodes = append(graph_nodes, new_node)
	}

	var DIRECTIONS []Position = []Position{Position{0, 1}, Position{0, -1}, Position{1, 0}, Position{-1, 0}}
	// Build connections
	var connections map[string][]Connection = make(map[string][]Connection)
	for _, node_from := range graph_nodes {
		connections[node_from.name] = make([]Connection, 0)
		if node_from.name == labyrinth.end_portal_name {
			// Ovjective has no out connections
			continue
		}

		var visited_cells map[Position]bool = map[Position]bool{node_from.position_leave: true}
		var current_explorers []Explorer = []Explorer{Explorer{node_from.position_leave, 0}}

		for len(current_explorers) > 0 {
			// Still exploring

			var new_explorers []Explorer = make([]Explorer, 0)
			for _, explorer := range current_explorers {
				for _, direction := range DIRECTIONS {

					tmp_position_x := explorer.current_position.x + direction.x
					tmp_position_y := explorer.current_position.y + direction.y
					new_position := Position{tmp_position_x, tmp_position_y}

					_, visited := visited_cells[new_position]
					map_code, map_code_set := labyrinth.mapping[new_position]
					if visited || !map_code_set || map_code == "Wall" || map_code == "Unknown" {
						// Can't move through here
						continue
					}

					visited_cells[new_position] = true
					new_explorer := Explorer{new_position, explorer.distance + 1}
					new_explorers = append(new_explorers, new_explorer)

					prefix_name := strings.Split(labyrinth.portals[new_position].name, "_")[0]
					if map_code == "Portal" && (new_position == labyrinth.end_portal || !strings.HasPrefix(node_from.name, prefix_name)) && new_position != labyrinth.start_portal {
						// Valid new portal for connection
						var graph_node_to_name, action_type string
						if new_position == labyrinth.end_portal {
							// Connects to the end
							graph_node_to_name = labyrinth.end_portal_name
							action_type = "none"
						} else {
							// Other valid portal
							graph_node_to_name = labyrinth.portals[new_position].name
							action_type = PORTALTYPETOACTION[labyrinth.portals[new_position].portal_type]
						}

						new_connection := Connection{graph_node_to_name, action_type, new_explorer.distance}
						connections[node_from.name] = append(connections[node_from.name], new_connection)
					}
				}
			}

			current_explorers = new_explorers
		}
	}

	var shortest_path int
	if !recursion {
		shortest_path = run_dijkstra(labyrinth.start_portal_name, labyrinth.end_portal_name, connections)
	} else {
		shortest_path = run_bfs_rec(labyrinth.start_portal_name, labyrinth.end_portal_name, connections)
	}

	return shortest_path
}

func (labyrinth *Labyrinth) print_mapping() {
	var current_position Position = Position{labyrinth.top_left.x - 1, labyrinth.top_left.y - 1}

	for current_position.y <= labyrinth.bottom_right.y+1 {

		// Retrieve code converted
		code, is_set := labyrinth.mapping[current_position]
		if !is_set {
			code = "Nothing"
		}

		element := CODETORUNE[code]

		// Update position
		if current_position.x > labyrinth.bottom_right.x {
			fmt.Println()
			current_position.x = labyrinth.top_left.x - 1
			current_position.y = current_position.y + 1
		} else {
			fmt.Printf("%c ", element)
			current_position.x = current_position.x + 1
		}
	}
}

// ----------------------- Labyrinth Struct End -----------------------

// ----------------------- Main Points Struct Start -----------------------

var MIN_PORTAL_ID rune = 'A'
var MAX_PORTAL_ID rune = 'Z'

var POSITION_INVALID Position = Position{-1, -1}

type MainPoints struct {
	// Main square
	top_left     Position
	bottom_right Position
	// Inner square
	top_left_inner     Position
	bottom_right_inner Position
}

func calculate_points(file_lines []string) MainPoints {
	var mainPoints MainPoints = MainPoints{POSITION_INVALID, POSITION_INVALID, POSITION_INVALID, POSITION_INVALID}

	// Go through middle horizontal : left -> right
	middle_horizontal := len(file_lines) / 2
	already_saw_text_horizontal := false
	horizontal_length := len(file_lines[middle_horizontal])
	for index, elem := range file_lines[middle_horizontal] {

		if (elem < MIN_PORTAL_ID || elem > MAX_PORTAL_ID) && elem != ' ' && !already_saw_text_horizontal && index < horizontal_length/2 {
			// Went through the top
			mainPoints.top_left.x = index
			already_saw_text_horizontal = true

		} else if ((elem >= MIN_PORTAL_ID && elem <= MAX_PORTAL_ID) || elem == ' ') && already_saw_text_horizontal && index < horizontal_length/2 {
			// Finished top
			mainPoints.top_left_inner.x = index
			already_saw_text_horizontal = false

		} else if (elem < MIN_PORTAL_ID || elem > MAX_PORTAL_ID) && elem != ' ' && !already_saw_text_horizontal && index > horizontal_length/2 {
			// Went through the bottom
			mainPoints.bottom_right_inner.x = index - 1
			already_saw_text_horizontal = true

		} else if ((elem >= MIN_PORTAL_ID && elem <= MAX_PORTAL_ID) || elem == ' ') && already_saw_text_horizontal && index > horizontal_length/2 {
			// Finished bottom
			mainPoints.bottom_right.x = index - 1
			already_saw_text_horizontal = false
			break
		}
	}

	// Go through middle vertical : top -> bottom
	middle_vertical := len(file_lines[0]) / 2
	already_saw_text_vertical := false
	vertical_length := len(file_lines)
	for index, line := range file_lines {
		elem := []rune(line)[middle_vertical]

		if (elem < MIN_PORTAL_ID || elem > MAX_PORTAL_ID) && elem != ' ' && !already_saw_text_vertical && index < vertical_length/2 {
			// Went through the top
			mainPoints.top_left.y = index
			already_saw_text_vertical = true

		} else if ((elem >= MIN_PORTAL_ID && elem <= MAX_PORTAL_ID) || elem == ' ') && already_saw_text_vertical && index < vertical_length/2 {
			// Finished top
			mainPoints.top_left_inner.y = index
			already_saw_text_vertical = false

		} else if (elem < MIN_PORTAL_ID || elem > MAX_PORTAL_ID) && elem != ' ' && !already_saw_text_vertical && index > vertical_length/2 {
			// Went through the bottom
			mainPoints.bottom_right_inner.y = index - 1
			already_saw_text_vertical = true

		} else if ((elem >= MIN_PORTAL_ID && elem <= MAX_PORTAL_ID) || elem == ' ') && already_saw_text_vertical && index > vertical_length/2 {
			// Finished bottom
			mainPoints.bottom_right.y = index - 1
			already_saw_text_vertical = false
			break
		}
	}

	return mainPoints
}

// ----------------------- Main Points Struct End -----------------------

type Connection struct {
	name     string
	action   string
	distance int
}

func run_dijkstra(source string, target string, connections map[string][]Connection) int {
	type Node struct {
		distance int
		previous string
	}

	// Initialize nodes in correct form
	var nodes map[string]Node = make(map[string]Node)
	var vertexes []string = make([]string, 0)
	for node_name, _ := range connections {
		// Add node
		new_node := Node{math.MaxInt32, ""}
		nodes[node_name] = new_node
		// Add to vertexes
		vertexes = append(vertexes, node_name)
	}

	// Initialize source with dist 0
	source_node, _ := nodes[source]
	source_node.distance = 0
	nodes[source] = source_node

	// Run algorithm
	for len(vertexes) > 0 {

		// Sort and remove first element
		sort.Slice(vertexes, func(i, j int) bool { return nodes[vertexes[i]].distance < nodes[vertexes[j]].distance })
		node_name := vertexes[0]
		vertexes = vertexes[1:]

		if node_name == target {
			// We are only interest in target
			break
		}

		for _, connection := range connections[node_name] {

			alt := nodes[node_name].distance + connection.distance
			if alt < nodes[connection.name].distance {
				nodes[connection.name] = Node{alt, node_name}
			}
		}
	}

	return nodes[target].distance
}

func run_bfs_rec(source string, target string, connections map[string][]Connection) int {

	type Explorer struct {
		current_position string
		level            int
		distance         int
	}

	var ACTIONFUNCTIONS map[string]func(int) int = map[string]func(int) int{
		"none":      func(level int) int { return level },
		"increment": func(level int) int { return level + 1 },
		"decrement": func(level int) int { return level - 1 },
	}

	var reached_end_explorer Explorer = Explorer{"INVALID", -1, -1}
	var current_explorers []Explorer = []Explorer{Explorer{source, 0, 0}}
	for len(current_explorers) > 0 {

		new_explorers := make([]Explorer, 0)
		sort.Slice(current_explorers, func(i, j int) bool { return current_explorers[i].distance < current_explorers[j].distance })

		for _, explorer := range current_explorers {

			if reached_end_explorer.current_position != "INVALID" && explorer.distance >= reached_end_explorer.distance {
				// Not worth checking a explorer which already haves a greater distance
				continue
			}

			for _, connection := range connections[explorer.current_position] {

				if connection.action == "decrement" && explorer.level <= 0 {
					// Movement not allowed outter can only be accessed if level > 0
					continue
				}

				if connection.action == "none" && explorer.level != 0 {
					// Movement not allowed main ports can only be accessed by level 0
					continue
				}

				tmp_current_position := connection.name
				tmp_level := ACTIONFUNCTIONS[connection.action](explorer.level)
				tmp_distance := explorer.distance + connection.distance

				new_explorer := Explorer{tmp_current_position, tmp_level, tmp_distance}
				if connection.name != target {
					// Not the target
					new_explorers = append(new_explorers, new_explorer)

				} else if reached_end_explorer.current_position == "INVALID" || reached_end_explorer.distance > new_explorer.distance {
					// Better explorer reached it
					reached_end_explorer = new_explorer

				}
			}
		}

		current_explorers = new_explorers
	}

	return reached_end_explorer.distance
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	// Create scanner over file
	var file_text []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		file_text = append(file_text, line)
	}

	path_starts_at, path_ends_at := "AA", "ZZ"

	// Retrieve main points
	main_points := calculate_points(file_text)
	// Create Labyrinth
	position_0 := Position{0, 0}
	var labyrinth Labyrinth = Labyrinth{position_0, position_0, path_starts_at, POSITION_INVALID, path_ends_at, POSITION_INVALID, make(map[Position]string), make(map[Position]Portal)}
	labyrinth.import_mapping(file_text, main_points)
	//labyrinth.print_mapping()

	// Part 1
	smallest_path := labyrinth.smallest_path(false)
	fmt.Printf("The shortest path from ' %s ' to ' %s ' takes ' %d ' steps (part 1)\n", path_starts_at, path_ends_at, smallest_path)

	// Part 1
	smallest_path_rec := labyrinth.smallest_path(true)
	fmt.Printf("The shortest path with recursion from ' %s ' to ' %s ' takes ' %d ' steps (part 2)\n", path_starts_at, path_ends_at, smallest_path_rec)
}
