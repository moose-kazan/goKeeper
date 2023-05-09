package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func newMenuItem(label string, action func(), Icon fyne.Resource, Shortcut fyne.Shortcut) *fyne.MenuItem {
	m := fyne.NewMenuItem(label, action)
	m.Icon = Icon
	m.Shortcut = Shortcut
	return m
}

func BuildMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu(
			"File",
			newMenuItem("Open", actionMenuOpen, theme.DocumentIcon(), nil),
			newMenuItem("Quit", func() { a.Quit() }, theme.LogoutIcon(), nil),
		),
		fyne.NewMenu(
			"Settings",
			newMenuItem("Settings", actionSettings, theme.SettingsIcon(), nil),
		),
		fyne.NewMenu(
			"Help",
			newMenuItem("About", actionHelpAbout, theme.InfoIcon(), nil),
		),
	)
}
