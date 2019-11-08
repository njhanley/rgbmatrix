// Package rgbmatrix is a Go binding for https://github.com/hzeller/rpi-rgb-led-matrix
package rgbmatrix

// #cgo CFLAGS: -I${SRCDIR}/rpi-rgb-led-matrix/include
// #cgo LDFLAGS: -L${SRCDIR}/rpi-rgb-led-matrix/lib -lrgbmatrix -lstdc++ -lm
// #include <led-matrix-c.h>
import "C"

import (
	"errors"
	"image"
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
	*image.RGBA
	canvas *C.struct_LedCanvas
	matrix *C.struct_RGBLedMatrix
}

// New initializes and returns a Matrix.
// The rpi-rgb-led-matrix library may write to stderr.
func New(cfg Config) (*Matrix, error) {
	matrix := C.led_matrix_create_from_options(cfg.toRGBLedMatrixOptions(), nil, nil)
	if matrix == nil {
		return nil, errors.New("failed to initialize matrix")
	}
	return &Matrix{
		RGBA:   image.NewRGBA(image.Rect(0, 0, cfg.Columns*cfg.ChainLength, cfg.Rows*cfg.Parallel)),
		canvas: C.led_matrix_create_offscreen_canvas(matrix),
		matrix: matrix,
	}, nil
}

// Close frees allocated resources and resets the hardware.
func (ma *Matrix) Close() {
	C.led_matrix_delete(ma.matrix)
}

func (ma *Matrix) Swap() {
	for y := ma.Rect.Min.Y; y < ma.Rect.Max.Y; y++ {
		for x := ma.Rect.Min.X; x < ma.Rect.Max.X; x++ {
			rgba := ma.RGBAAt(x, y)
			C.led_canvas_set_pixel(
				ma.canvas,
				C.int(x),
				C.int(y),
				C.uint8_t(rgba.R),
				C.uint8_t(rgba.G),
				C.uint8_t(rgba.B),
			)
		}
	}
	ma.canvas = C.led_matrix_swap_on_vsync(ma.matrix, ma.canvas)
}
