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

	return copy_computer
}

// ----------------------- IntCode Computer Struct End -----------------------

// ----------------------- Drone Struct Start -----------------------

type Position struct {
	x int
	y int
}

type Drone struct {
	top_left     Position
	bottom_right Position
	mapping      map[Position]int
	computer     IntCodeComputer
}

func (drone *Drone) run_computer_on(position Position) int {
	code_set, is_set := drone.mapping[position]
	if is_set {
		return code_set
	}

	copy_computer := make_deep_copy(drone.computer)
	copy_computer.input = append(copy_computer.input, position.x)
	copy_computer.input = append(copy_computer.input, position.y)

	copy_computer.run()
	var code int = copy_computer.output[len(copy_computer.output)-1]
	return code
}

func (drone *Drone) build_map() {
	var current_position Position = drone.top_left

	for current_position.y <= drone.bottom_right.y {
		var code int = drone.run_computer_on(current_position)
		drone.mapping[current_position] = code

		// Update position
		if current_position.x == drone.bottom_right.x {
			current_position.x = drone.top_left.x
			current_position.y = current_position.y + 1
		} else {
			current_position.x = current_position.x + 1
		}

	}
}

func (drone *Drone) count_pull_force_on_line(line int) int {
	var count int = 0
	var current_position Position = Position{drone.top_left.x, line}
	var CODES map[int]func(int) int = map[int]func(int) int{
		0: func(i int) int { return i },
		1: func(i int) int { return i + 1 },
	}

	not_done := true
	for not_done {
		var code int = drone.run_computer_on(current_position)
		count = CODES[code](count)

		current_position.x = current_position.x + 1
		if code == 0 && count != 0 {
			not_done = false
		}
	}

	return count
}

func (drone *Drone) first_line_with(width int) int {
	var correct_line int = -1

	var inf_limit, inf_line int = drone.count_pull_force_on_line(drone.bottom_right.y), drone.bottom_right.y
	var new_attempt_line int = width * inf_line / inf_limit

	done := false
	for !done {

		count := drone.count_pull_force_on_line(new_attempt_line)

		if count >= width && drone.count_pull_force_on_line(new_attempt_line-1) >= width {
			new_attempt_line = new_attempt_line - 1
		} else if count >= width {
			correct_line = new_attempt_line
			done = true
		} else {
			new_attempt_line = new_attempt_line + 1
		}
	}

	return correct_line
}

func (drone *Drone) position_for_box(size int) Position {
	var current_line int = drone.first_line_with(size)
	var starting_x int = 0

	for true {

		current_position := Position{starting_x, current_line}
		found_first_pull, line_over := false, false
		for !line_over {

			code := drone.run_computer_on(current_position)
			if code == 0 {
				if found_first_pull {
					// Went back to no pull -> skip to next line
					line_over = true
				} else {
					// Code 0 and still no pull
					current_position.x = current_position.x + 1
				}
			} else {
				line_passed, box_fits := drone.box_fits(current_position, size)
				if box_fits {
					// Code 1 and fits
					return current_position
				} else if !line_passed {
					// Line is over
					line_over = true
				} else {
					// Code 1 and doesn't fit
					current_position.x = current_position.x + 1
				}
			}

			// Update starting x for next line
			if code == 1 && !found_first_pull {
				found_first_pull = true
				starting_x = current_position.x
			}
		}

		current_line = current_line + 1
	}

	return Position{-1, -1}
}

func (drone *Drone) box_fits(position Position, size int) (bool, bool) {
	// Check row
	tmp_position := Position{position.x + size - 1, position.y}
	var code int = drone.run_computer_on(tmp_position)
	all_pull_row := code == 1

	// Check column
	tmp_position = Position{position.x, position.y + size - 1}
	code = drone.run_computer_on(tmp_position)
	all_pull_column := code == 1

	return all_pull_row, all_pull_row && all_pull_column
}

func (drone *Drone) count_pull_positions() int {
	var count int = 0
	var current_position Position = drone.top_left
	var CODES map[int]func(int) int = map[int]func(int) int{
		0: func(i int) int { return i },
		1: func(i int) int { return i + 1 },
	}

	for current_position.y <= drone.bottom_right.y {
		// Retrieve code converted
		count = CODES[drone.mapping[current_position]](count)

		// Update position
		if current_position.x == drone.bottom_right.x {
			current_position.x = drone.top_left.x
			current_position.y = current_position.y + 1
		} else {
			current_position.x = current_position.x + 1
		}

	}

	return count
}

func (drone *Drone) print_map() {
	var current_position Position = drone.top_left
	var CODES map[int]rune = map[int]rune{
		0: '.',
		1: '#',
	}

	for current_position.y <= drone.bottom_right.y {
		// Retrieve code converted
		element := CODES[drone.mapping[current_position]]

		// Update position
		if current_position.x == drone.bottom_right.x {
			fmt.Println()
			current_position.x = drone.top_left.x
			current_position.y = current_position.y + 1
		} else {
			fmt.Printf("%c ", element)
			current_position.x = current_position.x + 1
		}

	}
}

// ----------------------- Drone Struct End -----------------------

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

		var computer IntCodeComputer = IntCodeComputer{"booting", make([]int, 0), 0, values_converted, 0, 0, 0, make([]int, 0)}
		position_0 := Position{0, 0}
		position_till := Position{49, 49}
		var drone Drone = Drone{position_0, position_till, make(map[Position]int), computer}

		drone.build_map()
		drone.print_map()
		fmt.Println("-------------------------------------------")

		// Part 1
		count := drone.count_pull_positions()
		fmt.Printf("Pull force on: ' %d ' (part 1)\n", count)

		// Part 2
		var size int = 100
		box_position := drone.position_for_box(size)
		result := box_position.x*10000 + box_position.y
		fmt.Printf("Position ( %d, %d) fits box with ' %d ' size with result '%d ' (part 2)\n", box_position.x, box_position.y, size, result)
	}
}
