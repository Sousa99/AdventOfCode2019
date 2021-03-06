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

// ----------------------- Droid Struct Start -----------------------

type Droid struct {
	computer IntCodeComputer
	action   string
}

func (droid *Droid) run(code []string) int {

	// Change code adequately
	var START_COMMAND []int = convert_string_to_aascii(droid.action + "\n")
	var code_transformed []int = make([]int, 0)
	for _, line_of_code := range code {
		line_converted := convert_string_to_aascii(line_of_code)
		code_transformed = append(code_transformed, line_converted...)
	}

	// Start computer
	droid.computer.input = append(droid.computer.input, code_transformed...)
	droid.computer.input = append(droid.computer.input, START_COMMAND...)
	droid.computer.run()

	// Parse output
	var output_value int = droid.computer.output[len(droid.computer.output)-1]
	if output_value > 128 {
		// Successful
		return output_value
	} else {
		// Print debug information
		fmt.Println(convert_aascii_to_string(droid.computer.output))
		return -1
	}
}

// ----------------------- Droid Struct End -----------------------

func convert_string_to_aascii(line string) []int {
	slice_of_runes := []rune(line)
	var slice_of_ints []int = make([]int, 0)

	for _, element := range slice_of_runes {
		slice_of_ints = append(slice_of_ints, int(element))
	}

	return slice_of_ints
}

func convert_aascii_to_string(line []int) string {
	var slice_of_runes []rune = make([]rune, 0)
	for _, element := range line {
		slice_of_runes = append(slice_of_runes, rune(element))
	}

	return string(slice_of_runes)
}

func read_file_as_code(file_name string) []string {
	file, _ := os.Open(file_name)
	defer file.Close()

	var lines []string = make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text() + "\n"
		lines = append(lines, line)
	}

	return lines
}

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

		var computer_walk IntCodeComputer = IntCodeComputer{"booting", make([]int, 0), 0, values_converted, 0, 0, 0, make([]int, 0)}
		var computer_run IntCodeComputer = make_deep_copy(computer_walk)
		var droid_walk Droid = Droid{computer_walk, "WALK"}
		var droid_run Droid = Droid{computer_run, "RUN"}
		var lines_of_code_walk []string = read_file_as_code("code_walk.txt")
		var lines_of_code_run []string = read_file_as_code("code_run.txt")

		// Part 1
		hull_damage_walk := droid_walk.run(lines_of_code_walk)
		fmt.Printf("Hull damage whith ' %s ': ' %d ' (part 1)\n", droid_walk.action, hull_damage_walk)

		// Part 2
		hull_damage_run := droid_run.run(lines_of_code_run)
		fmt.Printf("Hull damage whith ' %s ': ' %d ' (part 2)\n", droid_run.action, hull_damage_run)
	}
}
