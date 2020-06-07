package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

type LocalCacheLoader struct {
	src      Loader
	Incoming map[int64][]int64
	incm     sync.Mutex
	Outgoing map[int64][]int64
	outm     sync.Mutex
	Titles   map[int64]string
	titlem   sync.Mutex
}

func NewLocalCacheLoader(src Loader) *LocalCacheLoader {
	l := &LocalCacheLoader{
		src:      src,
		Incoming: make(map[int64][]int64),
		Outgoing: make(map[int64][]int64),
		Titles:   make(map[int64]string),
	}

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			l.Dump()
		}
	}()

	return l
}

func (l *LocalCacheLoader) Cached() int {
	l.incm.Lock()
	n := len(l.Incoming)
	l.incm.Unlock()
	return n
}

func (l *LocalCacheLoader) LoadIncomingReferences(id int64) ([]int64, error) {
	l.incm.Lock()
	refs, ok := l.Incoming[id]
	l.incm.Unlock()
	if ok {
		return refs, nil
	}

	refs, err := l.src.LoadIncomingReferences(id)
	if err != nil {
		return nil, err
	}

	l.incm.Lock()
	l.Incoming[id] = refs
	l.incm.Unlock()

	return refs, nil
}

func (l *LocalCacheLoader) LoadOutgoingReferences(id int64) ([]int64, error) {
	l.outm.Lock()
	refs, ok := l.Outgoing[id]
	l.outm.Unlock()
	if ok {
		return refs, nil
	}

	refs, err := l.src.LoadOutgoingReferences(id)
	if err != nil {
		return nil, err
	}

	l.outm.Lock()
	l.Outgoing[id] = refs
	l.outm.Unlock()

	return refs, nil
}

func (l *LocalCacheLoader) Dump() {
	l.outm.Lock()
	l.incm.Lock()
	defer l.outm.Unlock()
	defer log.Infof("LocalCache Dump created")
	defer l.incm.Unlock()

	data, err := json.Marshal(l)
	if err != nil {
		log.Errorf("Cannot marshal LocalCacheLoader: %s", err)
		return
	}
	err = ioutil.WriteFile(fmt.Sprintf("/tmp/wikibookgen/dump.%p.json", l), data, 0666)
	if err != nil {
		log.Errorf("cannot write LocalCacheLoader: %s", err)
		return
	}
}

func (l *LocalCacheLoader) ID(s string) (int64, error) {
	return l.src.ID(s)
}

func (l *LocalCacheLoader) Title(id int64) (string, error) {
	l.titlem.Lock()
	title, ok := l.Titles[id]
	l.titlem.Unlock()
	if ok {
		return title, nil
	}

	title, err := l.src.Title(id)
	if err != nil {
		return title, err
	}

	l.titlem.Lock()
	l.Titles[id] = title
	l.titlem.Unlock()

	return title, err
}
