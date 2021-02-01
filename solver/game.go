package solver

type Game struct {
	CarOrientations uint64
	Seen            map[uint64]bool
	Positions       []Position
}

// NewGame creates a new game given cars and their orientations.
func NewGame(cars, carOrientations uint64) (game Game) {
	game.CarOrientations = carOrientations

	game.Positions = []Position{
		NewPosition(cars, carOrientations),
	}

	game.Seen = map[uint64]bool{
		game.Positions[0].Cars: true,
	}

	return
}

// Solve solves a game, returns a position if it is solvable,
// or nil if the game is impossible.
func (game *Game) Solve() *Position {
	for len(game.Positions) != 0 {
		positions := game.Positions
		game.Positions = []Position{}

		for i := range positions {
			solution := positions[i].Step(game)
			if solution != nil {
				return solution
			}
		}
	}

	return nil
}
