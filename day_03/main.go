package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ----------------------- Directions Struct Start -----------------------
type Direction struct {
	column_variation int
	row_variation    int
	point_type       string
}

var Directions = map[string]Direction{
	"R": Direction{1, 0, "horizontal"},
	"L": Direction{-1, 0, "horizontal"},
	"U": Direction{0, 1, "vertical"},
	"D": Direction{0, -1, "vertical"},
}

// ----------------------- Directions Struct End -----------------------

// ----------------------- Indication Struct Start -----------------------

type Indication struct {
	direction Direction
	length    int
}

func (indication *Indication) get_column_variation() int { return indication.direction.column_variation }
func (indication *Indication) get_row_variation() int    { return indication.direction.row_variation }
func (indication *Indication) get_point_type() string    { return indication.direction.point_type }
func (indication *Indication) get_length() int           { return indication.length }

// ----------------------- Indication Struct End -----------------------

// ----------------------- Panel Point Struct Start -----------------------

var PointTypesChars = map[string]string{
	"origin":           string("\033[32m") + "o" + string("\033[0m"),
	"horizontal":       "-",
	"vertical":         "|",
	"bend":             "+",
	"intersection_eq":  "+",
	"intersection_dif": string("\033[31m") + "X" + string("\033[0m"),
}

type PanelPoint struct {
	cable_id   int
	point_type string
	distance   int
}

func convert_panel_points(points []PanelPoint) string {
	switch len(points) {
	case 0:
		return "."
	case 1:
		var point PanelPoint = points[0]
		return PointTypesChars[point.point_type]
	default:
		var cable_ids []int = make([]int, 0)
		for _, point := range points {

			var already_added bool = false
			for _, id := range cable_ids {
				if point.cable_id == id {
					already_added = true
				}
			}

			if !already_added {
				cable_ids = append(cable_ids, point.cable_id)
			}
		}

		if len(cable_ids) == 1 {
			return PointTypesChars["intersection_eq"]
		} else {
			return PointTypesChars["intersection_dif"]
		}
	}
}

func point_is_valid_intersection(points []PanelPoint) bool {
	if len(points) > 1 {
		var cable_ids []int = make([]int, 0)
		for _, point := range points {

			var already_added bool = false
			for _, id := range cable_ids {
				if point.cable_id == id {
					already_added = true
				}
			}

			if !already_added {
				cable_ids = append(cable_ids, point.cable_id)
			}
		}

		if len(cable_ids) == 1 {
			return false
		} else {
			return true
		}
	}

	return false
}

func get_min_distance(points []PanelPoint) int {
	var first_cable_points []PanelPoint = make([]PanelPoint, 0)

	for _, point := range points {

		var already_added_index int = -1
		// Check if cable_id already added
		for index, added_point := range first_cable_points {
			if added_point.cable_id == point.cable_id {
				already_added_index = index
			}
		}

		if already_added_index != -1 && point.distance < first_cable_points[already_added_index].distance {
			first_cable_points[already_added_index] = point
		} else if already_added_index == -1 {
			first_cable_points = append(first_cable_points, point)
		}
	}

	// Compute sum of distance of cables
	var distance int = 0
	for _, point := range first_cable_points {
		distance = distance + point.distance
	}

	return distance
}

// ----------------------- Panel Point Struct End -----------------------

// ----------------------- Panel Struct Start -----------------------

type Panel struct {
	min_row    int
	min_column int
	max_row    int
	max_column int

	cable_ids []int
	position  map[int]map[int][]PanelPoint
}

func (panel *Panel) initialize_panel() {
	for row := panel.min_row; row <= panel.max_row; row++ {
		panel.position[row] = make(map[int][]PanelPoint)

		for column := panel.min_column; column <= panel.max_column; column++ {
			panel.position[row][column] = make([]PanelPoint, 0)
		}
	}

	panel.position[0][0] = append(panel.position[0][0], PanelPoint{-1, "origin", 0})
}

func (panel *Panel) validate_position_in_grid(row int, column int) {
	// Adding missing rows
	if row < panel.min_row {
		panel.min_row = row
	} else if row > panel.max_row {
		panel.max_row = row
	}
	// Adding missing columns
	if column < panel.min_column {
		panel.min_column = column
	} else if column > panel.max_column {
		panel.max_column = column
	}

	// Check if row exists
	if _, ok := panel.position[row]; !ok {
		// If row doesn't exists
		panel.position[row] = make(map[int][]PanelPoint)
	}
	// Check if column exists
	if _, ok := panel.position[row][column]; !ok {
		// If row doesn't exists
		panel.position[row][column] = make([]PanelPoint, 0)
	}
}

func (panel *Panel) print_panel() {
	for row := panel.max_row; row >= panel.min_row; row-- {
		for column := panel.min_column; column <= panel.max_column; column++ {
			fmt.Print(convert_panel_points(panel.position[row][column]) + " ")
		}
		fmt.Println()
	}
}

func (panel *Panel) add_cable(cable_id int, cable []Indication) {
	var length_indications int = len(cable)
	var current_row, current_column int = 0, 0
	var distance int = 0

	panel.cable_ids = append(panel.cable_ids, cable_id)

	for indication_index, indication := range cable {

		var row_variation int = indication.get_row_variation()
		var column_variation int = indication.get_column_variation()
		var indication_point_type string = indication.get_point_type()
		var length int = indication.get_length()

		for i := 0; i < length; i++ {
			current_row = current_row + row_variation
			current_column = current_column + column_variation
			distance = distance + 1

			panel.validate_position_in_grid(current_row, current_column)

			var point_type string = indication_point_type
			if i == length-1 && indication_index != length_indications-1 {
				point_type = "bend"
			}

			var panelPoint PanelPoint = PanelPoint{cable_id, point_type, distance}
			panel.position[current_row][current_column] = append(panel.position[current_row][current_column], panelPoint)
		}

		fmt.Printf("\033[2K\rIndication Processment: %d%% completed", (indication_index+1)*100/length_indications)
	}

	fmt.Println()
}

func (panel *Panel) get_closest_intersection() (int, int, int) {
	var min_row, min_column int = 0, 0
	var min_distance int = -1

	for row := panel.max_row; row >= panel.min_row; row-- {
		for column := panel.min_column; column <= panel.max_column; column++ {
			if !point_is_valid_intersection(panel.position[row][column]) {
				continue
			}

			var distance int = int(math.Abs(float64(row))) + int(math.Abs(float64(column)))
			if min_distance == -1 || distance < min_distance {
				min_row, min_column = row, column
				min_distance = distance
			}
		}
	}
	return min_row, min_column, min_distance
}

func (panel *Panel) position_has_cable_through(row int, column int, cable_id int) bool {
	var points []PanelPoint = panel.position[row][column]
	for _, point := range points {
		if point.cable_id == cable_id {
			return true
		}
	}

	return false
}

func (panel *Panel) get_minimizing_delay() (int, int, int) {
	var min_row, min_column int = 0, 0
	var min_distance int = -1

	for row := panel.max_row; row >= panel.min_row; row-- {
		for column := panel.min_column; column <= panel.max_column; column++ {
			if point_is_valid_intersection(panel.position[row][column]) {
				var distance int = get_min_distance(panel.position[row][column])

				if min_distance == -1 || distance < min_distance {
					min_row, min_column = row, column
					min_distance = distance
				}
			}
		}
	}

	return min_row, min_column, min_distance
}

// ----------------------- Panel Struct End -----------------------

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var panel Panel = Panel{-1, -1, 1, 1, make([]int, 0), make(map[int]map[int][]PanelPoint)}
	panel.initialize_panel()
	var cable_id int = 0

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var cable []Indication = make([]Indication, 0)

		// Iterate over split
		var split []string = strings.Split(scanner.Text(), ",")
		for _, indication_string := range split {
			direction_code := string(indication_string[0])
			length, _ := strconv.Atoi(indication_string[1:])

			var direction Direction = Directions[direction_code]
			var indication Indication = Indication{direction, length}
			cable = append(cable, indication)
		}

		panel.add_cable(cable_id, cable)
		cable_id = cable_id + 1
	}

	//panel.print_panel()

	// Part 1
	var column, row, distance int = panel.get_closest_intersection()
	fmt.Println("Closest intersection (", column, ",", row, ") with distance of '", distance, "' (part 1)")

	// Part 2
	column, row, distance = panel.get_minimizing_delay()
	fmt.Println("Closest intersection (", column, ",", row, ") with distance of '", distance, "' (part 2)")
}
