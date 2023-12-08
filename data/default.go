package data

import "github.com/polis-interactive/slate-tv/internal/util"

var BoardConfiguration = []util.BoardConfiguration{
	util.NewBoardConfiguration(util.Board7x7, util.Orient0, util.Point{X: 0, Y: 0}),
	util.NewBoardConfiguration(util.Board7x7, util.Orient0, util.Point{X: 0, Y: 7}),
	util.NewBoardConfiguration(util.Board7x7, util.Orient0, util.Point{X: 7, Y: 7}),
	util.NewBoardConfiguration(util.Board7x7, util.Orient0, util.Point{X: 7, Y: 0}),
}

var BoardDisallowedPositions []util.Point
