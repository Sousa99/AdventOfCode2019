package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ----------------------- Program Struct Start -----------------------

type Program struct {
	pointer int
	codes   []int
}

func (program *Program) append_code(code int) {
	program.codes = append(program.codes, code)
}

func (program *Program) run() {
	var halt bool = false
	for !halt {

		// Run iteration
		var current_code int = program.codes[program.pointer]
		switch current_code {
		// Halting
		case 99:
			halt = true

		// Addition
		case 1:
			var first_position int = program.codes[program.pointer+1]
			var first_value int = program.codes[first_position]
			var second_position int = program.codes[program.pointer+2]
			var second_value int = program.codes[second_position]

			var writing_position int = program.codes[program.pointer+3]
			program.codes[writing_position] = first_value + second_value

			//Advance pointer
			program.pointer = program.pointer + 4

		// Multiplication
		case 2:
			var first_position int = program.codes[program.pointer+1]
			var first_value int = program.codes[first_position]
			var second_position int = program.codes[program.pointer+2]
			var second_value int = program.codes[second_position]

			var writing_position int = program.codes[program.pointer+3]
			program.codes[writing_position] = first_value * second_value

			//Advance pointer
			program.pointer = program.pointer + 4

		// Should not happen
		default:
			fmt.Println("Code not recognized")
			os.Exit(1)
		}
	}
}

// ----------------------- Program Struct End -----------------------

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

		var program Program = Program{0, make([]int, 0)}
		// Iterate over split
		var split []string = strings.Split(scanner.Text(), ",")
		for _, code_string := range split {
			code, _ := strconv.Atoi(code_string)
			program.append_code(code)
		}

		// Run program (part 1)
		program.codes[1] = 12
		program.codes[2] = 2
		program.run()
		fmt.Println("Result: '", program.codes[0], "' (part 1)")
	}
}
