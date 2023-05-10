package kdb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
	Id       string
	Title    string
	Parent   string
	Entry    *gokeepasslib.Entry
	ChildIds []string
}

type KDBIface interface {
	Tree() []*KDBItem
	Load(fileName string, filePass string) error
	GetChildIDs(s string) []string
	GetItemByID(s string) *KDBItem
	IsBranch(s string) bool
	SetDebug(d bool)
	processTreeBranch(parent *KDBItem, key int, val gokeepasslib.Group) []*KDBItem
	processTreeItem(parent *KDBItem, key int, val gokeepasslib.Entry) []*KDBItem
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
	k.logLine(fmt.Sprintf("GetChildIDs(%s)", s))
	v := k.GetItemByID(s)
	if v == nil {
		return make([]string, 0)
	}
	return v.ChildIds
}

func (k *KDB) GetItemByID(s string) *KDBItem {
	k.logLine(fmt.Sprintf("GetItemById(%s)", s))
	if k.treeDataById == nil {
		return nil
	}
	return k.treeDataById[s]
}

func (k *KDB) IsBranch(s string) bool {
	k.logLine(fmt.Sprintf("IsBranch(%s)", s))
	v := k.GetItemByID(s)
	if v == nil {
		return false
	}
	return v.Entry == nil
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
	i.ChildIds = make([]string, 0)
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	k.treeData = append(k.treeData, &i)
	k.treeDataById[i.Id] = &i

	for key, val := range k.db.Content.Root.Groups {
		for _, v := range k.processTreeBranch(&i, key, val) {
			k.treeData = append(k.treeData, v)
			k.treeDataById[v.Id] = v
		}
	}

	// Sort alphabeticlay!
	for i := 1; i < len(k.treeData); i++ {
		// Sort without register
		if strings.ToUpper(k.treeData[i].Title) < strings.ToUpper(k.treeData[i-1].Title) {
			for j := i; j > 0 && k.treeData[j].Title < k.treeData[j-1].Title; j-- {
				k.treeData[j], k.treeData[j-1] = k.treeData[j-1], k.treeData[j]
			}
		}
	}

	// Only on sorted data!
	// In other cases childs ids wil not be sorted alphabeticlay!
	for _, v := range k.treeData {
		if v.Id != "/" {
			k.GetItemByID(v.Parent).ChildIds = append(
				k.GetItemByID(v.Parent).ChildIds,
				v.Id,
			)
		}
	}
	return k.treeData
}

func (k *KDB) processTreeBranch(parent *KDBItem, key int, val gokeepasslib.Group) []*KDBItem {
	rv := make([]*KDBItem, 0)
	var i KDBItem
	i.Id = fmt.Sprintf("%sG%d/", parent.Id, key)
	i.Title = val.Name
	i.Parent = parent.Id
	i.ChildIds = make([]string, 0)
	//parent.ChildIds = append(parent.ChildIds, i.Id)
	rv = append(rv, &i)
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	for key, val := range val.Groups {
		for _, v := range k.processTreeBranch(&i, key, val) {
			rv = append(rv, v)
		}
	}
	for key, val := range val.Entries {
		for _, v := range k.processTreeItem(&i, key, val) {
			rv = append(rv, v)
		}
	}
	return rv
}

func (k *KDB) processTreeItem(parent *KDBItem, key int, val gokeepasslib.Entry) []*KDBItem {
	rv := make([]*KDBItem, 0)
	var i KDBItem
	i.Id = fmt.Sprintf("%sI%d/", parent.Id, key)
	i.Title = val.GetTitle()
	i.Entry = &val
	i.Parent = parent.Id
	//parent.ChildIds = append(parent.ChildIds, i.Id)
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
