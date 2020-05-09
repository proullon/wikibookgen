package cache

import (
	"sync"

	. "github.com/proullon/wikibookgen/api/model"
)

type LocalCacheLoader struct {
	src      Loader
	incoming map[int][]int
	incm     sync.Mutex
	outgoing map[int][]int
	outm     sync.Mutex
}

func NewLocalCacheLoader(src Loader) *LocalCacheLoader {
	l := &LocalCacheLoader{
		src:      src,
		incoming: make(map[int][]int),
		outgoing: make(map[int][]int),
	}

	return l
}

func (l *LocalCacheLoader) Cached() int {
	l.incm.Lock()
	n := len(l.incoming)
	l.incm.Unlock()
	return n
}

func (l *LocalCacheLoader) LoadIncomingReferences(id int) ([]int, error) {
	l.incm.Lock()
	refs, ok := l.incoming[id]
	l.incm.Unlock()
	if ok {
		return refs, nil
	}

	refs, err := l.src.LoadIncomingReferences(id)
	if err != nil {
		return nil, err
	}

	l.incm.Lock()
	l.incoming[id] = refs
	l.incm.Unlock()

	return refs, nil
}

func (l *LocalCacheLoader) LoadOutgoingReferences(id int) ([]int, error) {
	l.outm.Lock()
	refs, ok := l.outgoing[id]
	l.outm.Unlock()
	if ok {
		return refs, nil
	}

	refs, err := l.src.LoadOutgoingReferences(id)
	if err != nil {
		return nil, err
	}

	l.outm.Lock()
	l.outgoing[id] = refs
	l.outm.Unlock()

	return refs, nil
}

/*
func GetPage(db *sql.DB, lt string) (*Page, error) {
	indexm.Lock()
	if pageIndex == nil {
		pageIndex = make(map[string]*Page)
	}
	p, ok := pageIndex[lt]
	indexm.Unlock()
	if ok {
		hit++
		return p, nil
	}

	p = &Page{LowerTitle: lt}
	query := `SELECT page.page_id, ar.refered_page, ar.occurrence, ar.reference_index, p.lower_title
						FROM page
						JOIN article_reference AS ar ON page.page_id = ar.page_id
						JOIN page AS p ON p.page_id = ar.refered_page
						WHERE page.lower_title = $1`
	rows, err := db.Query(query, lt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := &Reference{}
		err = rows.Scan(&r.PageID, &r.ReferedPage, &r.Occurence, &r.Index, &r.LowerTitle)
		if err != nil {
			return nil, err
		}
		p.ID = r.PageID
		p.References = append(p.References, r)
	}

	indexm.Lock()
	pageIndex[lt] = p
	indexm.Unlock()
	return p, nil
}
*/
