package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// forcedTheme — это наша "обертка". Она берет стандартную тему,
// но принудительно меняет вариант (светлый/темный)
type forcedTheme struct {
	fyne.Theme
	variant fyne.ThemeVariant
}

// Color переопределяет цвет. Каким бы ни был запрос (v),
// мы подсовываем свой вариант (t.variant)
func (t *forcedTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return t.Theme.Color(n, t.variant)
}

// NewForcedTheme создает тему с принудительным режимом
func NewForcedTheme(variant fyne.ThemeVariant) fyne.Theme {
	return &forcedTheme{
		Theme:   theme.DefaultTheme(),
		variant: variant,
	}
}
