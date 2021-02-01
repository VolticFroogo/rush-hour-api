package v1

import (
	"github.com/VolticFroogo/rush-hour-api/solver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type solveReq struct {
	Cars            [16]uint8 `json:"cars"`
	CarOrientations [16]uint8 `json:"carOrientations"`
}

type solveRes struct {
	Status        int   `json:"status"`
	Steps         []int `json:"steps,omitempty"`
	ExecutionTime int64 `json:"executionTime"`
}

const (
	solveStatusSolved = iota
	solveStatusUnsolvable
)

func solve(c *gin.Context) {
	var req solveReq
	err := c.BindJSON(&req)
	if err != nil {
		throwError(c, http.StatusBadRequest, errBadJSON)
		return
	}

	var cars, carOrientations uint64 = 0, 0
	for i := 0; i < 16; i++ {
		cars |= uint64(req.Cars[i]) << (i * 4)
		carOrientations |= uint64(req.CarOrientations[i]) << (i * 4)
	}

	game := solver.NewGame(cars, carOrientations)

	before := time.Now()
	solution := game.Solve()
	delta := time.Now().Sub(before).Microseconds()

	if solution == nil {
		c.JSON(http.StatusOK, solveRes{
			Status:        solveStatusUnsolvable,
			ExecutionTime: delta,
		})
		return
	}

	steps := make([]int, len(solution.History))
	for i := range solution.History {
		steps[i] = int(solution.History[i])
	}

	c.JSON(http.StatusOK, solveRes{
		Status:        solveStatusSolved,
		Steps:         steps,
		ExecutionTime: delta,
	})
}
