package settings

import "fyne.io/fyne/v2"

var (
	startLoadVariants = []string{
		"None",
		"Last File",
	}
)

const (
	START_LOAD_NONE = iota
	START_LOAD_LAST
)

type goKeeperSettings struct {
	pref fyne.Preferences
}

type goKeeperSettingsIface interface {
	GetLastFile() string
	SetLastFile(fileName string) *goKeeperSettings
	GetStartLoadOption() int
	SetStartLoadOption(s string) *goKeeperSettings
}

func StartLoadOptions() []string {
	return startLoadVariants
}

func New(pref fyne.Preferences) *goKeeperSettings {
	var rv goKeeperSettings
	rv.pref = pref
	return &rv
}

func (p *goKeeperSettings) GetLastFile() string {
	return p.pref.String("lastFile")
}

func (p *goKeeperSettings) SetLastFile(fileName string) *goKeeperSettings {
	p.pref.SetString("lastFile", fileName)
	return p
}

func (p *goKeeperSettings) GetStartLoadOption() int {
	return p.pref.Int("loadOnStart")
}

func (p *goKeeperSettings) SetStartLoadOption(s string) *goKeeperSettings {
	for k, v := range startLoadVariants {
		if v == s {
			p.pref.SetInt("loadOnStart", k)
			break
		}
	}
	return p
}
