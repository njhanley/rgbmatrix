package types

import "image/color"

type None struct{}

type Pixel struct {
	X, Y int
	C    color.RGBA
}
