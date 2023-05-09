package main

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func BuildToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), actionMenuOpen),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), actionSettings),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), actionHelpAbout),
	)
}
