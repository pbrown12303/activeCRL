package crleditorfynegui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type fyneGuiTheme struct{}

var _ fyne.Theme = (*fyneGuiTheme)(nil)

func (m fyneGuiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == "inputBackground" {
		return color.Transparent
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (m fyneGuiTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m fyneGuiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m fyneGuiTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNamePadding {
		return 2.0
	}
	if name == theme.SizeNameInnerPadding {
		return 1.3
	}
	return theme.DefaultTheme().Size(name)
}
