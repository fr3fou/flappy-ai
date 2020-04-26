package main

import (
	"github.com/fr3fou/flappy-go/ai"
	"github.com/fr3fou/flappy-go/flappy"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(flappy.Width, flappy.Height, "Flappy Bird!")
	g := ai.New(500, 0.1, 0.1)

	rl.SetTargetFPS(60)

	hasStarted := false
	for !rl.WindowShouldClose() {
		if hasStarted {
			g.Update()
		}
		g.Draw()

		if !hasStarted && rl.IsKeyPressed(rl.KeySpace) {
			hasStarted = true
			continue
		}
	}
}
