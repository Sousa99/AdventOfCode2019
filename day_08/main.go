package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

// ----------------------- Picture Struct Start -----------------------

type Pixel = int
type Line = []Pixel
type Layer = []Line

type Picture struct {
	width  int
	height int
	layers []Layer
}

func (picture *Picture) create_from_raw(pixels []Pixel) {
	var n_layers = len(pixels) / (picture.width * picture.height)

	for layer := 0; layer < n_layers; layer++ {

		var new_layer Layer = make([][]Pixel, 0)
		for line := 0; line < picture.height; line++ {

			var new_line Line = make([]Pixel, 0)
			for pixel_pos := 0; pixel_pos < picture.width; pixel_pos++ {

				var index int = layer*picture.height*picture.width + line*picture.width + pixel_pos
				var pixel Pixel = pixels[index]
				new_line = append(new_line, pixel)
			}

			new_layer = append(new_layer, new_line)
		}

		picture.layers = append(picture.layers, new_layer)
	}
}

func (picture *Picture) get_result_layer_least_zeros() (int, int) {
	var selected_layer int = -1
	var selected_number_zeros int = -1
	var selected_result int = -1

	for index_layer, layer := range picture.layers {
		current_count := map[int]int{0: 0, 1: 0, 2: 0}

		// Count different numbers
		for _, line := range layer {
			for _, pixel := range line {
				if pixel >= 0 && pixel <= 2 {
					current_count[pixel] = current_count[pixel] + 1
				}
			}
		}

		// Check if new minimum
		if current_count[0] < selected_number_zeros || selected_number_zeros == -1 {
			selected_layer = index_layer
			selected_number_zeros = current_count[0]
			selected_result = current_count[1] * current_count[2]
		}
	}

	return selected_layer + 1, selected_result
}

func (picture *Picture) develop_image() Layer {
	// Initialize_array
	var developed_layer Layer = picture.layers[0]

	for _, layer := range picture.layers {
		for index_line, line := range layer {
			for index_pixel, pixel := range line {

				// If transparent
				if developed_layer[index_line][index_pixel] == 2 {
					developed_layer[index_line][index_pixel] = pixel
				}
			}
		}
	}

	return developed_layer
}

// ----------------------- Picture Struct End -----------------------

func save_layer_as_image(file string, layer Layer) {
	width := len(layer[0])
	height := len(layer)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	colorCoding := map[int]*image.Uniform{
		0: image.NewUniform(color.Black),
		1: image.NewUniform(color.White),
		2: image.NewUniform(color.Transparent),
	}

	for line_index, line := range layer {
		for pixel_index, pixel := range line {
			img.Set(pixel_index, line_index, colorCoding[pixel])
		}
	}

	f, _ := os.Create(file)
	png.Encode(f, img)
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

		var picture Picture = Picture{25, 6, make([]Layer, 0)}
		var pixels []int = make([]Pixel, 0)

		// Get pixels individually
		var line_read string = scanner.Text()
		for _, digit_string := range line_read {
			digit, _ := strconv.Atoi(string(digit_string))
			pixels = append(pixels, digit)
		}

		// Verify picture (part 1)
		picture.create_from_raw(pixels)
		layer, result := picture.get_result_layer_least_zeros()
		fmt.Printf("The layer ' %d ' has a result of ' %d ' (part 1)\n", layer, result)

		// Develop picture (part 2)
		var developed Layer = picture.develop_image()
		save_layer_as_image("output.png", developed)
	}
}
