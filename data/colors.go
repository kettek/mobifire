package data

import (
	"image/color"

	"github.com/kettek/termfire/messages"
)

// Color returns a color from a CF color message.
func Color(c messages.MessageColor) color.Color {
	switch c {
	case messages.MessageColorBlack:
		return color.Black
	case messages.MessageColorWhite:
		return color.White
	case messages.MessageColorNavy:
		return color.NRGBA{0, 0, 128, 255}
	case messages.MessageColorRed:
		return color.NRGBA{255, 0, 0, 255}
	case messages.MessageColorOrange:
		return color.NRGBA{255, 165, 0, 255}
	case messages.MessageColorBlue:
		return color.NRGBA{0, 0, 255, 255}
	case messages.MessageColorDarkOrange:
		return color.NRGBA{255, 140, 0, 255}
	case messages.MessageColorGreen:
		return color.NRGBA{0, 128, 0, 255}
	case messages.MessageColorLightGreen:
		return color.NRGBA{144, 238, 144, 255}
	case messages.MessageColorGrey:
		return color.NRGBA{128, 128, 128, 255}
	case messages.MessageColorBrown:
		return color.NRGBA{165, 42, 42, 255}
	case messages.MessageColorGold:
		return color.NRGBA{255, 215, 0, 255}
	case messages.MessageColorTan:
		return color.NRGBA{210, 180, 140, 255}
	}
	return color.Black
}
