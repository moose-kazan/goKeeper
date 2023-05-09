package main

import (
	"fmt"
	"gokeeperViewer/internal/settings"
	"net/url"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func actionHelpAbout() {
	urlEmailTitle := "moose@ylsoftware.com"
	urlEmail, _ := url.Parse(fmt.Sprintf("mailto:%s", urlEmailTitle))
	urlWSTitle := "https://github.com/moose-kazan/goKeeperViewer"
	urlWS, _ := url.Parse("https://github.com/moose-kazan/goKeeperViewer")
	aboutLayout := container.NewVBox(
		widget.NewLabelWithStyle("goKeeperViewer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Author", widget.NewLabel("Vadim Kalinnov")),
			widget.NewFormItem("E-Mail", widget.NewHyperlink(urlEmailTitle, urlEmail)),
			widget.NewFormItem("Website", widget.NewHyperlink(urlWSTitle, urlWS)),
			widget.NewFormItem("Description", widget.NewLabel("Simple viewer for KDBX (KeePass) files.")),
			widget.NewFormItem("OS", widget.NewLabel(fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH))),
		),
	)
	d := dialog.NewCustom(
		"About",
		"OK",
		aboutLayout,
		w,
	)
	d.Show()
}

func actionMenuOpen() {
	d := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
		if e != nil {
			dialog.NewError(e, w).Show()
			return
		}
		if u != nil {
			loadFile(u.URI())
		}

	}, w)
	d.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	d.Show()
}

func actionSettings() {
	selectItem := widget.NewSelect(
		settings.StartLoadOptions(),
		func(s string) {

		},
	)
	selectItem.SetSelectedIndex(settings.New(a.Preferences()).GetStartLoadOption())

	dialog.NewForm(
		"Settings",
		"OK",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem(
				"Load on start",
				selectItem,
			),
		},
		func(b bool) {
			if !b {
				return
			}
			settings.New(a.Preferences()).SetStartLoadOption(selectItem.Selected)
		},
		w,
	).Show()
	return
}
