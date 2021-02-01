package main

import (
	"fmt"
	"github.com/VolticFroogo/rush-hour-api/solver"
	"time"
)

const (
	game1Cars          uint64 = 0x3222777777774001
	game1Orientations  uint64 = 0xD0B8000000001853
	game40Cars         uint64 = 0x7203773004133413
	game40Orientations uint64 = 0x0D28000A01B9AC53
)

func main() {
	game := solver.NewGame(game40Cars, game40Orientations)
	fmt.Println("Initial bitmap:")
	game.Positions[0].DebugPrintBitmap()

	before := time.Now()
	solution := game.Solve()
	delta := time.Now().Sub(before)
	fmt.Println("\nTime taken to solve game:", delta)

	fmt.Println("Positions checked:", len(game.Seen))
	if solution == nil {
		fmt.Println("\nThis game is impossible.")
		return
	}

	fmt.Println("\nSolved bitmap:")
	solution.DebugPrintBitmap()

	fmt.Println("\nSolution steps:")
	solution.DebugPrintMoves()
}
