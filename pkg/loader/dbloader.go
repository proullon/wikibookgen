package loader

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

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
	begin := time.Now()
	defer func(b time.Time) {
		log.Debugf("LoadIncomingReferences: took %s", time.Since(b))
	}(begin)
	refs := make([]int64, 0)

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
	begin := time.Now()
	defer func(b time.Time) {
		log.Debugf("LoadOutgoingReferences: took %s", time.Since(b))
	}(begin)
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
	var err error
	var id int64

	for i := 0; i < 10; i++ {
		id, err = l.id(s)
		if err == nil {
			return id, nil
		}
		time.Sleep(1 * time.Second)
	}

	return 0, err
}

func (l *DBLoader) id(s string) (int64, error) {
	query := `SELECT page_id FROM page WHERE lower_title = $1`

	var id int64
	err := l.db.QueryRow(query, parsing.CleanupTitle(s)).Scan(&id)
	return id, err
}

func (l *DBLoader) Title(id int64) (string, error) {
	var err error
	var t string

	for i := 0; i < 10; i++ {
		t, err = l.title(id)
		if err == nil {
			return t, nil
		}
		log.Errorf("Title(%d): %s", id, err)
		time.Sleep(1 * time.Second)
	}

	return "", err
}

func (l *DBLoader) title(id int64) (string, error) {
	query := `SELECT title FROM page WHERE page_id = $1`

	var title string
	err := l.db.QueryRow(query, id).Scan(&title)
	return title, err
}

func (l *DBLoader) Search(value string) ([]string, error) {
	var titles []string

	if len(value) < 3 {
		return titles, nil
	}

	query := `SELECT title FROM page WHERE lower_title LIKE $1 LIMIT 100`
	rows, err := l.db.Query(query, fmt.Sprintf("%s%%", value))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}

		titles = append(titles, t)
	}

	return titles, nil
}

func (l *DBLoader) Content(id int64) (string, error) {
	var c string
	var err error

	for i := 0; i < 10; i++ {
		query := `SELECT content FROM page_content WHERE page_id = $1`
		err := l.db.QueryRow(query, id).Scan(&c)
		if err != nil && err == sql.ErrNoRows {
			return "", err
		}
		if err == nil {
			return c, nil
		}
	}

	return "", err
}
