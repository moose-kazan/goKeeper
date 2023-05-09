package kdb

import (
	"errors"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/tobischo/gokeepasslib/v3"
)

type KDB struct {
	db           *gokeepasslib.Database
	treeData     []*KDBItem
	treeDataById map[string]*KDBItem
	debug        bool
}

type KDBItem struct {
	Id     string
	Title  string
	Parent string
	Entry  *gokeepasslib.Entry
}

type KDBIface interface {
	Tree() []*KDBItem
	Load(fileName string, filePass string) error
	GetChildIDs(s string) []string
	GetItemByID(s string) *KDBItem
	IsBranch(s string) bool
	SetDebug(d bool)
	processTreeBranch(id string, key int, val gokeepasslib.Group) []*KDBItem
	processTreeItem(id string, key int, val gokeepasslib.Entry) []*KDBItem
	getCredentials(filePass string, keyFileName fyne.URI) (*gokeepasslib.DBCredentials, error)
	logLine(s string)
}

func (k *KDB) SetDebug(d bool) {
	k.debug = d
}

func (k *KDB) getCredentials(filePass string, keyFileName fyne.URI) (*gokeepasslib.DBCredentials, error) {
	var rv *gokeepasslib.DBCredentials = nil
	var err error = nil
	if filePass != "" && keyFileName != nil {
		rv, err = gokeepasslib.NewPasswordAndKeyCredentials(filePass, keyFileName.Path())
	} else if filePass != "" {
		rv = gokeepasslib.NewPasswordCredentials(filePass)
	} else if keyFileName != nil {
		rv, err = gokeepasslib.NewKeyCredentials(keyFileName.Path())
	} else {
		err = errors.New("No credentials provided!")
	}

	return rv, err
}

func (k *KDB) GetChildIDs(s string) []string {
	rv := make([]string, 0)
	for _, v := range k.treeData {
		if v.Parent == s {
			rv = append(rv, v.Id)
		}
	}
	return rv
}

func (k *KDB) GetItemByID(s string) *KDBItem {
	return k.treeDataById[s]
}

func (k *KDB) IsBranch(s string) bool {
	k.logLine("KDB IsBranch")
	if k.treeData != nil {
		k.logLine("KDB IsBranch Not Nil")
		for _, v := range k.treeData {
			if v.Id == s {
				k.logLine("KDB IsBranch Found")
				return v.Entry == nil
			}
		}
	}
	return false
}

func (k *KDB) logLine(s string) {
	if k.debug {
		log.Println(s)
	}
}

func (k *KDB) Tree() []*KDBItem {
	if k.treeData != nil {
		return k.treeData
	}
	k.treeData = make([]*KDBItem, 0)
	k.treeDataById = make(map[string]*KDBItem)

	var i KDBItem
	i.Id = "/"
	i.Title = "Root"
	i.Parent = ""
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	k.treeData = append(k.treeData, &i)
	k.treeDataById[i.Id] = &i

	for key, val := range k.db.Content.Root.Groups {
		for _, v := range k.processTreeBranch(i.Id, key, val) {
			k.treeData = append(k.treeData, v)
			k.treeDataById[v.Id] = v
		}
	}

	for i := 1; i < len(k.treeData); i++ {
		if k.treeData[i].Title < k.treeData[i-1].Title {
			for j := i; j > 0 && k.treeData[j].Title < k.treeData[j-1].Title; j-- {
				k.treeData[j], k.treeData[j-1] = k.treeData[j-1], k.treeData[j]
			}
		}
	}
	return k.treeData
}

func (k *KDB) processTreeBranch(id string, key int, val gokeepasslib.Group) []*KDBItem {
	rv := make([]*KDBItem, 0)
	var i KDBItem
	i.Id = fmt.Sprintf("%sG%d/", id, key)
	i.Title = val.Name
	i.Parent = id
	rv = append(rv, &i)
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	for key, val := range val.Groups {
		for _, v := range k.processTreeBranch(i.Id, key, val) {
			rv = append(rv, v)
		}
	}
	for key, val := range val.Entries {
		for _, v := range k.processTreeItem(i.Id, key, val) {
			rv = append(rv, v)
		}
	}
	return rv
}

func (k *KDB) processTreeItem(id string, key int, val gokeepasslib.Entry) []*KDBItem {
	rv := make([]*KDBItem, 0)
	var i KDBItem
	i.Id = fmt.Sprintf("%sI%d/", id, key)
	i.Title = val.GetTitle()
	i.Entry = &val
	i.Parent = id
	rv = append(rv, &i)
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	return rv
}

func (k *KDB) Load(fileName fyne.URI, filePass string, keyFileName fyne.URI) error {
	file, err := os.Open(fileName.Path())
	if err != nil {
		return err
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials, err = k.getCredentials(filePass, keyFileName)

	if err != nil {
		return err
	}

	err = gokeepasslib.NewDecoder(file).Decode(db)

	if err != nil {
		return err
	}
	db.UnlockProtectedEntries()

	k.db = db

	return nil
}

func New() *KDB {
	var k KDB
	k.debug = false
	return &k
}
