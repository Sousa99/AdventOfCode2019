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
			fmt.Printf("Tag value not recognized\n")
			os.Exit(1)
		}

		tags = append(tags, tag)
		value = value / 10
	}

	// Only valid opcodes
	if (opcode < 1 || opcode > 9) && opcode != 99 {
		fmt.Printf("Opcode value not recognized\n")
		os.Exit(1)
	}

	return Opcode{opcode, tags}
}

// ----------------------- Opcode Struct End -----------------------

// ----------------------- IntCode Computer Struct Start -----------------------

type IntCodeComputer struct {
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
	var halt bool = false

	for !halt {
		computer.extend_memory(computer.memory_pointer)
		var current_opcode Opcode = get_opcode(computer.memory[computer.memory_pointer])
		switch current_opcode.opcode {
		// Halting
		case 99:
			halt = true

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
			var arguments []int = computer.transform_to_arguments(current_opcode, 1, 1)
			computer.extend_memory((arguments[0]))
			computer.memory[arguments[0]] = computer.input[computer.input_pointer]

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 2
			computer.input_pointer = computer.input_pointer + 1

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

// ----------------------- IntCode Computer Struct End -----------------------

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

		// Iterate over split
		var split []string = strings.Split(scanner.Text(), ",")
		var values_converted []int = make([]int, 0)
		for _, code := range split {
			code_converted, _ := strconv.Atoi(code)
			values_converted = append(values_converted, code_converted)
		}

		// Test Mode Computer
		var inputs_test []int = make([]int, 0)
		inputs_test = append(inputs_test, 1)
		var outputs_test []int = make([]int, 0)
		var computer_test IntCodeComputer = IntCodeComputer{inputs_test, 0, values_converted, 0, 0, 0, outputs_test}

		// Sensor Boost Mode Computer
		var inputs_boost []int = make([]int, 0)
		inputs_boost = append(inputs_boost, 2)
		var outputs_boost []int = make([]int, 0)
		var computer_boost IntCodeComputer = IntCodeComputer{inputs_boost, 0, values_converted, 0, 0, 0, outputs_boost}

		// Part 1
		computer_test.run()
		var diagnostic_code_test int = computer_test.output[len(computer_test.output)-1]
		fmt.Printf("Boost Keycode in test mode: ' %d ' (part 1)\n", diagnostic_code_test)

		// Part 2
		computer_boost.run()
		var diagnostic_code_boost int = computer_boost.output[len(computer_boost.output)-1]
		fmt.Printf("Boost Keycode in sensor boost mode: ' %d ' (part 2)\n", diagnostic_code_boost)
	}
}
