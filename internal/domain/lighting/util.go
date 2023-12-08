package lighting

import (
	"github.com/polis-interactive/slate-tv/internal/util"
)

func generateLights(boards []util.BoardConfiguration, disallowed []util.Point) ([]util.Light, int) {
	// outputs
	lights := make([]util.Light, 0)
	lastLight := 0

	// temporary variables
	ledCount := 49
	maxNominalY := 6
	for _, config := range boards {
		if config.Type == util.Board7x7 {
			ledCount = 49
			maxNominalY = 6
		} else {
			ledCount = 7
			maxNominalY = 0
		}
		for i := 0; i < ledCount; i++ {
			yPos := i / 7
			isOddYPos := (yPos % 2) != 0
			xPos := i % 7
			if isOddYPos {
				xPos = 6 - xPos
			}

			position := config.StartingPoint
			isAllowed := true

			switch config.Orientation {
			case util.Orient0:
				position.AlterPoint(position.X+yPos, position.Y+(6-xPos))
			case util.Orient90:
				position.AlterPoint(position.X+xPos, position.Y+yPos)
			case util.Orient180:
				position.AlterPoint(position.X+(maxNominalY-yPos), position.Y+xPos)
			case util.Orient270:
				position.AlterPoint(position.X+(6-xPos), position.Y+(maxNominalY-yPos))
			}

			for _, p := range disallowed {
				if position.IsEqual(p) {
					isAllowed = false
					break
				}
			}

			l := util.Light{
				Position: position,
				Pixel:    lastLight,
				Show:     isAllowed,
				Color:    util.Color{},
			}
			lights = append(lights, l)
			lastLight += 1
		}
	}
	return lights, lastLight
}
