package ai

import (
	"log"

	"github.com/fr3fou/flappy-go/flappy"
	"github.com/fr3fou/gone/gone"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	*flappy.Game
	Birds            []*Bird
	PopulationAmount int
	MutationRate     float64
	CrossoverOdds    float64
	Generation       int
}

var BrainLayers = []gone.Layer{
	{
		Nodes:     5,
		Activator: gone.Sigmoid(),
	},
	{
		Nodes:     8,
		Activator: gone.Sigmoid(),
	},
	{
		Nodes:     2,
		Activator: gone.Sigmoid(), // ideally should be softmax
	},
}

// New creates an AI simulation
func New(populationAmount int, mutationRate, crossoverOdds float64) *Game {
	birds := make([]*Bird, populationAmount)

	// Initialize population
	for i := range birds {
		birds[i] = NewBird(flappy.BirdSize*2, flappy.Height/2, gone.New(1, gone.MSE(), BrainLayers...))
	}

	g := &Game{
		Game:             flappy.NewGame(),
		Birds:            birds,
		Generation:       1,
		CrossoverOdds:    crossoverOdds,
		PopulationAmount: populationAmount,
		MutationRate:     mutationRate,
	}
	g.Init()

	return g
}

func (g *Game) Init() {
	g.Ground = flappy.NewGround()

	g.Pipes = make([]*flappy.Pipe, flappy.PipesBuffer)
	initialOffset := flappy.HorizontalGap + flappy.BirdSize*2*2
	offset := initialOffset
	for i := range g.Pipes {
		g.Pipes[i] = flappy.NewPipe((i)*flappy.PipeWidth+offset, 0)
		offset = flappy.HorizontalGap*(i+1) + initialOffset
	}
}

func (g *Game) Update() {
	for _, pipe := range g.Pipes {
		pipe.X -= flappy.Speed
	}

	for _, bird := range g.Birds {
		if !bird.Alive {
			continue
		}

		for _, pipe := range g.Pipes {
			if pipe.CollidesWith(bird.Rectangle) {
				bird.Alive = false
				break
			}
		}

		if g.Ground.CollidesWith(bird.Rectangle) || bird.AboveSky() {
			bird.Alive = false
		}
	}

	for _, bird := range g.Birds {
		if !bird.Alive {
			continue
		}

		// the next pipe is the closest one if the closestPipe is behind the bird
		closestPipe := g.Pipes[0]
		if g.Pipes[1].X < bird.X {
			closestPipe = g.Pipes[1]
		}

		if bird.ShouldJump(closestPipe) {
			bird.Jump()
		}

		bird.Update()
	}

	if g.Pipes[0].IsOffscreen() {
		newPipeX := int(g.Pipes[len(g.Pipes)-1].ToInt32().X + flappy.HorizontalGap + flappy.PipeWidth)
		g.Pipes = g.Pipes[1:]
		g.Pipes = append(g.Pipes, flappy.NewPipe(newPipeX, 0))
	}

	allDead := true
	for _, bird := range g.Birds {
		if bird.Alive {
			allDead = false
			break
		}
	}

	if allDead {
		g.Init()
		err := g.NextGeneration()
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.SkyBlue)

	for _, bird := range g.Birds {
		if bird.Alive {
			bird.Draw(100)
		}
	}

	for _, pipe := range g.Pipes {
		pipe.Draw()
	}

	g.Ground.Draw()
	rl.EndDrawing()
}
