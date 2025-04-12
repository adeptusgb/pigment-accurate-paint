package gui

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	Rect    rl.Rectangle
	Color   rl.Color
	Hover   bool
	Pressed bool
}

func NewButton(x, y, width, height float32, color rl.Color) Button {
	return Button{
		Rect:  rl.Rectangle{X: x, Y: y, Width: width, Height: height},
		Color: color,
	}
}

func (b *Button) Update(mousePos rl.Vector2) bool {
	b.Hover = rl.CheckCollisionPointRec(mousePos, b.Rect)
	b.Pressed = b.Hover && rl.IsMouseButtonPressed(rl.MouseLeftButton)
	return b.Pressed
}

func (b *Button) Draw() {
	color := b.Color
	if b.Hover {
		// Lighten color when hovered
		color.R = uint8(math.Min(float64(color.R)+20, 255))
		color.G = uint8(math.Min(float64(color.G)+20, 255))
		color.B = uint8(math.Min(float64(color.B)+20, 255))
	}
	rl.DrawRectangleRec(b.Rect, color)
	rl.DrawRectangleLinesEx(b.Rect, 1, rl.Black)
}
