package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const DEBUG bool = true

func find_valid_passwords(start int, end int, mode int) []int {
	var valid_passwords []int = make([]int, 0)
	var current_value int = start

	for current_value <= end {
		if DEBUG {
			fmt.Printf("\033[2K\rValidation Process: %d%%", (current_value-start)*100/(end-start))
		}

		var last_digit int = -1
		var digit_count int = 0

		var double_criteria bool = false
		var repetitions []int = make([]int, 0)

		var temp int = current_value
		for temp != 0 {

			var digit int = temp % 10
			// Invalid scenario
			if last_digit != -1 && digit > last_digit {
				break
			}

			// Repetition of digits handling
			if digit == last_digit {
				double_criteria = true
				repetitions[len(repetitions)-1] = repetitions[len(repetitions)-1] + 1
			} else {
				repetitions = append(repetitions, 1)
			}

			digit_count = digit_count + 1
			last_digit = digit
			temp = temp / 10
		}

		// Some check failed
		if (mode == 0 && !double_criteria) || (mode == 1 && !any(repetitions, 2)) {
			current_value = current_value + 1
			continue
		} else if temp != 0 {
			current_value = current_value + int(math.Pow10(digit_count-1))
			continue
		}

		// Add valid password
		valid_passwords = append(valid_passwords, current_value)
		current_value = current_value + 1
	}

	if DEBUG {
		fmt.Printf("\033[2K\rValidation Process: %d%%\n", (current_value-start)*100/(end-start))
	}

	return valid_passwords
}

func any(list []int, value int) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
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
		var split []string = strings.Split(scanner.Text(), "-")
		start, _ := strconv.Atoi(split[0])
		end, _ := strconv.Atoi(split[1])

		// Part 1
		var valid_passwords []int = find_valid_passwords(start, end, 0)
		fmt.Printf("Number of valid passwords: %d (part 1)\n", len(valid_passwords))

		// Part 1
		valid_passwords = find_valid_passwords(start, end, 1)
		fmt.Printf("Number of valid passwords: %d (part 2)\n", len(valid_passwords))
	}
}
