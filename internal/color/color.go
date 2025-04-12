package color

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// HSV to RGB conversion
func HsvToRgb(h, s, v float32) rl.Color {
	if s == 0 {
		// Gray
		value := uint8(v * 255)
		return rl.Color{R: value, G: value, B: value, A: 255}
	}

	h = float32(math.Mod(float64(h), 360))
	h /= 60 // sector 0 to 5
	i := int(h)
	f := h - float32(i) // fractional part of h

	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float32
	switch i {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}

	return rl.Color{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}
