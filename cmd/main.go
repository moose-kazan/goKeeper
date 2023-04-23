package main

import (
	"gokeeper/internal/kdb"
	"log"
	"os"
	"path/filepath"

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
	dialog.ShowInformation("About goKeeper", "Simple viewer for KDBX (KeePass) files.", w)
}

func actionMenuOpen() {
	d := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
		if e != nil {
			dialog.NewError(e, w).Show()
			return
		}
		if u != nil {
			loadFile(u.URI().Path())
		}

	}, w)
	d.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	d.Show()
}

func loadFile(fileName string) {
	pwdEntry := widget.NewPasswordEntry()
	dialog.NewForm(
		"Enter password",
		"OK",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("File Name", widget.NewLabel(filepath.Base(fileName))),
			widget.NewFormItem("Password", pwdEntry),
		},
		func(b bool) {
			if !b {
				return
			}

			tmpDb := kdb.New()
			err := tmpDb.Load(fileName, pwdEntry.Text)

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
		},
		w,
	).Show()
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	a = app.New()
	w = a.NewWindow("goKeeper")
	w.Resize(fyne.NewSize(640, 480))

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu(
			"File",
			newMenuItem("Open", actionMenuOpen, theme.DocumentIcon(), nil),
			newMenuItem("Quit", func() { a.Quit() }, theme.LogoutIcon(), nil),
		),
		fyne.NewMenu(
			"Help",
			newMenuItem("About", actionHelpAbout, theme.InfoIcon(), nil),
		),
	)
	w.SetMainMenu(mainMenu)

	/*
		toolbar := widget.NewToolbar(
			widget.NewToolbarAction(theme.DocumentIcon(), actionMenuOpen),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.HelpIcon(), actionHelpAbout),
		)
	*/

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
		for _, v := range passwordDetails.Items {
			if v.Text == "Title" {
				v.Widget.(*widget.Label).SetText(item.Entry.GetTitle())
			} else if v.Text == "Password" {
				v.Widget.(*widget.Entry).SetText(item.Entry.GetTitle())
			} else if v.Text == "URL" {
				log.Println(item.Entry.Values)
			}
		}
	}

	passwordDetails = widget.NewForm(
		widget.NewFormItem("Title", widget.NewLabel("")),
		widget.NewFormItem("URL", widget.NewLabel("")),
		widget.NewFormItem("Password", widget.NewPasswordEntry()),
	)
	passwordDetails.Hide()

	content := container.NewGridWithColumns(
		2,
		passwordTree,
		passwordDetails,
	)

	w.SetContent(content)
	w.ShowAndRun()

	/*
		if len(os.Args) > 1 {
			loadFile(os.Args[1])
		}
	*/

	/*
		content := container.NewBorder(
			toolbar,
			nil,
			nil,
			nil,
			widget.NewLabel("Hello!"))
	*/

}
