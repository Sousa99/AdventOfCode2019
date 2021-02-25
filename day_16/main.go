package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

// ----------------------- System Struct Start -----------------------

type System struct {
	input        []int
	base_pattern []int
}

func (system *System) create_phase(phase int) []int {
	base_size := len(system.base_pattern)
	number_phases := len(system.input)

	var new_phase []int = make([]int, 0)
	current_pointer := 0

	for len(new_phase) <= number_phases+1 {
		for iteration := 0; iteration <= phase; iteration++ {
			element := system.base_pattern[current_pointer%base_size]
			new_phase = append(new_phase, element)
		}

		current_pointer = current_pointer + 1
	}

	return new_phase[1 : number_phases+1]
}

func (system *System) run_n_phases(number_times int) string {
	var current_input []int = make([]int, len(system.input))
	copy(current_input, system.input)
	number_phases := len(system.input)

	for time := 0; time < number_times; time++ {

		var new_input []int = make([]int, 0)
		for phase_number := 0; phase_number < number_phases; phase_number++ {

			var phase []int = system.create_phase(phase_number)
			var current_value int = 0
			for value_index, value := range current_input {
				current_value = current_value + value*phase[value_index]
			}

			current_value = int(math.Abs(float64(current_value))) % 10
			new_input = append(new_input, current_value)
		}

		current_input = new_input
	}

	// Transform slice into a single int
	var final_value string = ""
	for _, value := range current_input {
		final_value = final_value + strconv.Itoa(value)
	}

	return final_value
}

func (system *System) run_n_phases_truncated(number_times int, truncated int) string {
	var original_input_size int = len(system.input)
	var input_size int = original_input_size - truncated

	var current_input []int = make([]int, input_size)
	copy(current_input, system.input[truncated:])

	for time := 0; time < number_times; time++ {

		var new_input []int = make([]int, len(current_input))
		last_value := 0
		for index := input_size - 1; index >= 0; index-- {

			if original_input_size-index > original_input_size/2 {
				// If second half
				saved_value := current_input[index]
				last_value = int(math.Abs(float64(last_value+saved_value))) % 10
				new_input[index] = last_value

			} else {
				// If first half
				var phase []int = system.create_phase(original_input_size - index - 1)
				var current_value int = 0
				for value_index, value := range current_input {
					current_value = current_value + value*phase[value_index]
				}

				current_value = int(math.Abs(float64(current_value))) % 10
				new_input[index] = current_value
			}
		}

		current_input = new_input
	}

	// Transform slice into a single int
	var final_value string = ""
	for _, value := range current_input {
		final_value = final_value + strconv.Itoa(value)
	}

	return final_value
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
		inputs := make([]int, 0)
		for _, code := range line {
			value, _ := strconv.Atoi(string(code))
			inputs = append(inputs, value)
		}

		// Part 1
		base_pattern := []int{0, 1, 0, -1}
		var system System = System{inputs, base_pattern}

		var number_of_iterations int = 100
		var number_digits int = 8
		output := system.run_n_phases(number_of_iterations)[:number_digits]
		fmt.Printf("After ' %d ' iterations: ' %s ' (part 1)\n", number_of_iterations, output)

		// Part 2
		var digits_for_offset int = 7
		var repetitions int = 10000

		initial_offset := 0
		for index := 0; index < digits_for_offset; index++ {
			scale := digits_for_offset - 1 - index
			initial_offset = initial_offset + int(math.Pow10(scale))*inputs[index]
		}

		new_inputs_rep := make([]int, 0)
		for rep := 0; rep < repetitions; rep++ {
			new_inputs_rep = append(new_inputs_rep, inputs...)
		}

		base_pattern = []int{0, 1, 0, -1}
		var system_rep System = System{new_inputs_rep, base_pattern}

		output = system_rep.run_n_phases_truncated(number_of_iterations, initial_offset)[:number_digits]
		fmt.Printf("After ' %d ' iterations: ' %s ' (part 2)\n", number_of_iterations, output)
	}
}
