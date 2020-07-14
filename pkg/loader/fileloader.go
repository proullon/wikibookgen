package loader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	log "github.com/sirupsen/logrus"
)

type FileLoader struct {
	Incoming map[int64][]int64
	incm     sync.Mutex
	Outgoing map[int64][]int64
	outm     sync.Mutex
	Titles   map[int64]string
	titlem   sync.Mutex
	Contents map[int64]string
	contentm sync.Mutex
}

func NewFileLoader(filepath string) (*FileLoader, error) {
	l := &FileLoader{
		Incoming: make(map[int64][]int64),
		Outgoing: make(map[int64][]int64),
		Titles:   make(map[int64]string),
		Contents: make(map[int64]string),
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, l)
	if err != nil {
		return nil, err
	}

	log.Infof("NewFileLoader: %d contents available", len(l.Contents))
	return l, nil
}

func (l *FileLoader) LoadIncomingReferences(id int64) ([]int64, error) {
	l.incm.Lock()
	refs, ok := l.Incoming[id]
	l.incm.Unlock()
	if ok {
		return refs, nil
	}

	return nil, sql.ErrNoRows
}

func (l *FileLoader) LoadOutgoingReferences(id int64) ([]int64, error) {
	l.outm.Lock()
	refs, ok := l.Outgoing[id]
	l.outm.Unlock()
	if ok {
		return refs, nil
	}

	return nil, sql.ErrNoRows
}

func (l *FileLoader) ID(s string) (int64, error) {
	return 0, fmt.Errorf("FileLoader.ID is not implemented")
}

func (l *FileLoader) Title(id int64) (string, error) {
	l.titlem.Lock()
	title, ok := l.Titles[id]
	l.titlem.Unlock()
	if ok {
		return title, nil
	}

	return "", fmt.Errorf("not found")
}

func (l *FileLoader) Search(value string) ([]string, error) {
	titles := make([]string, 0)
	return titles, nil
}

func (l *FileLoader) Content(id int64) (string, error) {
	l.contentm.Lock()
	content, ok := l.Contents[id]
	l.contentm.Unlock()
	if ok {
		return content, nil
	}

	return "", fmt.Errorf("not found")
}
