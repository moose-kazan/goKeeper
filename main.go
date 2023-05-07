package main

import (
	"fmt"
	"gokeeperViewer/fynefilechooser"
	"gokeeperViewer/internal/kdb"
	"gokeeperViewer/internal/settings"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	a               fyne.App
	w               fyne.Window
	passwordTree    *widget.Tree
	passwordDetails *widget.Form
	db              *kdb.KDB
)

func newMenuItem(label string, action func(), Icon fyne.Resource, Shortcut fyne.Shortcut) *fyne.MenuItem {
	m := fyne.NewMenuItem(label, action)
	m.Icon = Icon
	m.Shortcut = Shortcut
	return m
}

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

func loadFile(fileName fyne.URI) {
	pwdEntry := widget.NewPasswordEntry()
	keyFileChooser := fynefilechooser.NewFileChooser(w, storage.NewExtensionFileFilter([]string{".keyx", ".key"}))
	d := dialog.NewForm(
		"Enter password",
		"OK",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("File Name", widget.NewLabel(filepath.Base(fileName.Path()))),
			widget.NewFormItem("Password", pwdEntry),
			widget.NewFormItem("Key File", keyFileChooser),
		},
		func(b bool) {
			if !b {
				return
			}

			tmpDb := kdb.New()
			err := tmpDb.Load(fileName, pwdEntry.Text, keyFileChooser.GetURI())

			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			_ = tmpDb.Tree()

			db = tmpDb
			//db.SetDebug(true)
			passwordTree.Root = "/"
			passwordTree.Refresh()

			//log.Println(db.Content.Root.Groups[0].Groups[0].Entries[0].GetTitle())
			//log.Println(db.Content.Root.Groups[0].Groups[0].Entries[0].GetPassword())
			//log.Println(fileName)
			//log.Println(pwdEntry.Text)
			settings.New(a.Preferences()).SetLastFile(fileName.String())
		},
		w,
	)
	d.Show()
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	a = app.NewWithID("goKeeperViewer")
	w = a.NewWindow("goKeeperViewer")
	w.Resize(fyne.NewSize(640, 480))

	mainMenu := fyne.NewMainMenu(
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
	w.SetMainMenu(mainMenu)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), actionMenuOpen),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), actionSettings),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), actionHelpAbout),
	)

	passwordTree = widget.NewTree(
		func(s string) []string {
			return db.GetChildIDs(s)
		},
		func(s string) bool {
			if db == nil {
				return false
			}
			d := db.IsBranch(s)
			return d
		},
		func(b bool) fyne.CanvasObject {
			if b {
				return widget.NewLabel("")
			}
			return widget.NewLabel("")
		},
		func(s string, b bool, co fyne.CanvasObject) {
			item := db.GetItemByID(s)
			co.(*widget.Label).SetText(item.Title)
		},
	)
	passwordTree.OnSelected = func(uid widget.TreeNodeID) {
		item := db.GetItemByID(uid)
		if item.Entry == nil {
			return
		}
		passwordDetails.Show()
		//log.Println(item.Entry)
		// TODO: Process all Entry fields dinaicaly
		for _, v := range passwordDetails.Items {
			if v.Text == "Title" {
				v.Widget.(*widget.Label).SetText(item.Entry.GetTitle())
			} else if v.Text == "Password" {
				v.Widget.(*widget.Entry).SetText(item.Entry.GetTitle())
			} else if v.Text == "URL" {
				v.Widget.(*widget.Hyperlink).SetURLFromString(item.Entry.GetContent("URL"))
				v.Widget.(*widget.Hyperlink).SetText(item.Entry.GetContent("URL"))
			} else if v.Text == "UserName" {
				v.Widget.(*widget.Label).SetText(item.Entry.GetContent("UserName"))
			} else if v.Text == "Notes" {
				v.Widget.(*widget.Label).SetText(item.Entry.GetContent("Notes"))
			}
		}
	}

	passwordDetails = widget.NewForm(
		widget.NewFormItem("Title", widget.NewLabel("")),
		widget.NewFormItem("URL", widget.NewHyperlink("", nil)),
		widget.NewFormItem("UserName", widget.NewLabel("")),
		widget.NewFormItem("Password", widget.NewPasswordEntry()),
		widget.NewFormItem("Notes", widget.NewLabel("")),
	)
	passwordDetails.Hide()

	content := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		container.NewGridWithColumns(
			2,
			passwordTree,
			passwordDetails,
		),
	)

	w.SetContent(content)

	if len(os.Args) > 1 {
		loadFile(storage.NewFileURI(os.Args[1]))
	} else if settings.New(a.Preferences()).GetStartLoadOption() == settings.START_LOAD_LAST {
		var fileName = settings.New(a.Preferences()).GetLastFile()
		if fileName != "" {
			loadFile(storage.NewURI(fileName))
		}
	}

	w.ShowAndRun()

	/*
		content := container.NewBorder(
			toolbar,
			nil,
			nil,
			nil,
			widget.NewLabel("Hello!"))
	*/

}
