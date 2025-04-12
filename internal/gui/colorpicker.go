package gui

import (
	"math"

	clr "pigmentaccuratepaint/internal/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ColorPicker struct {
	Rect          rl.Rectangle
	ColorWheel    rl.Rectangle
	ValueSlider   rl.Rectangle
	PreviewRect   rl.Rectangle
	SelectedColor rl.Color
	SelectedH     float32
	SelectedS     float32
	SelectedV     float32
	IsDraggingHue bool
	IsDraggingVal bool
}

func NewColorPicker(x, y float32) ColorPicker {
	size := rl.Vector2{X: 200, Y: 270}
	wheelRect := rl.Rectangle{X: x + 10, Y: y + 10, Width: 180, Height: 180}
	valSliderRect := rl.Rectangle{X: x + 10, Y: y + 200, Width: 180, Height: 20}
	previewRect := rl.Rectangle{X: x + 10, Y: y + 230, Width: 30, Height: 30}

	return ColorPicker{
		Rect:          rl.Rectangle{X: x, Y: y, Width: size.X, Height: size.Y},
		ColorWheel:    wheelRect,
		ValueSlider:   valSliderRect,
		PreviewRect:   previewRect,
		SelectedColor: rl.Black,
		SelectedH:     0,
		SelectedS:     0,
		SelectedV:     0,
	}
}

func (cp *ColorPicker) Update(mousePos rl.Vector2) {
	// Check for color wheel interaction
	inColorWheel := rl.CheckCollisionPointRec(mousePos, cp.ColorWheel)

	if inColorWheel && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		centerX := cp.ColorWheel.X + cp.ColorWheel.Width/2
		centerY := cp.ColorWheel.Y + cp.ColorWheel.Height/2
		radius := cp.ColorWheel.Width / 2

		dx := mousePos.X - centerX
		dy := mousePos.Y - centerY

		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if distance > radius {
			// Normalize to edge of wheel
			dx = dx * radius / distance
			dy = dy * radius / distance
			distance = radius
		}

		// Calculate hue and saturation
		hue := float32(math.Atan2(float64(dy), float64(dx))) * 180 / float32(math.Pi)
		if hue < 0 {
			hue += 360
		}
		saturation := distance / radius

		cp.SelectedH = hue
		cp.SelectedS = saturation
		cp.IsDraggingHue = true
	} else {
		cp.IsDraggingHue = cp.IsDraggingHue && rl.IsMouseButtonDown(rl.MouseLeftButton)
	}

	// Check for value slider interaction
	inValueSlider := rl.CheckCollisionPointRec(mousePos, cp.ValueSlider)

	if inValueSlider && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		// Calculate value based on x position in slider
		relativeX := mousePos.X - cp.ValueSlider.X
		cp.SelectedV = relativeX / cp.ValueSlider.Width
		if cp.SelectedV < 0 {
			cp.SelectedV = 0
		}
		if cp.SelectedV > 1 {
			cp.SelectedV = 1
		}
		cp.IsDraggingVal = true
	} else {
		cp.IsDraggingVal = cp.IsDraggingVal && rl.IsMouseButtonDown(rl.MouseLeftButton)
	}

	// Update selected color
	cp.SelectedColor = clr.HsvToRgb(cp.SelectedH, cp.SelectedS, cp.SelectedV)
}

func (cp *ColorPicker) Draw() {
	rl.DrawRectangleRec(cp.Rect, rl.ColorAlpha(rl.LightGray, 0.9))
	rl.DrawRectangleLinesEx(cp.Rect, 1, rl.Black)

	// Draw color wheel
	centerX := cp.ColorWheel.X + cp.ColorWheel.Width/2
	centerY := cp.ColorWheel.Y + cp.ColorWheel.Height/2
	radius := cp.ColorWheel.Width / 2

	// Draw color wheel
	for x := int32(cp.ColorWheel.X); x < int32(cp.ColorWheel.X+cp.ColorWheel.Width); x++ {
		for y := int32(cp.ColorWheel.Y); y < int32(cp.ColorWheel.Y+cp.ColorWheel.Height); y++ {
			dx := float32(x) - centerX
			dy := float32(y) - centerY

			distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			if distance <= radius {
				hue := float32(math.Atan2(float64(dy), float64(dx))) * 180 / float32(math.Pi)
				if hue < 0 {
					hue += 360
				}
				saturation := distance / radius

				color := clr.HsvToRgb(hue, saturation, 1.0)
				rl.DrawPixel(x, y, color)
			}
		}
	}

	// Draw selected position on color wheel
	selectedPosX := centerX + radius*cp.SelectedS*float32(math.Cos(float64(cp.SelectedH)*math.Pi/180))
	selectedPosY := centerY + radius*cp.SelectedS*float32(math.Sin(float64(cp.SelectedH)*math.Pi/180))
	rl.DrawCircleV(rl.Vector2{X: selectedPosX, Y: selectedPosY}, 5, rl.White)
	rl.DrawCircleLines(int32(selectedPosX), int32(selectedPosY), 5, rl.Black)

	// Draw value slider
	for x := int32(cp.ValueSlider.X); x < int32(cp.ValueSlider.X+cp.ValueSlider.Width); x++ {
		relX := float32(x-int32(cp.ValueSlider.X)) / cp.ValueSlider.Width
		rl.DrawLine(
			x,
			int32(cp.ValueSlider.Y),
			x,
			int32(cp.ValueSlider.Y+cp.ValueSlider.Height),
			clr.HsvToRgb(cp.SelectedH, cp.SelectedS, relX),
		)
	}

	// Draw value slider marker
	markerX := cp.ValueSlider.X + cp.SelectedV*cp.ValueSlider.Width
	rl.DrawRectangle(
		int32(markerX-2),
		int32(cp.ValueSlider.Y),
		4,
		int32(cp.ValueSlider.Height),
		rl.White,
	)
	rl.DrawRectangleLines(
		int32(markerX-2),
		int32(cp.ValueSlider.Y),
		4,
		int32(cp.ValueSlider.Height),
		rl.Black,
	)

	// Draw selected color preview
	rl.DrawRectangleRec(cp.PreviewRect, cp.SelectedColor)
}

func (cp *ColorPicker) GetSelectedColor() rl.Color {
	return cp.SelectedColor
}
