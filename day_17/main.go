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

// ----------------------- InterfaceASCII Struct Start -----------------------

type Object = string

var ConvertCode map[int]Object = map[int]Object{
	46:  "free",
	35:  "scaffold",
	94:  "droid_up",
	118: "droid_down",
	62:  "droid_right",
	60:  "droid_left",
}

var ObjectCodes map[Object]string = map[string]string{
	"free":         ".",
	"scaffold":     "#",
	"intersection": "O",
	"droid_up":     "^",
	"droid_down":   "v",
	"droid_right":  ">",
	"droid_left":   "<",
	"visited":      "X",
}

type Position struct {
	x int
	y int
}

var DirectionCode map[Object]Position = map[Object]Position{
	"droid_up":    Position{0, -1},
	"droid_down":  Position{0, 1},
	"droid_right": Position{1, 0},
	"droid_left":  Position{-1, 0},
}

var RotateRight map[Object]Object = map[Object]Object{
	"droid_up":    "droid_right",
	"droid_down":  "droid_left",
	"droid_right": "droid_down",
	"droid_left":  "droid_up",
}
var RotateLeft map[Object]Object = map[Object]Object{
	"droid_up":    "droid_left",
	"droid_down":  "droid_right",
	"droid_right": "droid_up",
	"droid_left":  "droid_down",
}

type InterfaceASCII struct {
	top_left       Position
	bottom_right   Position
	droid_position Position
	droid_type     string
	mapping        map[Position]Object
	computer       IntCodeComputer
}

func (ascii *InterfaceASCII) build_map() {
	ascii.computer.run()
	output := ascii.computer.output

	index_cutoff, last_new_line := -1, false
	var x_position, y_position int = 0, 0
	for index, code := range output {

		// Update limits
		if x_position > ascii.bottom_right.x {
			ascii.bottom_right.x = x_position
		}
		if y_position > ascii.bottom_right.y {
			ascii.bottom_right.y = y_position
		}

		code_converted := ConvertCode[code]
		// Converting code to symbol
		if code == 10 {
			// New line
			if last_new_line {
				index_cutoff = index
				break
			}

			x_position = 0
			y_position = y_position + 1
			last_new_line = true
		} else if code_converted == "droid_up" || code_converted == "droid_down" || code_converted == "droid_left" || code_converted == "droid_right" {
			// Droid
			last_new_line = false
			ascii.droid_type = code_converted
			ascii.droid_position = Position{x_position, y_position}
			ascii.mapping[Position{x_position, y_position}] = "scaffold"
			x_position = x_position + 1

		} else {
			// Else place object
			last_new_line = false
			ascii.mapping[Position{x_position, y_position}] = code_converted
			x_position = x_position + 1
		}
	}

	ascii.computer.output = ascii.computer.output[index_cutoff+1:]
}

func (ascii *InterfaceASCII) compute_intersections() int {

	var value int = 0
	for position, object := range ascii.mapping {
		if object != "free" {
			object_up, is_set_up := ascii.mapping[Position{position.x, position.y + 1}]
			object_down, is_set_down := ascii.mapping[Position{position.x, position.y - 1}]
			object_right, is_set_right := ascii.mapping[Position{position.x + 1, position.y}]
			object_left, is_set_left := ascii.mapping[Position{position.x - 1, position.y}]

			all_set := is_set_up && is_set_down && is_set_right && is_set_left
			if !all_set {
				continue
			}

			all_not_free := object_up != "free" && object_down != "free" && object_right != "free" && object_left != "free"
			if !all_not_free {
				continue
			}

			ascii.mapping[position] = "intersection"
			value = value + position.x*position.y
		}
	}

	return value
}

func (ascii *InterfaceASCII) compute_trajectory() []string {
	RIGHT := "R"
	LEFT := "L"
	ascii.mapping[ascii.droid_position] = "visited"

	var trajectory []string = make([]string, 0)
	var number_set int = -1
	for ascii.count_unvisited_cells() != 0 {

		// Try to move forward
		direction := DirectionCode[ascii.droid_type]
		temp_x := ascii.droid_position.x + direction.x
		temp_y := ascii.droid_position.y + direction.y
		temp_position := Position{temp_x, temp_y}
		value, is_set := ascii.mapping[temp_position]
		if is_set && value != "free" {
			// Move forward
			ascii.droid_position = temp_position
			ascii.mapping[ascii.droid_position] = "visited"
			if number_set != -1 {
				number_set = number_set + 1
			} else {
				number_set = 1
			}
			continue
		}

		// Need to rotate
		if number_set != -1 {
			trajectory = append(trajectory, strconv.Itoa(number_set))
		}
		number_set = -1

		rotate_right_direction := DirectionCode[RotateRight[ascii.droid_type]]
		temp_x = ascii.droid_position.x + rotate_right_direction.x
		temp_y = ascii.droid_position.y + rotate_right_direction.y
		temp_position = Position{temp_x, temp_y}
		value, is_set = ascii.mapping[temp_position]
		if is_set && value != "free" {
			// Rotate Right
			ascii.droid_type = RotateRight[ascii.droid_type]
			trajectory = append(trajectory, RIGHT)
			continue
		}

		rotate_left_direction := DirectionCode[RotateLeft[ascii.droid_type]]
		temp_x = ascii.droid_position.x + rotate_left_direction.x
		temp_y = ascii.droid_position.y + rotate_left_direction.y
		temp_position = Position{temp_x, temp_y}
		value, is_set = ascii.mapping[temp_position]
		if is_set && value != "free" {
			// Rotate Left
			ascii.droid_type = RotateLeft[ascii.droid_type]
			trajectory = append(trajectory, LEFT)
			continue
		}

		fmt.Println("Fuck it needs to turn back")
	}

	if number_set != -1 {
		trajectory = append(trajectory, strconv.Itoa(number_set))
	}

	return trajectory
}

func (ascii *InterfaceASCII) print_map() {
	var current_position Position = Position{0, 0}

	// Iterate until last line
	for current_position.y <= ascii.bottom_right.y {
		code := ascii.mapping[current_position]
		if ascii.droid_position == current_position {
			code = ascii.droid_type
		}

		fmt.Printf("%s ", ObjectCodes[code])

		if current_position.x == ascii.bottom_right.x {
			// Last pixel on line
			fmt.Println()
			current_position.x = ascii.top_left.x
			current_position.y = current_position.y + 1
		} else {
			// Not the last pixel
			current_position.x = current_position.x + 1
		}
	}
}

func (ascii *InterfaceASCII) count_unvisited_cells() int {
	var count int = 0

	// Iterate until last line
	for _, code := range ascii.mapping {
		if code == "scaffold" || code == "intersection" {
			count = count + 1
		}
	}

	return count
}

func (ascii *InterfaceASCII) start_moving(codification []string, patterns [][]string) int {
	var NEW_LINE int = 10
	var COMMA int = int(',')
	var VIDEO_FEED int = int('n')

	var input [][]int = make([][]int, 0)

	// Transform codification
	codification_transformed := transform_codification_to_input(codification, COMMA, NEW_LINE)
	input = append(input, codification_transformed)

	// Add patterns to input
	for _, pattern := range patterns {
		pattern_transformed := transform_codification_to_input(pattern, COMMA, NEW_LINE)
		input = append(input, pattern_transformed)
	}

	// No video feed back
	video_feed := make([]int, 0)
	video_feed = append(video_feed, VIDEO_FEED)
	video_feed = append(video_feed, NEW_LINE)
	input = append(input, video_feed)

	// Start computer
	for _, input_line := range input {
		print_output_string(ascii.computer.output)
		ascii.computer.output = make([]int, 0)
		fmt.Printf("%+v\n", input_line)
		ascii.computer.input = append(ascii.computer.input, input_line...)
		ascii.computer.run()
	}

	return ascii.computer.output[len(ascii.computer.output)-1]
}

// ----------------------- InterfaceASCII Struct End -----------------------

func get_codification(trajectory []string, max_line_length int, codes []string) ([]string, [][]string) {
	code_length := make([]int, 0, 3)
	for index := 0; index < len(codes); index++ {
		code_length = append(code_length, 1)
	}

	var not_done bool = false
	for !not_done {

		copy_trajectory := make([]string, 0, len(trajectory))
		copy_trajectory = append(copy_trajectory, trajectory...)

		parameters := make([][]string, 3)
		for parameter_index := 0; parameter_index < len(codes); parameter_index++ {
			current_param := make([]string, 0)

			current_index := 0
			// Skip codes set
			for current_index < len(copy_trajectory) && check_code_in_codes(copy_trajectory[current_index], codes) {
				current_index = current_index + 1
			}

			// Create param
			param_size := 0
			for current_index < len(copy_trajectory) && !check_code_in_codes(copy_trajectory[current_index], codes) && len(current_param) < code_length[parameter_index] && param_size+len(current_param)-1 < max_line_length {
				element := copy_trajectory[current_index]
				current_param = append(current_param, element)
				param_size = param_size + len(element)
				current_index = current_index + 1
			}

			// "Remove" from trajectory
			copy_trajectory = replace_in_trajectory(copy_trajectory, current_param, codes[parameter_index])
			parameters[parameter_index] = current_param
		}

		simplified_codification := simplify_for_codification(copy_trajectory)
		size_codification := len(simplified_codification)
		if check_trajectory_completed(copy_trajectory, codes) && size_codification+size_codification-1 < max_line_length {
			fmt.Printf("Trajectory: \t%v\n", copy_trajectory)
			fmt.Printf("Simplified: \t%v\n", simplified_codification)
			for index, param := range parameters {
				fmt.Printf("Parameter %s: \t%v\n", codes[index], param)
			}

			return simplified_codification, parameters
		}

		// Increment code_length
		not_done = true
		for index := 0; index < len(code_length); index++ {
			if code_length[index]+(code_length[index]-1) < max_line_length {
				code_length[index] = code_length[index] + 1
				not_done = false
				break
			} else {
				code_length[index] = 1
			}
		}
	}

	fmt.Println("Not a single codification was found")
	return make([]string, 0), make([][]string, 0)
}

func replace_in_trajectory(trajectory []string, pattern []string, code string) []string {
	for index := 0; index < len(trajectory)-len(pattern)+1; index++ {

		valid_pattern := true
		for pattern_index := 0; pattern_index < len(pattern); pattern_index++ {
			if trajectory[index+pattern_index] != pattern[pattern_index] {
				valid_pattern = false
				break
			}
		}

		// If valid sum to count and check
		if valid_pattern {
			// Replace in trajectory
			for pattern_index := 0; pattern_index < len(pattern); pattern_index++ {
				trajectory[index+pattern_index] = code
			}
			// Update index
			index = index - 1 + len(pattern)
		}
	}

	return trajectory
}

func simplify_for_codification(trajectory []string) []string {
	copy_trajectory := make([]string, 0)
	copy_trajectory = append(copy_trajectory, trajectory...)

	index := 1
	for index < len(copy_trajectory) {
		if copy_trajectory[index-1] == copy_trajectory[index] {
			new_trajectory := make([]string, 0)
			new_trajectory = append(new_trajectory, copy_trajectory[:index]...)
			if index+1 < len(copy_trajectory) {
				new_trajectory = append(new_trajectory, copy_trajectory[index+1:]...)
			}
			copy_trajectory = new_trajectory

		} else {
			index = index + 1
		}
	}

	return copy_trajectory
}

func check_trajectory_completed(trajectory []string, codes []string) bool {
	for _, code := range trajectory {
		if !check_code_in_codes(code, codes) {
			return false
		}
	}

	return true
}

func check_code_in_codes(code string, codes []string) bool {
	for _, valid_code := range codes {
		if code == valid_code {
			return true
		}
	}

	return false
}

func transform_codification_to_input(codification []string, sep int, end int) []int {
	transformed := make([]int, 0)

	for _, code := range codification {
		for _, letter_code := range code {
			transformed = append(transformed, int(letter_code))
		}
		transformed = append(transformed, sep)
	}
	transformed[len(transformed)-1] = end

	return transformed
}

func print_output_string(output []int) {
	var final string = ""
	for _, code := range output {
		final = final + string(rune(code))
	}

	fmt.Printf("- %s", final)
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

		inputs := make([]int, 0)
		outputs := make([]int, 0)
		var computer IntCodeComputer = IntCodeComputer{"booting", inputs, 0, values_converted, 0, 0, 0, outputs}

		position_0 := Position{0, 0}
		var ascii InterfaceASCII = InterfaceASCII{position_0, position_0, position_0, "unknown", make(map[Position]string), computer}

		// Part 1
		ascii.build_map()
		alignment := ascii.compute_intersections()
		//ascii.print_map()
		fmt.Printf("Sum of Alignment Parameters: ' %d ' (part 1)\n", alignment)
		fmt.Println("--------------------------------------------")

		// Part 2
		trajectory := ascii.compute_trajectory()
		fmt.Printf("Trajectory: \t%v\n", trajectory)

		var max_size int = 20
		var codes []string = []string{"A", "B", "C"}
		codification, parameters := get_codification(trajectory, max_size, codes)
		fmt.Println("--------------------------------------------")
		dust := ascii.start_moving(codification, parameters)
		fmt.Println("--------------------------------------------")
		fmt.Printf("Dust collected: ' %d ' (part 2)\n", dust)
	}
}
