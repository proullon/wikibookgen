package wikibookgen

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

type WikibookGen struct {
	db  *sql.DB
	gen Generator
}

func New(db *sql.DB, gen Generator) *WikibookGen {
	wg := &WikibookGen{
		db:  db,
		gen: gen,
	}

	wg.startJobRoutine()
	return wg
}

func (wg *WikibookGen) startJobRoutine() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			wg.jobRoutine()
		}
	}()
}

// job routine looks for new generation job
// then start a generation and update job table
func (wg *WikibookGen) jobRoutine() {
	tx, err := wg.db.Begin()
	if err != nil {
		log.Errorf("job: cannot start tx: %s", err)
		return
	}
	defer tx.Rollback()

	var j Job
	query := `SELECT id, subject, model FROM job WHERE status = $1 LIMIT 1`
	err = tx.QueryRow(query, CREATED).Scan(&j.ID, &j.Subject, &j.Model)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("job: cannot query waiting job: %s", err)
		}
		return
	}

	query = `UPDATE job SET status = $1 WHERE id = $2`
	_, err = tx.Exec(query, ONGOING, j.ID)
	if err != nil {
		log.Errorf("job: cannot set job %s as %s: %s", j.ID, ONGOING, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("job: cannot commit tx for job %s: %s", j.ID, err)
		return
	}

	log.Infof("coucou start %v", j)
	wg.gen.Generate(j)
}

func (wg *WikibookGen) QueueGenerationJob(subject, model, lang string) (string, error) {

	err := wg.ValidateSubject(subject, lang)
	if err != nil {
		return "", err
	}

	err = wg.ValidateModel(model)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO job (subject, model, creation_date, status) VALUES ($1, $2, NOW(), $3) RETURNING id`

	var id string
	err = wg.db.QueryRow(query, subject, model, CREATED).Scan(&id)
	return id, err
}

func (wg *WikibookGen) ValidateModel(s string) error {

	ps := Model(s)
	switch ps {
	case ABSTRACT, TOUR, ENCYCLOPEDIA:
		return nil
	default:
		return fmt.Errorf("invalid model '%s'", s)
	}
}

func (wg *WikibookGen) ValidateSubject(s string, lang string) error {
	_, err := wg.gen.Find(s, lang)
	if err != nil {
		return err
	}
	return nil
}

func (wg *WikibookGen) LoadOrder(uuid string) (string, string, error) {
	query := `SELECT status, book_id FROM job WHERE id = $1`

	var status, bookID string
	err := wg.db.QueryRow(query, uuid).Scan(&status, &bookID)

	return status, bookID, err
}

func (wg *WikibookGen) Load(uuid string) (*Wikibook, error) {
	return nil, nil
}

func (wg *WikibookGen) List() ([]*Wikibook, error) {
	return nil, nil
}
