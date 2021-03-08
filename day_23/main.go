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
				computer.state = "awaiting input"
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

// ----------------------- CommunicatingModule Struct Start -----------------------

type Packet struct {
	to_id int
	x     int
	y     int
}

type CommunicatingModule struct {
	id             int
	computer       IntCodeComputer
	cache_received []Packet
}

func (module *CommunicatingModule) add_packet(new_packet Packet) {
	module.cache_received = append(module.cache_received, new_packet)
}

// ----------------------- CommunicatingModule Struct End -----------------------

// ----------------------- NATModule Struct Start -----------------------

var INVALID_PACKET Packet = Packet{-1, -1, -1}

type NATModule struct {
	id      int
	send_to int
	packet  Packet
}

// ----------------------- NATModule Struct End -----------------------

// ----------------------- System Struct Start -----------------------

func new_system(number_modules int, mock_computer IntCodeComputer, target_id int, when_idle_send int) System {
	var new_nat NATModule = NATModule{target_id, when_idle_send, INVALID_PACKET}
	var new_system System = System{new_nat, make(map[int]CommunicatingModule)}

	for module_id := 0; module_id < number_modules; module_id++ {

		copy_computer := make_deep_copy(mock_computer)
		copy_computer.input = append(copy_computer.input, module_id)

		new_module := CommunicatingModule{module_id, copy_computer, make([]Packet, 0)}
		new_system.modules[module_id] = new_module
	}

	return new_system
}

type System struct {
	nat     NATModule
	modules map[int]CommunicatingModule
}

func (system *System) run(ids []int) bool {
	var idle bool = true
	// Actually running
	for _, module_id := range ids {

		var module CommunicatingModule = system.modules[module_id]
		module.computer.run()
		// If waiting for input
		if module.computer.state == "awaiting input" {

			if len(module.cache_received) > 0 {
				// Packets waiting processing
				packet := module.cache_received[0]
				module.cache_received = module.cache_received[1:]
				module.computer.input = append(module.computer.input, packet.x)
				module.computer.input = append(module.computer.input, packet.y)

			} else {
				// No packets for processing
				module.computer.input = append(module.computer.input, -1)
			}
		}

		// If has output, process it
		if len(module.computer.output) >= 3 {
			for len(module.computer.output) >= 3 {

				new_packet := Packet{module.computer.output[0], module.computer.output[1], module.computer.output[2]}
				module.computer.output = module.computer.output[3:]

				if new_packet.to_id == system.nat.id {
					// Sent for the target port
					system.nat.packet = new_packet
				} else {
					idle = false
					// Sent to another module
					sent_to_module := system.modules[new_packet.to_id]
					sent_to_module.add_packet(new_packet)
					system.modules[new_packet.to_id] = sent_to_module
				}
			}
		}

		system.modules[module_id] = module
	}

	return idle
}

func (system *System) run_until_packet_for_target() {
	var ids []int = make([]int, 0)
	for module_id, _ := range system.modules {
		ids = append(ids, module_id)
	}

	for system.nat.packet == INVALID_PACKET {
		system.run(ids)
	}
}

func (system *System) run_until_second_in_a_row_idle() {
	var ids []int = make([]int, 0)
	for module_id, _ := range system.modules {
		ids = append(ids, module_id)
	}

	var last_idle bool = false
	for true {
		// Actually running
		var idle bool = system.run(ids)

		if idle && last_idle {
			// Two in a row idle
			break

		} else if idle {
			// System is idle
			sent_to_module := system.modules[system.nat.send_to]
			sent_to_module.computer.input = append(sent_to_module.computer.input, system.nat.packet.x)
			sent_to_module.computer.input = append(sent_to_module.computer.input, system.nat.packet.y)
			system.modules[system.nat.send_to] = sent_to_module
		}

		last_idle = idle
	}
}

// ----------------------- System Struct End -----------------------

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

		// Variable Set
		var TARGET_PORT int = 255
		var TARGET_PORT_FOR_NAT int = 0
		var NUMBER_MODULES int = 50

		var mock_computer IntCodeComputer = IntCodeComputer{"booting", make([]int, 0), 0, values_converted, 0, 0, 0, make([]int, 0)}
		var system System = new_system(NUMBER_MODULES, mock_computer, TARGET_PORT, TARGET_PORT_FOR_NAT)

		// Part 1
		system.run_until_packet_for_target()
		packet := system.nat.packet
		fmt.Printf("The first packet for ' %d ' has a Y ' %d ' (part 1)\n", TARGET_PORT, packet.y)

		// Part 2
		system.run_until_second_in_a_row_idle()
		packet = system.nat.packet
		fmt.Printf("The second in a row idle packe in NAT has a Y ' %d ' (part 2)\n", packet.y)
	}
}
