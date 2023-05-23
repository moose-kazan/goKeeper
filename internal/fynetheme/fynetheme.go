package fynetheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type FyneTheme struct {
	colorsLight map[fyne.ThemeColorName]color.Color
}

func New() *FyneTheme {
	var ft FyneTheme
	ft.colorsLight = make(map[fyne.ThemeColorName]color.Color)
	ft.colorsLight[theme.ColorNameBackground] = color.RGBA{0xFF, 0xDA, 0xB9, 0xff}
	ft.colorsLight[theme.ColorNameOverlayBackground] = color.RGBA{0xFF, 0xDA, 0xB9, 0xff}
	ft.colorsLight[theme.ColorNameForeground] = color.RGBA{0x00, 0x00, 0x00, 0xff}
	ft.colorsLight[theme.ColorNamePrimary] = color.RGBA{0xFF, 0xA5, 0x00, 0xff}
	ft.colorsLight[theme.ColorNameSuccess] = color.RGBA{0xFF, 0xA5, 0x00, 0xff}
	ft.colorsLight[theme.ColorNameButton] = color.RGBA{0xF6, 0xE4, 0xD1, 0xFF}
	ft.colorsLight[theme.ColorNameMenuBackground] = color.RGBA{0xFF, 0xDA, 0xB9, 0xff}
	ft.colorsLight[theme.ColorNameInputBackground] = color.RGBA{0xF6, 0xE4, 0xD1, 0xFF}
	ft.colorsLight[theme.ColorNameDisabled] = color.RGBA{0x5A, 0x5A, 0x5A, 0xFF}
	return &ft
}

func (ft FyneTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	if variant == theme.VariantLight {
		if _, ok := ft.colorsLight[name]; ok {
			return ft.colorsLight[name]
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (ft FyneTheme) Font(name fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(name)
}
func (ft FyneTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (ft FyneTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
