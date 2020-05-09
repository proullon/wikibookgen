package loader

import (
	"database/sql"
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

func (l *DBLoader) LoadIncomingReferences(id int) ([]int, error) {
	refs := make([]int, 0)

	// WIP no index right now
	//return refs, nil

	query := `SELECT page_id FROM article_reference WHERE refered_page = $1`
	rows, err := l.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ref int
	for rows.Next() {
		err = rows.Scan(&ref)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (l *DBLoader) LoadOutgoingReferences(id int) ([]int, error) {
	query := `SELECT refered_page FROM article_reference WHERE page_id = $1`
	rows, err := l.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []int
	var ref int
	for rows.Next() {
		err = rows.Scan(&ref)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}
