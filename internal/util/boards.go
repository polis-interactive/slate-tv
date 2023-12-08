package util

type BoardType int

const (
	Board7x7 BoardType = 0
	Board1x7           = 1
)

type BoardOrientation int

const (
	Orient0   BoardOrientation = 0
	Orient90  BoardOrientation = 1
	Orient180 BoardOrientation = 2
	Orient270 BoardOrientation = 3
)

type BoardConfiguration struct {
	Type          BoardType
	Orientation   BoardOrientation
	StartingPoint Point
}

func NewBoardConfiguration(t BoardType, o BoardOrientation, startingPoint Point) BoardConfiguration {
	return BoardConfiguration{
		Type:          t,
		Orientation:   o,
		StartingPoint: startingPoint,
	}
}
