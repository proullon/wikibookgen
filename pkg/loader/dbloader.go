package loader

import (
	"database/sql"

	"github.com/proullon/wikibookgen/pkg/parsing"
)

type DBLoader struct {
	db *sql.DB
}

func NewDBLoader(db *sql.DB) *DBLoader {
	l := &DBLoader{
		db: db,
	}

	return l
}

func (l *DBLoader) LoadIncomingReferences(id int64) ([]int64, error) {
	refs := make([]int64, 0)

	// WIP no index right now
	//return refs, nil

	query := `SELECT page_id FROM article_reference WHERE refered_page = $1`
	rows, err := l.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ref int64
	for rows.Next() {
		err = rows.Scan(&ref)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (l *DBLoader) LoadOutgoingReferences(id int64) ([]int64, error) {
	query := `SELECT refered_page FROM article_reference WHERE page_id = $1`
	rows, err := l.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []int64
	var ref int64
	for rows.Next() {
		err = rows.Scan(&ref)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (l *DBLoader) ID(s string) (int64, error) {
	query := `SELECT page_id FROM page WHERE lower_title = $1`

	var id int64
	err := l.db.QueryRow(query, parsing.CleanupTitle(s)).Scan(&id)
	return id, err
}
