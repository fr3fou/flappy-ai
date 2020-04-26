package main

import (
	"fmt"

	"github.com/fr3fou/flappy-ai/ai"
	"github.com/fr3fou/flappy-go/flappy"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(flappy.Width, flappy.Height, "Flappy AI!")
	g := ai.New(500, 0.1, 0.1)

	rl.SetTargetFPS(60)
	sliderValue := float32(60)
	prevSliderValue := sliderValue

	hasStarted := false
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		if hasStarted {
			g.Update()
		}
		g.Draw()

		if !hasStarted && rl.IsKeyPressed(rl.KeySpace) {
			hasStarted = true
			continue
		}

		raygui.Label(rl.NewRectangle(15, 50, 200, 20), "FPS")
		sliderValue = raygui.Slider(rl.NewRectangle(15, 70, 200, 20), sliderValue, 60, 240)
		raygui.Label(rl.NewRectangle(15+200+5, 70, 20, 20), fmt.Sprintf("%.0f", sliderValue))
		if prevSliderValue != sliderValue {
			rl.SetTargetFPS(int32(sliderValue))
			prevSliderValue = sliderValue
		}
		rl.EndDrawing()
	}
}
