package ai

import (
	"log"
	"math"
	"math/rand"

	"github.com/fr3fou/flappy-go/flappy"
	"github.com/fr3fou/gone/gone"
)

func (g *Game) NextGeneration() error {
	g.Generation++
	log.Printf("Generation Number %d...", g.Generation)

	// penalize worse scores
	// reward higher scores
	sum := 0.0
	for _, bird := range g.Birds {
		score := int(math.Pow(float64(bird.Score), 2))
		bird.Score = score
		sum += float64(score)
	}

	// Maybe there's a better way?
	for _, bird := range g.Birds {
		bird.Fitness = float64(bird.Score) / sum
	}

	// Make the new population
	newBirds := make([]*Bird, g.PopulationAmount)
	for i := 0; i < g.PopulationAmount; i++ {
		// Determine whether to do crossover or a copy of the best one
		var child *Bird
		if rand.Float64() > g.CrossoverOdds {
			child = g.NaturalSelection()
		} else {
			firstParent := g.NaturalSelection()
			secondParent := g.NaturalSelection()

			childBrain, err := firstParent.Brain.Crossover(secondParent.Brain)
			if err != nil {
				return err
			}

			child = NewBird(flappy.BirdSize*2, flappy.Height/2, childBrain)
		}

		child.Brain.Mutate(
			gone.GaussianMutation(g.MutationRate, 1, 0),
		)
		newBirds[i] = child
	}
	g.Birds = newBirds

	return nil
}

func (g *Game) NaturalSelection() *Bird {
	// https://github.com/CodingTrain/Toy-Neural-Network-JS/blob/master/examples/neuroevolution-flappybird/ga.js#L63-L87
	i := 0
	probability := rand.Float64()
	for probability > 0 {
		probability -= g.Birds[i].Fitness
		i++
	}
	i--
	bird := g.Birds[i]
	return NewBird(flappy.BirdSize*2, flappy.Height/2, bird.Brain.Copy())
}
