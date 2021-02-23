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

// ----------------------- IntCode Computer Struct End -----------------------

// ----------------------- Game Struct End -----------------------

type Position struct {
	x int
	y int
}

type Element struct {
	name       string
	print_code string
}

var ObjectCodes map[int]Element = map[int]Element{
	0: Element{"Empty", " "},
	1: Element{"Wall", "X"},
	2: Element{"Block", "u"},
	3: Element{"HorizontalPaddle", "-"},
	4: Element{"Ball", "o"},
}

type Game struct {
	score        int
	top_left     Position
	bottom_right Position
	space        map[Position]int
	computer     IntCodeComputer
	paddle_x     int
	ball_x       int
}

func (game *Game) run_computer() {
	game.computer.run()
	output := game.computer.output

	for index := 0; index < len(output); index = index + 3 {
		var x_position int = output[index]
		var y_position int = output[index+1]
		var code int = output[index+2]

		if x_position == -1 && y_position == 0 {
			// Scoring
			game.score = code
			continue
		}

		game.space[Position{x_position, y_position}] = code

		// Update X Limits
		if x_position < game.top_left.x {
			game.top_left.x = x_position
		} else if x_position > game.bottom_right.x {
			game.bottom_right.x = x_position
		}
		// Update Y Limits
		if y_position < game.top_left.y {
			game.top_left.y = y_position
		} else if y_position > game.bottom_right.y {
			game.bottom_right.y = y_position
		}

		// Check if ball or paddle
		if ObjectCodes[code].name == "Ball" {
			game.ball_x = x_position
		} else if ObjectCodes[code].name == "HorizontalPaddle" {
			game.paddle_x = x_position
		}
	}
}

func (game *Game) run_game() int {
	var current_input int = 0
	var game_finished bool = false

	for !game_finished {
		game.computer.input = append(game.computer.input, current_input)
		game.space = make(map[Position]int)
		game.run_computer()

		if game.ball_x == game.paddle_x {
			// Stay put
			current_input = 0
		} else if game.ball_x < game.paddle_x {
			// Go left
			current_input = -1
		} else if game.ball_x > game.paddle_x {
			// Go right
			current_input = 1
		}

		// Check if game finished
		number_blocks := game.number_elements("Block")
		fmt.Printf("\033[2K\rNumber of blocks remaining: %d", number_blocks)
		game_finished = number_blocks == 0
	}
	fmt.Println()

	return game.score
}

func (game *Game) print_space() {
	for y_index := game.top_left.y; y_index <= game.bottom_right.y; y_index++ {
		for x_index := game.top_left.x; x_index <= game.bottom_right.x; x_index++ {
			code, is_set := game.space[Position{x_index, y_index}]
			if !is_set {
				// By omission is a free space
				code = 0
			}

			fmt.Printf("%s ", ObjectCodes[code].print_code)
		}

		fmt.Println()
	}
}

func (game *Game) number_elements(element string) int {
	var count int = 0
	for _, code := range game.space {
		if ObjectCodes[code].name == element {
			count = count + 1
		}
	}

	return count
}

// ----------------------- Game Struct End -----------------------

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

		var inputs []int = make([]int, 0)
		var output []int = make([]int, 0)
		var computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, values_converted, 0, 0, 0, output}
		var space map[Position]int = make(map[Position]int)
		var game Game = Game{0, Position{0, 0}, Position{0, 0}, space, computer, 0, 0}
		game.run_computer()
		//game.print_space()

		// Part 1
		var object string = "Block"
		var count int = game.number_elements(object)
		fmt.Printf("Number of ' %s ': ' %d ' (part 1)\n", object, count)
		fmt.Println("-------------------------------------------------")

		inputs = make([]int, 0)
		output = make([]int, 0)
		values_converted[0] = 2
		computer = IntCodeComputer{"booting", inputs, 0, values_converted, 0, 0, 0, output}
		space = make(map[Position]int)
		game = Game{0, Position{0, 0}, Position{0, 0}, space, computer, 0, 0}
		game.run_computer()

		// Part 2
		var score int = game.run_game()
		fmt.Printf("Final score: ' %d ' (part 2)\n", score)
	}
}
