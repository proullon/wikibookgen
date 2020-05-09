package loader

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"sync"

	. "github.com/proullon/wikibookgen/api/model"
)

type FileLoader struct {
	incoming map[int64][]int64
	incm     sync.Mutex
	outgoing map[int64][]int64
	outm     sync.Mutex
}

func NewFileLoader(filepath string) (*FileLoader, error) {
	l := &FileLoader{
		incoming: make(map[int64][]int64),
		outgoing: make(map[int64][]int64),
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	vertices := make(map[int64]*Vertex)
	err = json.Unmarshal(content, &vertices)
	if err != nil {
		return nil, err
	}

	for k, v := range vertices {
		l.incoming[k] = v.Referers
		l.outgoing[k] = v.References
	}

	return l, nil
}

func (l *FileLoader) LoadIncomingReferences(id int64) ([]int64, error) {
	l.incm.Lock()
	refs, ok := l.incoming[id]
	l.incm.Unlock()
	if ok {
		return refs, nil
	}

	return nil, sql.ErrNoRows
}

func (l *FileLoader) LoadOutgoingReferences(id int64) ([]int64, error) {
	l.outm.Lock()
	refs, ok := l.outgoing[id]
	l.outm.Unlock()
	if ok {
		return refs, nil
	}

	return nil, sql.ErrNoRows
}
