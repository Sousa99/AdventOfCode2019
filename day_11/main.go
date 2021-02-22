package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

// ----------------------- Robot Struct End -----------------------

type Position struct {
	x int
	y int
}

var Directions = map[string]Position{
	"Up":    {0, 1},
	"Down":  {0, -1},
	"Right": {1, 0},
	"Left":  {-1, 0},
}

var DirectionString = map[string]string{
	"Up":    "^",
	"Down":  "v",
	"Right": ">",
	"Left":  "<",
}

var RotateRight = map[string]string{
	"Up":    "Right",
	"Down":  "Left",
	"Right": "Down",
	"Left":  "Up",
}
var RotateLeft = map[string]string{
	"Up":    "Left",
	"Down":  "Right",
	"Right": "Up",
	"Left":  "Down",
}

type Robot struct {
	position    Position
	direction   string
	bottom_left Position
	top_right   Position
	panel       map[Position]int
	computer    IntCodeComputer
}

func (robot *Robot) run() {
	var current_state string = robot.computer.state
	for current_state != "halted" {

		// Read current tile
		current_tile, tile_exists := robot.panel[robot.position]
		if !tile_exists {
			robot.panel[robot.position] = 0
			current_tile = 0
		}

		// Update and run robot
		robot.computer.input = append(robot.computer.input, current_tile)
		robot.computer.run()

		// Retrieve values
		current_state = robot.computer.state
		len_output := len(robot.computer.output)
		paint, turn_value := robot.computer.output[len_output-2], robot.computer.output[len_output-1]

		// Robot actuates
		robot.panel[robot.position] = paint
		if turn_value == 1 {
			// Rotate right
			robot.direction = RotateRight[robot.direction]
		} else {
			// Rotate left
			robot.direction = RotateLeft[robot.direction]
		}
		position_variance := Directions[robot.direction]
		robot.position.x = robot.position.x + position_variance.x
		robot.position.y = robot.position.y + position_variance.y

		// Update x's for printing
		if robot.position.x < robot.bottom_left.x {
			robot.bottom_left.x = robot.position.x
		} else if robot.position.x > robot.top_right.x {
			robot.top_right.x = robot.position.x
		}
		// Update y's for printing
		if robot.position.y < robot.bottom_left.y {
			robot.bottom_left.y = robot.position.y
		} else if robot.position.y > robot.top_right.y {
			robot.top_right.y = robot.position.y
		}
	}
}

func (robot *Robot) print_panel() {
	var current_position Position = Position{robot.bottom_left.x, robot.top_right.y}
	var completed bool = false

	for !completed {

		// Print out element
		panel_value, panel_set := robot.panel[current_position]
		var panel_string string
		if current_position == robot.position {
			// Robot overtop
			panel_string = DirectionString[robot.direction]
		} else if !panel_set || panel_value == 0 {
			// Not set or black
			panel_string = "."
		} else {
			// White
			panel_string = "#"
		}
		fmt.Printf("%s ", panel_string)

		// Last point in pannel
		if current_position.x == robot.top_right.x && current_position.y == robot.bottom_left.y {
			fmt.Println()
			completed = true
			continue
		}

		// Next line
		if current_position.x == robot.top_right.x {
			fmt.Println()
			current_position.x = robot.bottom_left.x
			current_position.y = current_position.y - 1
		} else {
			current_position.x = current_position.x + 1
		}
	}
}

func (robot *Robot) save_panel_as_image(file string) {
	width := robot.top_right.x - robot.bottom_left.x + 1
	height := robot.top_right.y - robot.bottom_left.y + 1

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	colorCoding := map[int]*image.Uniform{
		0: image.NewUniform(color.Black),
		1: image.NewUniform(color.White),
	}

	for line_index := robot.bottom_left.y; line_index <= robot.top_right.y; line_index++ {
		for pixel_index := robot.bottom_left.x; pixel_index <= robot.top_right.x; pixel_index++ {
			value, is_set := robot.panel[Position{pixel_index, line_index}]
			if !is_set {
				value = 0
			}

			fixed_line_index := robot.top_right.y - line_index
			fixed_pixel_index := pixel_index - robot.bottom_left.x
			img.Set(fixed_pixel_index, fixed_line_index, colorCoding[value])
		}
	}

	f, _ := os.Create(file)
	png.Encode(f, img)
}

// ----------------------- Robot Struct End -----------------------

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
		var inputs []int = make([]int, 0)
		var outputs []int = make([]int, 0)
		var robot_computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, values_converted, 0, 0, 0, outputs}
		var robot Robot = Robot{Position{0, 0}, "Up", Position{0, 0}, Position{0, 0}, make(map[Position]int), robot_computer}

		robot.run()
		//robot.print_panel()
		least_painted_cells := len(robot.panel)
		fmt.Printf("The robot painted at least ' %d ' cells (part 1)\n", least_painted_cells)

		// Part 2
		var fixed_inputs []int = make([]int, 0)
		var fixed_outputs []int = make([]int, 0)
		var fixed_robot_computer IntCodeComputer = IntCodeComputer{"booting", fixed_inputs, 0, values_converted, 0, 0, 0, fixed_outputs}
		var fixed_panel map[Position]int = make(map[Position]int)
		fixed_panel[Position{0, 0}] = 1
		var fixed_robot Robot = Robot{Position{0, 0}, "Up", Position{0, 0}, Position{0, 0}, fixed_panel, fixed_robot_computer}

		fixed_robot.run()
		fixed_robot.print_panel()
		fixed_robot.save_panel_as_image("output.png")
	}
}
