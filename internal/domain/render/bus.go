package render

import "github.com/polis-interactive/slate-tv/internal/util"

type Bus interface {
	GetLightCount() int
	CopyLightsToColorBuffer(buff []util.Color) error
	CopyLightsToUint32Buffer(buff []uint32) error
}
