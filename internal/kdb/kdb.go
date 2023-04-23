package kdb

import (
	"fmt"
	"log"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

type KDB struct {
	db       *gokeepasslib.Database
	treeData []*KDBItem
	debug    bool
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
	logLine(s string)
}

func (k *KDB) SetDebug(d bool) {
	k.debug = d
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
	for _, v := range k.treeData {
		if v.Id == s {
			return v
		}
	}
	return nil
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
	var i KDBItem
	i.Id = "/"
	i.Title = "Root"
	i.Parent = ""
	k.logLine(fmt.Sprintf("%s - \"%s\"", i.Id, i.Title))
	k.treeData = append(k.treeData, &i)

	for key, val := range k.db.Content.Root.Groups {
		for _, v := range k.processTreeBranch(i.Id, key, val) {
			k.treeData = append(k.treeData, v)
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

func (k *KDB) Load(fileName string, filePass string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(filePass)
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
