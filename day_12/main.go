package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ----------------------- SpaceSystem Struct Start -----------------------

type Vector3D struct {
	x int
	y int
	z int
}

type Mass struct {
	position Vector3D
	velocity Vector3D
}

func (mass *Mass) change_velocity_attraction(other_position Vector3D) {
	// Update component X
	if mass.position.x != other_position.x {
		aux_x := other_position.x - mass.position.x
		mass.velocity.x += aux_x / int(math.Abs(float64(aux_x)))
	}
	// Update component Y
	if mass.position.y != other_position.y {
		aux_y := other_position.y - mass.position.y
		mass.velocity.y += aux_y / int(math.Abs(float64(aux_y)))
	}
	// Update component X
	if mass.position.z != other_position.z {
		aux_z := other_position.z - mass.position.z
		mass.velocity.z += aux_z / int(math.Abs(float64(aux_z)))
	}
}

func (mass *Mass) get_energy() (int, int) {
	var potential float64 = math.Abs(float64(mass.position.x)) + math.Abs(float64(mass.position.y)) + math.Abs(float64(mass.position.z))
	var kinetic float64 = math.Abs(float64(mass.velocity.x)) + math.Abs(float64(mass.velocity.y)) + math.Abs(float64(mass.velocity.z))
	return int(potential), int(kinetic)
}

type SpaceSystem struct {
	time            int
	masses          []Mass
	original_masses []Mass
}

func (system *SpaceSystem) add_mass(position Vector3D) {
	var new_mass Mass = Mass{position, Vector3D{0, 0, 0}}
	system.masses = append(system.masses, new_mass)
	system.original_masses = append(system.masses, new_mass)
}

func (system *SpaceSystem) run_timestep() {
	// Update velocities
	for index, mass := range system.masses {
		for _, other_mass := range system.masses {
			mass.change_velocity_attraction(other_mass.position)
		}
		system.masses[index] = mass
	}

	// Update position based on velocities
	for index, mass := range system.masses {
		mass.position.x = mass.position.x + mass.velocity.x
		mass.position.y = mass.position.y + mass.velocity.y
		mass.position.z = mass.position.z + mass.velocity.z
		system.masses[index] = mass
	}

	system.time = system.time + 1
}

func (system *SpaceSystem) run_t_steps(t int) {
	for current_timestep := 0; current_timestep < t; current_timestep++ {
		system.run_timestep()
	}
}

func (system *SpaceSystem) get_energy() int {
	var total_energy int = 0
	for _, mass := range system.masses {
		potential, kinetic := mass.get_energy()
		total_energy = total_energy + potential*kinetic
	}

	return total_energy
}

func (system *SpaceSystem) like_original(set_reps Vector3D) Vector3D {
	var rep_x, rep_y, rep_z bool = true, true, true
	var number_masses int = len(system.masses)
	for index := 0; index < number_masses; index++ {
		rep_x = rep_x && system.masses[index].position.x == system.original_masses[index].position.x
		rep_x = rep_x && system.masses[index].velocity.x == system.original_masses[index].velocity.x

		rep_y = rep_y && system.masses[index].position.y == system.original_masses[index].position.y
		rep_y = rep_y && system.masses[index].velocity.y == system.original_masses[index].velocity.y

		rep_z = rep_z && system.masses[index].position.z == system.original_masses[index].position.z
		rep_z = rep_z && system.masses[index].velocity.z == system.original_masses[index].velocity.z
	}

	if set_reps.x == -1 && rep_x {
		set_reps.x = system.time
	}
	if set_reps.y == -1 && rep_y {
		set_reps.y = system.time
	}
	if set_reps.z == -1 && rep_z {
		set_reps.z = system.time
	}

	return set_reps
}

func (system *SpaceSystem) steps_till_rep() int {
	var steps_by_axis Vector3D = Vector3D{-1, -1, -1}
	var like_original bool = false

	for !like_original {
		system.run_timestep()
		steps_by_axis = system.like_original(steps_by_axis)

		like_original = steps_by_axis.x != -1 && steps_by_axis.y != -1 && steps_by_axis.z != -1
	}

	return least_common_multiple(steps_by_axis.x, steps_by_axis.y, steps_by_axis.z)
}

func (system *SpaceSystem) print_state() {
	for _, mass := range system.masses {
		var position Vector3D = mass.position
		var velocity Vector3D = mass.velocity

		fmt.Printf("pos=<x=%d, y=%d, z=%d>, ", position.x, position.y, position.z)
		fmt.Printf("vel=<x=%d, y=%d, z=%d>\n", velocity.x, velocity.y, velocity.z)
	}
}

// ----------------------- SpaceSystem Struct End -----------------------

func greatest_common_divider(a int, b int) int {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}

	return a
}

func least_common_multiple(a int, b int, integers ...int) int {
	result := a * b / greatest_common_divider(a, b)

	for i := 0; i < len(integers); i++ {
		result = least_common_multiple(result, integers[i])
	}

	return result
}

func main() {

	// ----------------- SETUP INPUT TXT -----------------
	// Trying to open file
	file, _ := os.Open("input.txt")
	// Defer closing of file
	defer file.Close()
	// ----------------- FINISHED INPUT TXT -----------------

	var spaceSystem SpaceSystem = SpaceSystem{0, make([]Mass, 0), make([]Mass, 0)}
	var spaceSystem_rep SpaceSystem = SpaceSystem{0, make([]Mass, 0), make([]Mass, 0)}

	// Create scanner over file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var line string = scanner.Text()
		line = strings.ReplaceAll(line, "<x=", "")
		line = strings.ReplaceAll(line, " y=", "")
		line = strings.ReplaceAll(line, " z=", "")
		line = strings.ReplaceAll(line, ">", "")

		var split []string = strings.Split(line, ",")
		var split_numbers []int = make([]int, 0)
		for _, elem := range split {
			number, _ := strconv.Atoi(elem)
			split_numbers = append(split_numbers, number)
		}

		var x, y, z int = split_numbers[0], split_numbers[1], split_numbers[2]
		var position Vector3D = Vector3D{x, y, z}
		spaceSystem.add_mass(position)
		spaceSystem_rep.add_mass(position)
	}

	// Part 1
	var number_steps int = 1000
	spaceSystem.run_t_steps(number_steps)
	total_energy := spaceSystem.get_energy()
	fmt.Printf("Total energy after ' %d ' steps: ' %d ' (part 1)\n", number_steps, total_energy)

	// Part 2
	steps_taken := spaceSystem_rep.steps_till_rep()
	fmt.Printf("It takes ' %d ' steps until the universe repeats itself (part 2)\n", steps_taken)
}
