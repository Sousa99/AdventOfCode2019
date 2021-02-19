package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ----------------------- System Struct Start -----------------------

type System struct {
	orbits  map[string][]string
	orbited map[string][]string
}

func (system *System) add_orbit(center string, orbited string) {
	var center_already_added bool = false
	var orbited_already_added bool = false
	for mass, _ := range system.orbits {
		center_already_added = center_already_added || mass == center
		orbited_already_added = orbited_already_added || mass == orbited
	}

	// If center doesn't exist add it
	if !center_already_added {
		system.orbits[center] = make([]string, 0)
		system.orbited[center] = make([]string, 0)
	}
	// If orbited doesn't exist add it
	if !orbited_already_added {
		system.orbits[orbited] = make([]string, 0)
		system.orbited[orbited] = make([]string, 0)
	}

	// If A ) B then orbits[A] includes B
	system.orbits[center] = append(system.orbits[center], orbited)
	// If A ) B then orbits[B] includes A
	system.orbited[orbited] = append(system.orbits[orbited], center)
}

func (system *System) compute_orbit_depths(initial string) map[string]int {
	var to_explore []string = make([]string, 0)
	to_explore = append(to_explore, initial)
	var depths map[string]int = make(map[string]int)
	var current_depth int = 0

	for len(to_explore) != 0 {
		var new_to_explore []string = make([]string, 0)

		for _, mass := range to_explore {
			depths[mass] = current_depth
			new_to_explore = append(new_to_explore, system.orbits[mass]...)
		}

		to_explore = new_to_explore
		current_depth = current_depth + 1
	}

	return depths
}

func (system *System) compute_transfers(initial string, end string) []string {
	var transfers map[string][]string = make(map[string][]string)
	transfers[initial] = []string{initial}
	transfers[end] = []string{end}

	var final_path []string = nil
	for final_path == nil {

		var new_transfers map[string][]string = make(map[string][]string)
		for from_point, path := range transfers {

			// Dealing with connections
			var connections []string = make([]string, 0)
			connections = append(connections, system.orbits[from_point]...)
			connections = append(connections, system.orbited[from_point]...)

			for _, to_point := range connections {
				set_path, already_set_path := transfers[to_point]
				if already_set_path && set_path[0] == path[0] {
					// Closest path already discovered
					continue
				}

				var update_path []string = make([]string, 0)
				update_path = append(update_path, path...)
				update_path = append(update_path, to_point)

				if already_set_path && set_path[0] != path[0] {
					// Connection point discovered
					final_path = make([]string, 0)
					final_path = append(final_path, set_path...)
					for i := len(update_path) - 2; i >= 0; i-- {
						final_path = append(final_path, update_path[i])
					}

					break
				} else {
					new_transfers[to_point] = update_path
				}
			}
		}

		// Add new transfers
		for key, value := range new_transfers {
			transfers[key] = value
		}
	}

	return final_path
}

// ----------------------- System Struct End -----------------------

func compute_sum_depths(depths map[string]int) int {
	var sum_depth int = 0
	for _, depth := range depths {
		sum_depth = sum_depth + depth
	}

	return sum_depth
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var system System = System{make(map[string][]string), make(map[string][]string)}

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Iterate over split
		var split []string = strings.Split(scanner.Text(), ")")
		var center string = split[0]
		var orbited string = split[1]
		system.add_orbit(center, orbited)
	}

	// Part 1
	depths := system.compute_orbit_depths("COM")
	sum_depth := compute_sum_depths(depths)
	fmt.Printf("Number of direct / indirect orbits: ' %d ' (part 1)\n", sum_depth)

	// Part 2
	transfers := system.compute_transfers("YOU", "SAN")
	number_transfers := len(transfers) - 3
	fmt.Printf("Number of transfers: ' %d ' (part 2)\n", number_transfers)
}
