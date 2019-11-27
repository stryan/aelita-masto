package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type textDB struct {
	filePath string
	ids      []string
}

func CreateTextDB(path string) textDB {
	var i []string
	return textDB{filePath: path, ids: i}
}

func (t *textDB) Open() bool {
	b, err := ioutil.ReadFile(t.filePath)
	if err != nil {
		return false
	}
	t.ids = strings.Split(string(b), "\n")
	return true
}

func (t *textDB) Close() bool {
	os.Remove(t.filePath)
	save := strings.Join(t.ids, "\n")
	err := ioutil.WriteFile(t.filePath, []byte(save), 0755)
	if err != nil {
		return false
	}
	return true
}

func (t *textDB) Check(id string) bool {
	for _, i := range t.ids {
		if id == i {
			return true
		}
	}
	return false
}

func (t *textDB) Add(id string) bool {
	t.ids = append(t.ids, id)
	return true
}

func (t *textDB) Sync() bool {
	os.Remove(t.filePath)
	save := strings.Join(t.ids, "\n")
	err := ioutil.WriteFile(t.filePath, []byte(save), 0755)
	if err != nil {
		return false
	}
	return true
}
