package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var win_width, win_height = 800, 600
var square_size = 25
var n_rows, n_cols = win_height / square_size, win_width / square_size
var bg_color = color.RGBA{102, 205, 170, 0xff}   // MediumAquamarine
var snake_color = color.RGBA{0, 128, 128, 0xff}  // Teal
var food_color = color.RGBA{255, 0, 0, 0xff}     // Red
var fill_color = color.RGBA{255, 255, 255, 0xff} // While

var display_start = true

// Set the game speed (the lower number, the bigger speed.)
var speed_delay = 100

// Direction in which the snake is heading
var dx, dy = 1, 0

// Segments of the snake body

var best_score = 0
var score = 0
var crashed = false

var background = get_background(bg_color, win_width, win_height)

type Snake struct {
	size     image.Rectangle
	Position image.Point
	Status   bool
	Points   int
}

// this only works for iTerm2!
func printImage(img image.Image) {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Printf("\x1b[2;0H\x1b]1337;File=inline=1:%s\a", imgBase64Str)
}

func get_background(color color.RGBA, width, height int) *image.RGBA {
	up_left := image.Point{0, 0}
	low_right := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{up_left, low_right})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color)
		}
	}
	return img
}

func draw_arena(dst *image.RGBA, color color.RGBA, n_rows, n_cols, square_size int) {
	up_left := image.Point{0, 0}
	low_right := image.Point{square_size * n_cols, square_size * n_rows}
	sr := image.Rectangle{up_left, low_right}

	src := image.NewRGBA(sr)
	// Set color for each pixel.
	for x := 0; x < square_size*n_cols; x += square_size {
		for y := 0; y < square_size*n_rows; y++ {
			src.Set(x, y, color)
		}
	}

	for y := 0; y < square_size*n_rows; y++ {
		src.Set(square_size*n_cols-1, y, color)
	}

	// Set color for each pixel.
	for y := 0; y < square_size*n_rows; y += square_size {
		for x := 0; x < square_size*n_cols; x++ {
			src.Set(x, y, color)
		}
	}
	for x := 0; x < square_size*n_cols; x++ {
		src.Set(x, square_size*n_rows-1, color)
	}

	dp := image.Point{0, 0}
	r := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(dst, r, src, sr.Min, draw.Src)
}

func draw_food(dst *image.RGBA, color color.RGBA, px, py, square_size int) {
	up_left := image.Point{0, 0}
	low_right := image.Point{square_size, square_size}
	sr := image.Rectangle{up_left, low_right}

	src := image.NewRGBA(sr)
	// Set color for each pixel.
	for x := 0; x < square_size; x++ {
		for y := 0; y < square_size; y++ {
			src.Set(x, y, color)
		}
	}

	dp := image.Point{px * square_size, py * square_size}
	r := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(dst, r, src, sr.Min, draw.Src)
}

func draw_snake(dst *image.RGBA, color color.RGBA, px, py, square_size int) {
	up_left := image.Point{0, 0}
	low_right := image.Point{square_size, square_size}
	sr := image.Rectangle{up_left, low_right}

	src := image.NewRGBA(sr)
	// Set color for each pixel.
	for x := 0; x < square_size; x++ {
		for y := 0; y < square_size; y++ {
			src.Set(x, y, color)
		}
	}

	dp := image.Point{px * square_size, py * square_size}
	r := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(dst, r, src, sr.Min, draw.Src)
}

func draw_stats(dst *image.RGBA, color color.RGBA, width, height int) {
	addLabel(dst, color, 10, height-10, "Use arrow keys to change snake direction!")
	addLabel(dst, color, width/2+10, height-10, fmt.Sprintf("Best Score: %d", best_score))
	addLabel(dst, color, width-100, height-10, fmt.Sprintf("Score: %d", score))

}
func addLabel(img *image.RGBA, color color.RGBA, x, y int, label string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func main() {
	dst := get_background(fill_color, win_width, win_height)
	draw_arena(dst, bg_color, n_rows, n_cols, square_size)
	draw_food(dst, food_color, 10, 10, square_size)
	draw_stats(dst, fill_color, win_width, win_height)
	printImage(dst)
}
