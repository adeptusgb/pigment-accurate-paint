package main

import (
	"fmt"

	gui "pigmentaccuratepaint/internal/gui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth, screenHeight := int32(1500), int32(1000)
	canvasWidth, canvasHeight := int32(800), int32(800)
	rl.InitWindow(screenWidth, screenHeight, "Drawing App")
	rl.SetTargetFPS(60)

	// Create a canvas to draw on
	canvas := gui.NewCanvas(300, 100, canvasWidth, canvasHeight)

	// Brush settings
	brushColor := rl.Black
	brushRadius := float32(5)

	// Create color picker
	colorPicker := gui.NewColorPicker(10, 10)

	// Create clear canvas button
	clearButton := gui.NewButton(10, 300, 200, 30, rl.Red)

	// Create brush size buttons
	increaseSizeButton := gui.NewButton(80, 350, 50, 50, rl.Gray)
	decreaseSizeButton := gui.NewButton(10, 350, 50, 50, rl.Gray)

	for !rl.WindowShouldClose() {
		mousePos := rl.GetMousePosition()

		// Update color picker
		colorPicker.Update(mousePos)
		brushColor = colorPicker.GetSelectedColor()

		// Update brush size buttons
		if increaseSizeButton.Update(mousePos) {
			brushRadius += 1
			if brushRadius > 50 {
				brushRadius = 50
			}
		}

		if decreaseSizeButton.Update(mousePos) {
			brushRadius -= 1
			if brushRadius < 1 {
				brushRadius = 1
			}
		}

		// Update canvas
		canvas.Update(mousePos, brushColor, brushRadius)

		// Update clear button
		if clearButton.Update(mousePos) {
			rl.BeginTextureMode(canvas.RenderTexture)
			rl.ClearBackground(rl.White)
			rl.EndTextureMode()
		}

		// Draw everything
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw canvas
		canvas.Draw()

		// Draw UI elements
		colorPicker.Draw()
		clearButton.Draw()
		rl.DrawText("Clear", int32(clearButton.Rect.X+75), int32(clearButton.Rect.Y+8), 20, rl.White)

		increaseSizeButton.Draw()
		rl.DrawText("+", int32(increaseSizeButton.Rect.X+17), int32(increaseSizeButton.Rect.Y+11), 32, rl.Black)

		decreaseSizeButton.Draw()
		rl.DrawText("-", int32(decreaseSizeButton.Rect.X+17), int32(decreaseSizeButton.Rect.Y+11), 32, rl.Black)
		// Display current brush size
		sizeText := "Size: " + fmt.Sprintf("%d", int32(brushRadius))
		rl.DrawText(sizeText, 35, 410, 20, rl.Black)

		rl.EndDrawing()
	}

	rl.UnloadRenderTexture(canvas.RenderTexture)
	rl.CloseWindow()
}
