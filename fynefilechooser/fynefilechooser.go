package fynefilechooser

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var defaultTitle = "(none)"

type FileChooser struct {
	widget.Button
	uri    fyne.URI
	parent fyne.Window
	filter storage.FileFilter
}

type FileChooserIface interface {
	GetURI() fyne.URI
}

func NewFileChooser(parent fyne.Window, filter storage.FileFilter) *FileChooser {
	rv := &FileChooser{}
	rv.ExtendBaseWidget(rv)
	rv.IconPlacement = widget.ButtonIconTrailingText
	rv.SetIcon(theme.FolderOpenIcon())
	rv.SetText(defaultTitle)
	rv.parent = parent
	rv.filter = filter

	return rv
}

func (fc *FileChooser) GetURI() fyne.URI {
	return fc.uri
}

func (fc *FileChooser) Tapped(*fyne.PointEvent) {
	d := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
		if e != nil {
			dialog.NewError(e, fc.parent).Show()
			return
		}
		if u != nil {
			fc.uri = u.URI()
			fc.SetText(filepath.Base(fc.uri.String()))
		} else {
			fc.uri = nil
			fc.SetText(defaultTitle)
		}

	}, fc.parent)
	d.SetFilter(fc.filter)
	d.Show()

}
