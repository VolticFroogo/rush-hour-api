package solver

import (
	"fmt"
)

type Position struct {
	Cars    uint64
	Bitmap  uint64
	History []uint8
}

func NewPosition(cars uint64, carOrientations uint64) (position Position) {
	position.Cars = cars

	// Add cars to bitmap.
	for i := 0; i < 16; i++ {
		pos := int((position.Cars >> (i * 4)) & 0xF)

		// If the car is off the grid, skip.
		// This is an error check, but also allows for cars to be manually skipped.
		// In which case their position should be 0x7 (111).
		if pos >= 6 {
			continue
		}

		orientation := int((carOrientations >> (i * 4)) & 0xF)

		position.Bitmap |= positionsToBitmask(pos, orientation)
		position.Bitmap |= positionsToBitmask(pos+1, orientation)

		if i >= 12 {
			position.Bitmap |= positionsToBitmask(pos+2, orientation)
		}
	}

	return
}

// Step adds all children positions to the queue (if they haven't already been computed),
// and returns a position if it has found a solution.
func (position *Position) Step(game *Game) *Position {
	for i := 0; i < 16; i++ {
		pos := int((position.Cars >> (i * 4)) & 0xF)

		// If the car is off the grid, skip.
		if pos >= 6 {
			continue
		}

		orientation := int((game.CarOrientations >> (i * 4)) & 0xF)

		for j := pos - 1; j >= 0; j-- {
			shouldBreak, solved := position.check(game, pos, orientation, 0, i, j)
			if shouldBreak {
				break
			}

			if solved != nil {
				return solved
			}
		}

		// If the car is a truck, give make it have an extra width.
		var width int
		if i < 12 {
			width = 1
		} else {
			width = 2
		}

		for j := pos + 1; j < 6-width; j++ {
			shouldBreak, solved := position.check(game, pos, orientation, width, i, j)
			if shouldBreak {
				break
			}

			if solved != nil {
				return solved
			}
		}
	}

	return nil
}

// check checks if a car move is possible.
// If a move is blocked, it will return true, nil; the loop should then be broken.
// If a solution is found, it will return false, pos, with the pos being the solution.
func (position *Position) check(game *Game, pos, orientation, width, i, j int) (bool, *Position) {
	// If there is a car blocking us from moving here, break.
	if position.Bitmap&positionsToBitmask(j+width, orientation) != 0 {
		return true, nil
	}

	// Get the new car position.
	cars := (position.Cars &^ (0xF << (i * 4))) | (uint64(j) << (i * 4))

	// If we have already computed this position before, skip.
	if game.Seen[cars] {
		return false, nil
	}

	// Add this position to seen.
	game.Seen[cars] = true

	size := len(position.History)
	history := make([]uint8, size+1)
	copy(history, position.History)
	history[size] = uint8((i << 4) | j)

	// Create the new position.
	newPosition := Position{
		Cars:    cars,
		Bitmap:  position.Bitmap,
		History: history,
	}

	// Clear the bits where the car was.
	newPosition.Bitmap &^= positionsToBitmask(pos, orientation)
	newPosition.Bitmap &^= positionsToBitmask(pos+1, orientation)
	// If the car is a truck, clear a third bit.
	if i >= 12 {
		newPosition.Bitmap &^= positionsToBitmask(pos+2, orientation)
	}

	// Add the new car's position to the bitmap.
	newPosition.Bitmap |= positionsToBitmask(j, orientation)
	newPosition.Bitmap |= positionsToBitmask(j+1, orientation)
	// If the car is a truck, add a third bit.
	if i >= 12 {
		newPosition.Bitmap |= positionsToBitmask(j+2, orientation)
	}

	// If the position is solved, return it.
	if newPosition.IsSolved() {
		return false, &newPosition
	}

	// Add this new position to the queue of positions to check.
	game.Positions = append(game.Positions, newPosition)
	return false, nil
}

// positionsToBitmask takes a position and orientation,
// returning a single high bit for where it is in the bitmap.
func positionsToBitmask(pos, orientation int) uint64 {
	if orientation&0x8 == 0 {
		return 1 << (pos + (orientation&0x7)*6)
	}

	return 1 << ((orientation & 0x7) + pos*6)
}

// IsSolved returns whether or not a position is in a solved state.
func (position *Position) IsSolved() bool {
	for i := position.Cars&0xF + 2; i < 6; i++ {
		if position.Bitmap&(1<<(i+6*3)) != 0 {
			return false
		}
	}

	return true
}

// DebugPrintBitmap prints the bitmap in a human-readable way to console.
func (position *Position) DebugPrintBitmap() {
	fmt.Printf("%06b\n%06b\n%06b\n%06b\n%06b\n%06b\n",
		(position.Bitmap>>30)&0x3F,
		(position.Bitmap>>24)&0x3F,
		(position.Bitmap>>18)&0x3F,
		(position.Bitmap>>12)&0x3F,
		(position.Bitmap>>6)&0x3F,
		(position.Bitmap>>0)&0x3F,
	)
}

func (position *Position) DebugPrintMoves() {
	for i := range position.History {
		fmt.Printf("Car %2d to %d\n", position.History[i]>>4, position.History[i]&0xF)
	}
}
