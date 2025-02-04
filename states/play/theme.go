package play

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameShadow {
		return color.RGBA{0, 0, 0, 0} // Blank out shadows for chat area.
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInlineIcon {
		return theme.DefaultTheme().Size(name) * 2 // I guess
	} else if name == theme.SizeNameInnerPadding || name == theme.SizeNameLineSpacing {
		return 0
	} else if name == theme.SizeNameSeparatorThickness {
		return 0
	}
	return theme.DefaultTheme().Size(name)
}
