package bus

import (
	"github.com/polis-interactive/slate-tv/internal/util"
	"sync"
)

func (b *bus) GetLightCount() int {
	return b.lightingService.GetLightCount()
}

func (b *bus) CopyLightsToColorBuffer(rawPbOut []util.Color) error {
	lights, preLockedLightsMutex := b.lightingService.GetLights()
	pbIn, preLockedGraphicsMutex := b.graphicsService.GetPb()
	defer func(lightsMu *sync.RWMutex, graphicsMu *sync.RWMutex) {
		lightsMu.RUnlock()
		graphicsMu.RUnlock()
	}(preLockedLightsMutex, preLockedGraphicsMutex)
	for _, l := range lights {
		if !l.Show {
			continue
		}
		rawPbOut[l.Pixel] = pbIn.GetPixel(&l.Position)
	}
	return nil
}

func (b *bus) CopyLightsToUint32Buffer(rawUint32BuffOut []uint32) error {
	lights, preLockedLightsMutex := b.lightingService.GetLights()
	pbIn, preLockedGraphicsMutex := b.graphicsService.GetPb()
	defer func(lightsMu *sync.RWMutex, graphicsMu *sync.RWMutex) {
		lightsMu.RUnlock()
		graphicsMu.RUnlock()
	}(preLockedLightsMutex, preLockedGraphicsMutex)
	for _, l := range lights {
		if !l.Show {
			continue
		}
		rawUint32BuffOut[l.Pixel] = pbIn.GetPixelPointer(&l.Position).ToBits()
	}
	return nil
}
