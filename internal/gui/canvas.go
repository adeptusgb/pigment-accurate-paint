package gui

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Canvas represents a drawable area on the screen
type Canvas struct {
	RenderTexture rl.RenderTexture2D
	Rect          rl.Rectangle

	rawPoints []rl.Vector2
	lastPoint rl.Vector2
	dragging  bool
}

// NewCanvas creates a new canvas with the specified position and size.
// It initializes the render texture and clears it with a white background
func NewCanvas(x, y float32, width, height int32) Canvas {
	canvas := Canvas{
		RenderTexture: rl.LoadRenderTexture(width, height),
		Rect:          rl.Rectangle{X: x, Y: y, Width: float32(width), Height: float32(height)},

		rawPoints: []rl.Vector2{},
		lastPoint: rl.Vector2{},
		dragging:  false,
	}
	rl.BeginTextureMode(canvas.RenderTexture)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()

	return canvas
}

func (c *Canvas) Draw() {
	// raylib uses a different coordinate system for textures, since they run on OpenGL.
	// OpenGL has the origin at the bottom left corner of the screen,
	// so we need to flip the Y-axis to correct raylib's coordinate system for textures.
	srcRect := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(c.RenderTexture.Texture.Width),
		Height: -float32(c.RenderTexture.Texture.Height),
	}

	rl.DrawTexturePro(
		c.RenderTexture.Texture,
		srcRect,
		c.Rect,
		rl.Vector2{X: 0, Y: 0},
		0.0,
		rl.White,
	)
}

func (c *Canvas) Update(mousePos rl.Vector2, brushColor rl.Color, brushRadius float32) {
	if !rl.IsMouseButtonDown(rl.MouseLeftButton) {
		c.dragging = false
		c.rawPoints = []rl.Vector2{} // Clear raw points when mouse is released
		c.lastPoint = rl.Vector2{}   // Reset last point

		return
	}

	// Check if the mouse is inside the canvas
	if !rl.CheckCollisionPointRec(mousePos, c.Rect) {
		// todo: should draw up to the edge of the canvas and then reset the state
		c.dragging = false
		c.rawPoints = []rl.Vector2{} // Clear raw points when mouse is released
		c.lastPoint = rl.Vector2{}   // Reset last point

		fmt.Println("Mouse was clicked outside the canvas!")
		return
	}

	relativeMousePos := rl.Vector2{
		X: mousePos.X - c.Rect.X,
		Y: mousePos.Y - c.Rect.Y,
	}

	// Drawing logic
	// todo: move this to a drawing package
	maxPointSpacing := brushRadius * 4 // Maximum distance between points before interpolation

	if !c.dragging {
		// Start new stroke
		c.dragging = true
		c.rawPoints = []rl.Vector2{relativeMousePos}
		c.lastPoint = relativeMousePos

		// Draw the first point as a circle
		rl.BeginTextureMode(c.RenderTexture)
		rl.DrawCircleV(relativeMousePos, brushRadius, brushColor)
		rl.EndTextureMode()
	} else {
		// Check if we need to interpolate points due to fast movement
		if distance(c.lastPoint, relativeMousePos) > maxPointSpacing {
			interpolated := interpolatePoints(c.lastPoint, relativeMousePos, maxPointSpacing)

			// Add interpolated points to raw points
			for _, point := range interpolated {
				c.rawPoints = append(c.rawPoints, point)

				// Draw each interpolated point
				rl.BeginTextureMode(c.RenderTexture)
				rl.DrawCircleV(point, brushRadius, brushColor)
				rl.EndTextureMode()
			}

			c.lastPoint = relativeMousePos
		} else {
			// Add current point normally
			c.rawPoints = append(c.rawPoints, relativeMousePos)
			c.lastPoint = relativeMousePos

			// Draw point
			rl.BeginTextureMode(c.RenderTexture)
			rl.DrawCircleV(relativeMousePos, brushRadius, brushColor)
			rl.EndTextureMode()
		}
	}

	// If enough points, smooth and render
	if len(c.rawPoints) >= 4 {
		smooth := smoothPoints(c.rawPoints)

		rl.BeginTextureMode(c.RenderTexture)
		// Draw circles along the smoothed path
		for i := 0; i < len(smooth); i++ {
			rl.DrawCircleV(smooth[i], brushRadius, brushColor)
		}
		rl.EndTextureMode()
	}
}

// Linear interpolation between two points
func lerp(p1, p2 rl.Vector2, t float32) rl.Vector2 {
	return rl.Vector2{
		X: p1.X + (p2.X-p1.X)*t,
		Y: p1.Y + (p2.Y-p1.Y)*t,
	}
}

// Calculate distance between two points
func distance(p1, p2 rl.Vector2) float32 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Generate points between p1 and p2 with maximum spacing
func interpolatePoints(p1, p2 rl.Vector2, maxSpacing float32) []rl.Vector2 {
	dist := distance(p1, p2)
	if dist <= maxSpacing {
		return []rl.Vector2{p2} // Just return the end point
	}

	steps := int(math.Ceil(float64(dist / maxSpacing)))
	var result []rl.Vector2

	for i := 1; i <= steps; i++ {
		t := float32(i) / float32(steps)
		result = append(result, lerp(p1, p2, t))
	}

	return result
}

// CatmullRom interpolation
func CatmullRom(p0, p1, p2, p3 rl.Vector2, t float32) rl.Vector2 {
	t2 := t * t
	t3 := t2 * t
	return rl.Vector2{
		X: 0.5 * ((2 * p1.X) +
			(-p0.X+p2.X)*t +
			(2*p0.X-5*p1.X+4*p2.X-p3.X)*t2 +
			(-p0.X+3*p1.X-3*p2.X+p3.X)*t3),
		Y: 0.5 * ((2 * p1.Y) +
			(-p0.Y+p2.Y)*t +
			(2*p0.Y-5*p1.Y+4*p2.Y-p3.Y)*t2 +
			(-p0.Y+3*p1.Y-3*p2.Y+p3.Y)*t3),
	}
}

// smoothPoints generates interpolated points
func smoothPoints(points []rl.Vector2) []rl.Vector2 {
	var smooth []rl.Vector2
	const steps = 30 // More steps = smoother lines
	if len(points) < 4 {
		return points
	}
	for i := 0; i < len(points)-1; i++ {
		p0 := points[int(math.Max(float64(i-1), 0))]
		p1 := points[i]
		p2 := points[i+1]
		p3 := points[int(math.Min(float64(i+2), float64(len(points)-1)))]
		// Interpolate between p1 and p2
		for step := 0; step < steps; step++ {
			t := float32(step) / float32(steps)
			smooth = append(smooth, CatmullRom(p0, p1, p2, p3, t))
		}
	}
	smooth = append(smooth, points[len(points)-1])
	return smooth
}
