package loader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type FileLoader struct {
	Incoming map[int64][]int64
	incm     sync.Mutex
	Outgoing map[int64][]int64
	outm     sync.Mutex
}

func NewFileLoader(filepath string) (*FileLoader, error) {
	l := &FileLoader{
		Incoming: make(map[int64][]int64),
		Outgoing: make(map[int64][]int64),
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, l)
	if err != nil {
		return nil, err
	}

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
