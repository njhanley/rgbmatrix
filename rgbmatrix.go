// Package rgbmatrix is a Go binding for https://github.com/hzeller/rpi-rgb-led-matrix
package rgbmatrix

// #cgo CFLAGS: -I${SRCDIR}/rpi-rgb-led-matrix/include
// #cgo LDFLAGS: -L${SRCDIR}/rpi-rgb-led-matrix/lib -lrgbmatrix -lstdc++ -lm
// #include <led-matrix-c.h>
import "C"

import (
	"errors"
	"image"
	"image/color"
)

// Config is a Matrix configuration.
type Config struct {
	Rows              int // Number of rows on a single panel
	Columns           int // Number of columns on a single panel
	ChainLength       int // Number of daisy-chained panels
	Parallel          int // Number of parallel chains
	Brightness        int // Brightness percentage
	PWMBits           int
	PWMLSBNanoseconds int
}

func (c Config) toRGBLedMatrixOptions() *C.struct_RGBLedMatrixOptions {
	return &C.struct_RGBLedMatrixOptions{
		rows:                C.int(c.Rows),
		cols:                C.int(c.Columns),
		chain_length:        C.int(c.ChainLength),
		parallel:            C.int(c.Parallel),
		brightness:          C.int(c.Brightness),
		pwm_bits:            C.int(c.PWMBits),
		pwm_lsb_nanoseconds: C.int(c.PWMLSBNanoseconds),
	}
}

var DefaultConfig = Config{
	Rows:              32,
	Columns:           32,
	ChainLength:       1,
	Parallel:          1,
	Brightness:        100,
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
}

// Matrix represents an LED matrix.
// It is not safe to call methods on a Matrix concurrently.
type Matrix struct {
	canvas *C.struct_LedCanvas
	matrix *C.struct_RGBLedMatrix
	width  int
	height int
}

// New initializes and returns a Matrix.
// The rpi-rgb-led-matrix library may write to stderr.
func New(cfg Config) (*Matrix, error) {
	matrix := C.led_matrix_create_from_options(cfg.toRGBLedMatrixOptions(), nil, nil)
	if matrix == nil {
		return nil, errors.New("failed to initialize matrix")
	}

	return &Matrix{
		canvas: C.led_matrix_create_offscreen_canvas(matrix),
		matrix: matrix,
		width:  cfg.Columns * cfg.ChainLength,
		height: cfg.Rows * cfg.Parallel,
	}, nil
}

// Close frees allocated resources and resets the hardware.
func (ma *Matrix) Close() {
	C.led_matrix_delete(ma.matrix)
}

// Size of the canvas in pixels.
func (ma *Matrix) Size() image.Point {
	return image.Point{ma.width, ma.height}
}

// SwapOnVSync swaps the front and back canvases.
// Call this method after modifying the back canvas to display your changes.
func (ma *Matrix) SwapOnVSync() {
	ma.canvas = C.led_matrix_swap_on_vsync(ma.matrix, ma.canvas)
}

// Clear the back canvas.
func (ma *Matrix) Clear() {
	C.led_canvas_clear(ma.canvas)
}

func toRGB(c color.Color) (r, g, b C.uint8_t) {
	cl := color.RGBAModel.Convert(c).(color.RGBA)
	return C.uint8_t(cl.R), C.uint8_t(cl.G), C.uint8_t(cl.B)
}

// Fill the back canvas with the given color.
func (ma *Matrix) Fill(c color.Color) {
	r, g, b := toRGB(c)
	C.led_canvas_fill(ma.canvas, r, g, b)
}

// Set the color of a pixel on the back canvas.
func (ma *Matrix) Set(x, y int, c color.Color) {
	r, g, b := toRGB(c)
	C.led_canvas_set_pixel(ma.canvas, C.int(x), C.int(y), r, g, b)
}

// Draw an image on the back canvas.
func (ma *Matrix) Draw(m image.Image) {
	for y := 0; y < ma.height; y++ {
		for x := 0; x < ma.width; x++ {
			ma.Set(x, y, m.At(x, y))
		}
	}
}
