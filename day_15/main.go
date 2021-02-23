package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ----------------------- Opcode Struct Start -----------------------

type Opcode struct {
	opcode int
	tags   []int
}

func get_opcode(value int) Opcode {
	var opcode int = value % 100
	var tags []int = make([]int, 0)
	value = value / 100

	for value != 0 {
		var tag int = value % 10
		// Only valid tags
		if tag < 0 || tag > 2 {
			fmt.Printf("Tag value not recognized: ' %d '\n", tag)
			os.Exit(1)
		}

		tags = append(tags, tag)
		value = value / 10
	}

	// Only valid opcodes
	if (opcode < 1 || opcode > 9) && opcode != 99 {
		fmt.Printf("Opcode value not recognized: ' %d '\n", opcode)
		os.Exit(1)
	}

	return Opcode{opcode, tags}
}

// ----------------------- Opcode Struct End -----------------------

// ----------------------- IntCode Computer Struct Start -----------------------

type IntCodeComputer struct {
	state                string
	input                []int
	input_pointer        int
	memory               []int
	memory_pointer       int
	memory_default_value int
	relative_pointer     int
	output               []int
}

func (computer *IntCodeComputer) extend_memory(position int) {
	var missing_entries int = position - (len(computer.memory) - 1)
	for i := 0; i < missing_entries; i++ {
		computer.memory = append(computer.memory, computer.memory_default_value)
	}
}

func (computer *IntCodeComputer) transform_to_arguments(opcode Opcode, number_arg int, writing_args int) []int {
	var arguments []int = make([]int, 0, number_arg)
	// Iterate over arguments
	for arg_index := 0; arg_index < number_arg; arg_index++ {

		var type_argument int
		if arg_index >= len(opcode.tags) {
			// By omission type is 0
			type_argument = 0
		} else {
			// Else type is given by opcode tag
			type_argument = opcode.tags[arg_index]
		}

		computer.extend_memory(computer.memory_pointer + 1 + arg_index)
		var argument int = computer.memory[computer.memory_pointer+1+arg_index]

		if arg_index >= number_arg-writing_args {
			// Writing position
			if type_argument == 2 {
				// Relative mode writing
				var argument_value int = computer.relative_pointer + argument
				arguments = append(arguments, argument_value)
			} else if type_argument == 0 {
				// Position mode writing
				arguments = append(arguments, argument)
			}
		} else {
			// Reading position
			if type_argument == 2 {
				// Relative mode
				var position int = computer.relative_pointer + argument
				computer.extend_memory(position)
				var argument_value int = computer.memory[position]
				arguments = append(arguments, argument_value)
			} else if type_argument == 1 {
				// Immediate Mode
				arguments = append(arguments, argument)
			} else if type_argument == 0 {
				// Position Mode
				computer.extend_memory(argument)
				var argument_value int = computer.memory[argument]
				arguments = append(arguments, argument_value)
			}
		}
	}

	return arguments
}

func (computer *IntCodeComputer) run() {
	computer.state = "running"
	for computer.state == "running" {

		computer.extend_memory(computer.memory_pointer)
		var current_opcode Opcode = get_opcode(computer.memory[computer.memory_pointer])
		switch current_opcode.opcode {
		// Halting
		case 99:
			computer.state = "halted"

		// Addition
		case 1:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
			computer.extend_memory(arguments[2])
			computer.memory[arguments[2]] = arguments[0] + arguments[1]

			//Advance pointer
			computer.memory_pointer = computer.memory_pointer + 4

		// Multiplication
		case 2:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
			computer.extend_memory(arguments[2])
			computer.memory[arguments[2]] = arguments[0] * arguments[1]

			//Advance pointer
			computer.memory_pointer = computer.memory_pointer + 4

		// Input
		case 3:
			if computer.input_pointer >= len(computer.input) {
				// Computer must wait for more input
				computer.state = "paused"
			} else {
				var arguments []int = computer.transform_to_arguments(current_opcode, 1, 1)
				computer.extend_memory((arguments[0]))
				computer.memory[arguments[0]] = computer.input[computer.input_pointer]

				//Advance pointers
				computer.memory_pointer = computer.memory_pointer + 2
				computer.input_pointer = computer.input_pointer + 1
			}

		// Output
		case 4:
			var arguments []int = computer.transform_to_arguments(current_opcode, 1, 0)
			computer.output = append(computer.output, arguments[0])

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 2
			computer.state = "paused"

		// Jump if True
		case 5:
			var arguments []int = computer.transform_to_arguments(current_opcode, 2, 0)
			if arguments[0] != 0 {
				computer.memory_pointer = arguments[1]
			} else {
				computer.memory_pointer = computer.memory_pointer + 3
			}

		// Jump if False
		case 6:
			var arguments []int = computer.transform_to_arguments(current_opcode, 2, 0)
			if arguments[0] == 0 {
				computer.memory_pointer = arguments[1]
			} else {
				computer.memory_pointer = computer.memory_pointer + 3
			}

		// Less than
		case 7:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
			computer.extend_memory(arguments[2])
			if arguments[0] < arguments[1] {
				computer.memory[arguments[2]] = 1
			} else {
				computer.memory[arguments[2]] = 0
			}

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 4

		// Equals
		case 8:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
			computer.extend_memory(arguments[2])
			if arguments[0] == arguments[1] {
				computer.memory[arguments[2]] = 1
			} else {
				computer.memory[arguments[2]] = 0
			}

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 4

		// Adjust relative base
		case 9:
			var arguments []int = computer.transform_to_arguments(current_opcode, 1, 0)
			computer.relative_pointer = computer.relative_pointer + arguments[0]

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 2

		// Should not happen
		default:
			fmt.Println("Code not recognized")
			os.Exit(1)
		}
	}
}

func make_deep_copy(computer IntCodeComputer) IntCodeComputer {
	copy_computer := computer

	copy_computer.input = make([]int, len(computer.input))
	copy_computer.memory = make([]int, len(computer.memory))
	copy_computer.output = make([]int, len(computer.output))

	copy(copy_computer.input, computer.input)
	copy(copy_computer.memory, computer.memory)
	copy(copy_computer.output, computer.output)

	//fmt.Printf("%p\n", computer.memory)
	//fmt.Printf("%p\n", copy_computer.memory)
	return copy_computer
}

// ----------------------- IntCode Computer Struct End -----------------------

// ----------------------- Droid Struct Start -----------------------

type Position struct {
	x int
	y int
}

type MapPoint struct {
	value    string
	distance int
}

type Direction struct {
	name     string
	movement Position
}

var DirectionCodes map[int]Direction = map[int]Direction{
	1: Direction{"North", Position{0, 1}},
	2: Direction{"South", Position{0, -1}},
	3: Direction{"West", Position{-1, 0}},
	4: Direction{"East", Position{1, 0}},
}

var DirectionReverse map[int]int = map[int]int{
	1: 2,
	2: 1,
	3: 4,
	4: 3,
}

var ObjectCodes map[string]string = map[string]string{
	"Unknown":      " ",
	"Wall":         "#",
	"Droid":        "D",
	"FreeSpace":    ".",
	"OxygenSystem": "0",
	"Oxygenated":   "O",
	"Initial":      "x",
}

type SubDroid struct {
	position Position
	computer IntCodeComputer
}

type Droid struct {
	saved_position Position
	top_right      Position
	bottom_left    Position
	mapping        map[Position]MapPoint
	saved_computer IntCodeComputer
}

func (droid *Droid) run_droid_until_oxygen() (Position, int) {
	var initial_subdroid SubDroid = SubDroid{droid.saved_position, make_deep_copy(droid.saved_computer)}
	var subdroids []SubDroid = make([]SubDroid, 0)
	subdroids = append(subdroids, initial_subdroid)

	for len(subdroids) != 0 {

		var new_subdroids []SubDroid = make([]SubDroid, 0)
		// Iterate over sub_droids
		for _, subdroid := range subdroids {
			current_distance := droid.mapping[subdroid.position].distance

			// Check every direction
			for direction_code, direction := range DirectionCodes {
				tmp_position_x := subdroid.position.x + direction.movement.x
				tmp_position_y := subdroid.position.y + direction.movement.y
				new_position := Position{tmp_position_x, tmp_position_y}

				code_stored, is_set := droid.mapping[new_position]
				if is_set && code_stored.value != "Unknown" && code_stored.distance <= current_distance+1 {
					continue
				}

				// Valid direction to explore and move robot
				var new_subdroid SubDroid = SubDroid{subdroid.position, make_deep_copy(subdroid.computer)}
				new_subdroid.computer.input = append(new_subdroid.computer.input, direction_code)
				new_subdroid.computer.run()
				var status_code int = new_subdroid.computer.output[len(new_subdroid.computer.output)-1]
				new_subdroid.computer.output = make([]int, 0)

				switch status_code {
				case 0:
					// Hits a wall
					droid.mapping[new_position] = MapPoint{"Wall", -1}
				case 1:
					droid.mapping[new_position] = MapPoint{"FreeSpace", current_distance + 1}
					new_subdroid.position = new_position
					new_subdroids = append(new_subdroids, new_subdroid)
				case 2:
					droid.mapping[new_position] = MapPoint{"OxygenSystem", current_distance + 1}
					droid.saved_position = new_position
					droid.saved_computer = new_subdroid.computer
				}

				// Update limits for printing
				if new_position.x < droid.bottom_left.x {
					droid.bottom_left.x = new_position.x
				} else if new_position.x > droid.top_right.x {
					droid.top_right.x = new_position.x
				}
				if new_position.y < droid.bottom_left.y {
					droid.bottom_left.y = new_position.y
				} else if new_position.y > droid.top_right.y {
					droid.top_right.y = new_position.y
				}
			}
		}

		subdroids = new_subdroids
	}

	return droid.saved_position, droid.mapping[droid.saved_position].distance
}

func (droid *Droid) run_droid_to_oxigenate() int {
	var searching_positions []Position = make([]Position, 0)
	searching_positions = append(searching_positions, droid.saved_position)
	var time int = 0
	var completely_oxigenated bool = false

	for !completely_oxigenated {

		var new_search_positions []Position = make([]Position, 0)
		// Iterate over sub_droids
		for _, search_position := range searching_positions {

			// Check every direction
			for _, direction := range DirectionCodes {
				tmp_position_x := search_position.x + direction.movement.x
				tmp_position_y := search_position.y + direction.movement.y
				new_position := Position{tmp_position_x, tmp_position_y}

				code_stored, is_set := droid.mapping[new_position]
				if !is_set || code_stored.value != "FreeSpace" {
					continue
				}

				// Valid direction to explore and move position
				current_distance := droid.mapping[new_position].distance
				droid.mapping[new_position] = MapPoint{"Oxygenated", current_distance}
				new_search_positions = append(new_search_positions, new_position)
			}
		}

		searching_positions = new_search_positions
		time = time + 1
		completely_oxigenated = droid.count_number_object("FreeSpace") == 0
	}

	return time
}

func (droid *Droid) count_number_object(object string) int {
	var count int = 0
	for _, map_point := range droid.mapping {
		if map_point.value == object {
			count = count + 1
		}
	}

	return count
}

func (droid *Droid) print_mapping(show_droid bool) {
	var current_position Position = Position{droid.bottom_left.x, droid.top_right.y}
	for current_position.y >= droid.bottom_left.y {
		// Print mapping element
		object, is_set := droid.mapping[current_position]
		if !is_set {
			object.value = "Unknown"
		}
		if show_droid && droid.saved_position == current_position {
			object.value = "Droid"
		}
		fmt.Printf("%s ", ObjectCodes[object.value])

		// Update current_position
		if current_position.x == droid.top_right.x {
			current_position.x = droid.bottom_left.x
			current_position.y = current_position.y - 1
			fmt.Println()
		} else {
			current_position.x = current_position.x + 1
		}
	}
}

// ----------------------- Droid Struct End -----------------------

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		var split []string = strings.Split(line, ",")
		var values_converted []int = make([]int, 0)
		for _, code := range split {
			code_converted, _ := strconv.Atoi(code)
			values_converted = append(values_converted, code_converted)
		}

		inputs := make([]int, 0)
		outputs := make([]int, 0)
		var computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, values_converted, 0, 0, 0, outputs}

		position_0 := Position{0, 0}
		mapping := make(map[Position]MapPoint)
		mapping[position_0] = MapPoint{"FreeSpace", 0}
		var droid Droid = Droid{position_0, position_0, position_0, mapping, computer}

		// Part 1
		_, distance := droid.run_droid_until_oxygen()
		fmt.Printf("The minimum distance for the OxygenSystem: ' %d ' (part 1)\n", distance)

		// Part 2
		time := droid.run_droid_to_oxigenate()
		fmt.Printf("To fully oxigenate the space: ' %d ' minutes (part 2)\n", time)

	}
}
