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
		if tag != 1 && tag != 0 {
			fmt.Printf("Tag value not recognized")
			os.Exit(1)
		}

		tags = append(tags, tag)
		value = value / 10
	}

	// Only valid opcodes
	if (opcode < 1 || opcode > 8) && opcode != 99 {
		fmt.Printf("Tag value not recognized")
		os.Exit(1)
	}

	return Opcode{opcode, tags}
}

// ----------------------- Opcode Struct End -----------------------

// ----------------------- IntCode Computer Struct Start -----------------------

type IntCodeComputer struct {
	state          string
	input          []int
	input_pointer  int
	memory         []int
	memory_pointer int
	output         []int
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

		var argument int = computer.memory[computer.memory_pointer+1+arg_index]

		if arg_index >= number_arg-writing_args || type_argument == 1 {
			// Immediate Mode or Writing Arg
			arguments = append(arguments, argument)
		} else if type_argument == 0 {
			// Position Mode
			var argument_value int = computer.memory[argument]
			arguments = append(arguments, argument_value)
		}
	}

	return arguments
}

func (computer *IntCodeComputer) run() {
	computer.state = "running"
	for computer.state == "running" {
		var current_opcode Opcode = get_opcode(computer.memory[computer.memory_pointer])
		switch current_opcode.opcode {
		// Halting
		case 99:
			computer.state = "halted"

		// Addition
		case 1:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
			computer.memory[arguments[2]] = arguments[0] + arguments[1]

			//Advance pointer
			computer.memory_pointer = computer.memory_pointer + 4

		// Multiplication
		case 2:
			var arguments []int = computer.transform_to_arguments(current_opcode, 3, 1)
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
				computer.memory[arguments[0]] = computer.input[computer.input_pointer]

				//Advance pointers
				computer.memory_pointer = computer.memory_pointer + 2
				computer.input_pointer = computer.input_pointer + 1
			}

		// Output
		case 4:
			var arguments []int = computer.transform_to_arguments(current_opcode, 1, 1)
			computer.output = append(computer.output, computer.memory[arguments[0]])

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
			if arguments[0] == arguments[1] {
				computer.memory[arguments[2]] = 1
			} else {
				computer.memory[arguments[2]] = 0
			}

			//Advance pointers
			computer.memory_pointer = computer.memory_pointer + 4

		// Should not happen
		default:
			fmt.Println("Code not recognized")
			os.Exit(1)
		}
	}
}

// ----------------------- IntCode Computer Struct End -----------------------

// ----------------------- Amplifier Controller Struct End -----------------------

type AmplifierController struct {
	number_amplifiers int
	minimum_phase     int
	maximum_phase     int
	first_input       int
	code              []int
}

func (controller *AmplifierController) run_with_phase(phase_setting []int) int {
	var current_input = controller.first_input
	for _, phase_value := range phase_setting {
		var inputs []int = make([]int, 0, 2)
		inputs = append(inputs, phase_value)
		inputs = append(inputs, current_input)
		var memory []int = make([]int, 0)
		memory = append(memory, controller.code...)
		var output []int = make([]int, 0)

		var computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, memory, 0, output}
		computer.run()
		current_input = computer.output[len(computer.output)-1]
	}

	return current_input
}

func (controller *AmplifierController) run_with_phase_with_feedback(phase_setting []int) int {
	var amplifiers []IntCodeComputer = make([]IntCodeComputer, 0, 6)

	// Setup amplifiers
	for _, phase_value := range phase_setting {
		var inputs []int = make([]int, 0)
		inputs = append(inputs, phase_value)
		var memory []int = make([]int, 0)
		memory = append(memory, controller.code...)
		var output []int = make([]int, 0)

		var computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, memory, 0, output}
		amplifiers = append(amplifiers, computer)
	}

	// Run amplifiers
	var current_input = controller.first_input
	var halted bool = false

	for !halted {
		for index, _ := range amplifiers {
			// Retrieve amplifier
			amplifier := amplifiers[index]
			amplifier.input = append(amplifier.input, current_input)

			amplifier.run()
			current_input = amplifier.output[len(amplifier.output)-1]

			halted = halted || amplifier.state == "halted"
			// Save new version of amplifier
			amplifiers[index] = amplifier
		}

	}

	return current_input
}

func (controller *AmplifierController) get_maximum_thrust(with_feedback bool) ([]int, int) {
	var phase_setting []int = make([]int, 0)
	// Initialize phase setting
	for i := controller.minimum_phase; i <= controller.maximum_phase; i++ {
		phase_setting = append(phase_setting, i)
	}

	// Storing of maximums
	var maximum_thrust int = 0
	var maximum_phase_setting []int = make([]int, 0)

	var permutations [][]int = permutations(phase_setting)
	for _, permutation := range permutations {

		// Run Intcode
		var thrust_value int
		if !with_feedback {
			thrust_value = controller.run_with_phase(permutation)
		} else {
			thrust_value = controller.run_with_phase_with_feedback(permutation)
		}

		// Check if maximum
		if thrust_value > maximum_thrust {
			maximum_thrust = thrust_value
			maximum_phase_setting = permutation
		}
	}

	return maximum_phase_setting, maximum_thrust
}

// ----------------------- Amplifier Controller Struct End -----------------------

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
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

		// Iterate over split
		var split []string = strings.Split(scanner.Text(), ",")
		var values_converted []int = make([]int, 0)
		for _, code := range split {
			code_converted, _ := strconv.Atoi(code)
			values_converted = append(values_converted, code_converted)
		}

		// Part 1
		var amplifier_controller AmplifierController = AmplifierController{5, 0, 4, 0, values_converted}
		phase_setting, thrust_value := amplifier_controller.get_maximum_thrust(false)
		fmt.Println(phase_setting)
		fmt.Printf("Maximum thrust possible without feedback: ' %d ' (part 1)\n", thrust_value)

		fmt.Println("------------------------------------------")

		// Part 2
		var feedback_amplifier_controller AmplifierController = AmplifierController{5, 5, 9, 0, values_converted}
		phase_setting, thrust_value = feedback_amplifier_controller.get_maximum_thrust(true)
		fmt.Println(phase_setting)
		fmt.Printf("Maximum thrust possible with feedback: ' %d ' (part 2)\n", thrust_value)
	}
}
