package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// transparentTheme is the standard Fyne theme with a see-through background.
// Fyne clears the canvas with the background colour, so making it transparent
// leaves the desktop showing everywhere the eyes are not painted.
type transparentTheme struct {
	fyne.Theme
}

func (t *transparentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.Transparent
	}
	return t.Theme.Color(name, variant)
}
