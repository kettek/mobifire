package layouts

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct{}

func (m Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameShadow {
		return color.RGBA{0, 0, 0, 0} // Blank out shadows for chat area.
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m Theme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m Theme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInlineIcon {
		return theme.DefaultTheme().Size(name) * 2 // I guess
	} else if name == theme.SizeNameInnerPadding {
		return 0
	} else if name == theme.SizeNameSeparatorThickness {
		return 0
	}
	return theme.DefaultTheme().Size(name)
}

type NoPaddingTheme struct{}

func (m NoPaddingTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m NoPaddingTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m NoPaddingTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m NoPaddingTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInnerPadding || name == theme.SizeNamePadding {
		return 0
	}
	if name == theme.SizeNameInlineIcon {
		return theme.DefaultTheme().Size(name) * 2
	}
	return theme.DefaultTheme().Size(name)
}
