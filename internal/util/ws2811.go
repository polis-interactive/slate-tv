package util

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"math"
)

type GpioPinType int

const (
	gpio12 GpioPinType = 12
	gpio18             = 18
	gpio13             = 13
	gpio19             = 19
)

var GpioPinTypes = struct {
	GPIO12 GpioPinType
	GPIO18 GpioPinType
	GPIO13 GpioPinType
	GPIO19 GpioPinType
}{
	GPIO12: gpio12,
	GPIO18: gpio18,
	GPIO13: gpio13,
	GPIO19: gpio19,
}

type StripType int

const (
	ws2811rgb StripType = ws2811.WS2811StripRGB
	ws2811rbg           = ws2811.WS2811StripRBG
	ws2811grb           = ws2811.WS2811StripGRB
	ws2811gbr           = ws2811.WS2811StripGBR
	ws2811brg           = ws2811.WS2811StripBRG
	ws2811bgr           = ws2811.WS2811StripBGR

	sk6812rgbw = ws2811.SK6812StripRGBW
	sk6812rbgw = ws2811.SK6812StripRBGW
	sk6812grbw = ws2811.SK6812StripGRBW
	sk6812gbrw = ws2811.SK6812StrioGBRW
	sk6812brgw = ws2811.SK6812StrioBRGW
	sk6812bgrw = ws2811.SK6812StripBGRW
)

var StripTypes = struct {
	WS2811RGB StripType
	WS2811RBG StripType
	WS2811GRB StripType
	WS2811GBR StripType
	WS2811BRG StripType
	WS2811BGR StripType

	SK6812RGBW StripType
	SK6812RBGW StripType
	SK6812GRBW StripType
	SK6812GBRW StripType
	SK6812BRGW StripType
	SK6812BGRW StripType
}{
	WS2811RGB: ws2811rgb,
	WS2811RBG: ws2811rbg,
	WS2811GRB: ws2811grb,
	WS2811GBR: ws2811gbr,
	WS2811BRG: ws2811brg,
	WS2811BGR: ws2811bgr,

	SK6812RGBW: sk6812rgbw,
	SK6812RBGW: sk6812rbgw,
	SK6812GRBW: sk6812grbw,
	SK6812GBRW: sk6812gbrw,
	SK6812BRGW: sk6812brgw,
	SK6812BGRW: sk6812bgrw,
}

func MakeGammaTable(gamma float64) []byte {
	gt := make([]byte, 256)
	gmdMax := math.Pow(255, gamma)
	for i := 0; i < 256; i++ {
		gmd := math.Pow(float64(i), gamma)
		gmdNorm := math.Round(gmd / gmdMax * 255.0)
		gt[i] = byte(math.Min(255.0, gmdNorm))
	}
	return gt
}
