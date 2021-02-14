package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func computeFuel(mass int) int {
	var div float64 = float64(mass) / 3.0
	var round_down int = int(math.Floor(div))
	var subtract int = round_down - 2
	return subtract
}

func computeFuelNeededForFuel(fuel int) int {
	var extra_needed int = 0
	var sub_fuel int = computeFuel(fuel)

	for sub_fuel > 0 {
		extra_needed = sub_fuel + extra_needed
		sub_fuel = computeFuel(sub_fuel)
	}

	return extra_needed
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var total_fuel int = 0
	var total_fuel_considering_fuel int = 0

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mass, _ := strconv.Atoi(scanner.Text())

		// Fuel needed for module mass
		var sub_fuel int = computeFuel(mass)
		// Fuel needed for fuel added for the module mass
		var fuel_for_fuel int = computeFuelNeededForFuel(sub_fuel)

		total_fuel = total_fuel + sub_fuel
		total_fuel_considering_fuel = total_fuel_considering_fuel + sub_fuel + fuel_for_fuel
	}

	fmt.Println("Total fuel needed '", total_fuel, "' (part 1)")
	fmt.Println("Total fuel needed '", total_fuel_considering_fuel, "' (part 2)")
}
