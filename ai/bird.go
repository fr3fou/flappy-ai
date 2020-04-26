package ai

import (
	"math"

	"github.com/fr3fou/flappy-go/flappy"
	"github.com/fr3fou/gone/gone"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bird struct {
	flappy.Bird
	Brain   *gone.NeuralNetwork
	Score   int
	Fitness float64
	Alive   bool
}

func NewBird(x, y int, brain *gone.NeuralNetwork) *Bird {
	return &Bird{
		Bird:  *flappy.NewBird(x, y),
		Alive: true,
		Brain: brain,
	}
}

func (b *Bird) Update() {
	b.Score++
	b.Bird.Update()
}

func (b *Bird) ShouldJump(pipe *flappy.Pipe) bool {
	// don't ask why i used greek letters as inputs
	// it looks fancier
	// it looks better than x1, x2, x3, x4
	// idk
	// explanation of the variables here
	// https://cdn.discordapp.com/attachments/361910177961738244/703687897886228600/JPEG_20200425_222426.jpg
	alpha := float64(b.Y) / flappy.Height
	beta := math.Abs(float64(b.X-pipe.X)) / flappy.Width
	gamma := math.Abs(float64(b.Y-pipe.Height+flappy.PipeBorder)) / flappy.Height
	delta := math.Abs(float64(b.Y-flappy.Height-pipe.Height+flappy.PipeBorder)) / flappy.Height
	epsilon := float64(b.Velocity / flappy.MaxVelocity)

	inputs := []float64{alpha, beta, gamma, delta, epsilon}
	outputs := b.Brain.Predict(inputs)

	return outputs[0] > outputs[1]
}

func (b *Bird) Draw(opacity uint8) {
	if b.Alive {
		b.Bird.Draw(rl.NewColor(rl.Yellow.R, rl.Yellow.G, rl.Yellow.B, opacity))
	} else {
		b.Bird.Draw(rl.NewColor(rl.Red.R, rl.Red.G, rl.Red.B, opacity))
	}
}
