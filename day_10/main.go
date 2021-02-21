package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// ----------------------- Asteroid Map Struct Start -----------------------

type Point = string
type AsteroidMap struct {
	max_x  int
	max_y  int
	points map[int]map[int]Point
}

func (asteroid_map *AsteroidMap) add_point(x int, y int, point Point) {
	// Create line if it wasn't yet created
	_, row_exists := asteroid_map.points[y]
	if !row_exists {
		asteroid_map.points[y] = make(map[int]string)
	}

	// Add point
	asteroid_map.points[y][x] = point

	// Update max_x and max_y
	if x > asteroid_map.max_x {
		asteroid_map.max_x = x
	}
	if y > asteroid_map.max_y {
		asteroid_map.max_y = y
	}
}

func (asteroid_map *AsteroidMap) get_visibility_count_map() map[int]map[int]int {
	var visibility map[int]map[int]int = make(map[int]map[int]int)

	for y_index, line := range asteroid_map.points {
		for x_index, point := range line {
			// If it isn't asteroid don't even consider it
			if point != "#" {
				continue
			}

			sub_asteroid_map := asteroid_map.create_map_of_visbility(x_index, y_index)
			var count_visible_asteroids int = 0
			for _, line := range sub_asteroid_map {
				for _, point := range line {
					if point == "#" {
						count_visible_asteroids = count_visible_asteroids + 1
					}
				}
			}

			// Verify if line initialized
			_, line_exists := visibility[y_index]
			if !line_exists {
				visibility[y_index] = make(map[int]int)
			}

			// Add count
			visibility[y_index][x_index] = count_visible_asteroids
		}
	}

	return visibility
}

func (asteroid_map *AsteroidMap) create_map_of_visbility(x int, y int) map[int]map[int]string {
	var map_of_visbility map[int]map[int]string = make(map[int]map[int]string)
	var distance_asteroid_map map[int][3]int = make(map[int][3]int)

	for y_index := 0; y_index <= asteroid_map.max_y; y_index++ {
		map_of_visbility[y_index] = make(map[int]string)

		for x_index := 0; x_index <= asteroid_map.max_x; x_index++ {
			map_of_visbility[y_index][x_index] = asteroid_map.points[y_index][x_index]

			if asteroid_map.points[y_index][x_index] == "#" && (x != x_index || y != y_index) {
				var new_x float64 = float64(x_index - x)
				var new_y float64 = float64(y_index - y)
				var angle int = int(math.Round((math.Atan2(new_y, new_x) * (180.0 / math.Pi)) * 100))

				var distance int = int(math.Abs(new_y)) + int(math.Abs(new_x))

				// Check existence of angle
				value_saved, angle_exists := distance_asteroid_map[angle]
				if !angle_exists {
					distance_asteroid_map[angle] = [3]int{x_index, y_index, distance}
				} else {
					if value_saved[2] <= distance {
						// Existing value remains
						map_of_visbility[y_index][x_index] = "X"
					} else {
						// New value is closer
						map_of_visbility[value_saved[1]][value_saved[0]] = "X"
						distance_asteroid_map[angle] = [3]int{x_index, y_index, distance}
					}
				}

			}
		}
	}

	map_of_visbility[y][x] = "O"
	return map_of_visbility
}

func (asteroid_map *AsteroidMap) eliminate_asteroids(laser_position_x int, laser_position_y int) [][2]int {
	asteroid_map.points[laser_position_y][laser_position_x] = "O"
	var elimination_list [][2]int = make([][2]int, 0)

	var remaining_asteroids bool = true
	for remaining_asteroids {

		// Get asteroids to remove
		var distance_asteroid_map map[int][3]int = make(map[int][3]int)
		var angles_present []int = make([]int, 0)
		remaining_asteroids = false

		for y_index := 0; y_index <= asteroid_map.max_y; y_index++ {
			for x_index := 0; x_index <= asteroid_map.max_x; x_index++ {
				if asteroid_map.points[y_index][x_index] != "#" {
					continue
				}

				remaining_asteroids = true

				var new_x float64 = float64(x_index - laser_position_x)
				var new_y float64 = float64(y_index - laser_position_y)
				var angle_degrees float64 = math.Atan2(new_y, new_x) * (180.0 / math.Pi)
				var angle_degrees_fixed float64 = math.Mod(angle_degrees+180+270, 360)
				var angle int = int(math.Round(angle_degrees_fixed * 100))

				var distance int = int(math.Abs(new_y)) + int(math.Abs(new_x))

				// Check existence of angle
				value_saved, angle_exists := distance_asteroid_map[angle]
				if !angle_exists {
					distance_asteroid_map[angle] = [3]int{x_index, y_index, distance}
					angles_present = append(angles_present, angle)
				} else if value_saved[2] > distance {
					// New value is closer
					distance_asteroid_map[angle] = [3]int{x_index, y_index, distance}
				}
			}
		}

		// Delete asteroids and insert into elimination list
		sort.Ints(angles_present)
		for _, angle := range angles_present {
			info, _ := distance_asteroid_map[angle]
			asteroid_x, asteroid_y := info[0], info[1]

			// Update map
			asteroid_map.points[asteroid_y][asteroid_x] = "."
			// Add to removal list
			elimination_list = append(elimination_list, [2]int{asteroid_x, asteroid_y})
		}
	}

	return elimination_list
}

// ----------------------- Asteroid Map Struct End -----------------------

func compute_max(visibility map[int]map[int]int) (int, int, int) {
	var max_count, max_x, max_y int = -1, -1, -1

	for y, line := range visibility {
		for x, asteroids_count := range line {
			if asteroids_count > max_count {
				max_count, max_x, max_y = asteroids_count, x, y
			}
		}
	}

	return max_x, max_y, max_count
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var asteroid_map AsteroidMap = AsteroidMap{0, 0, make(map[int]map[int]string)}

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	var y int = 0
	for scanner.Scan() {

		// Iterate over line
		var line string = scanner.Text()
		for x, code := range line {
			asteroid_map.add_point(x, y, string(code))
		}

		y = y + 1
	}

	// Part 1
	visibility_map := asteroid_map.get_visibility_count_map()
	max_x, max_y, max_count := compute_max(visibility_map)
	fmt.Printf("( %d, %d) with a count of ' %d ' asteroids (part 1)\n", max_x, max_y, max_count)

	// Part 2
	nth_eliminated := 200
	elimination_list := asteroid_map.eliminate_asteroids(max_x, max_y)
	eliminated_x, eliminated_y := elimination_list[nth_eliminated-1][0], elimination_list[nth_eliminated-1][1]
	result := eliminated_x*100 + eliminated_y
	fmt.Printf("( %d, %d) is the %dth eliminated with a result of ' %d ' (part 2)\n", eliminated_x, eliminated_y, nth_eliminated, result)
}
