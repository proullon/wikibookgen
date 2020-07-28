package wikibookgen

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	query := `SELECT id, subject, model, language FROM job WHERE status = $1 LIMIT 1`
	err = tx.QueryRow(query, CREATED).Scan(&j.ID, &j.Subject, &j.Model, &j.Language)
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

	log.Infof("Starting generation job for %v", j)
	wg.gen.Generate(j)
}

func (wg *WikibookGen) QueueGenerationJob(subject, model, lang string) (string, error) {

	err := wg.ValidateLanguage(lang)
	if err != nil {
		return "", err
	}

	err = wg.ValidateSubject(subject, lang)
	if err != nil {
		return "", err
	}

	err = wg.ValidateModel(model)
	if err != nil {
		return "", err
	}

	if uuid, err := wg.LoadWikibook(subject, model, lang); err == nil {
		return uuid, nil
	}

	query := `INSERT INTO job (subject, model, creation_date, status, language) VALUES ($1, $2, NOW(), $3, $4) RETURNING id`

	var id string
	err = wg.db.QueryRow(query, subject, model, CREATED, lang).Scan(&id)
	return id, err
}

func (wg *WikibookGen) LoadWikibook(subject, model, lang string) (string, error) {
	query := `SELECT id FROM wikibook WHERE subject = $1 AND model = $2 AND lang = $3`

	var id string
	err := wg.db.QueryRow(query, subject, model, lang).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
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
		return fmt.Errorf("ValidateSubject(%s, %s): %s", s, lang, err)
	}
	return nil
}

func (wg *WikibookGen) ValidateLanguage(language string) error {
	switch language {
	case "fr", "en":
		return nil
	default:
		return fmt.Errorf("language not available")
	}
}

func (wg *WikibookGen) LoadOrder(uuid string) (string, string, error) {
	query := `SELECT status, book_id FROM job WHERE id = $1`

	var status, bookid string
	var bookID sql.NullString
	err := wg.db.QueryRow(query, uuid).Scan(&status, &bookID)
	if err != nil {
		return ``, ``, nil
	}

	if bookID.Valid {
		bookid = bookID.String
	}

	return status, bookid, err
}

func (wg *WikibookGen) Load(uuid string) (*Wikibook, error) {
	log.Infof("Loading %s", uuid)
	query := `SELECT subject, model, language, pages, table_of_content, gen_date FROM wikibook WHERE id = $1`

	w := &Wikibook{}
	var subject, model, language, toc, gendate string
	var pages int64
	err := wg.db.QueryRow(query, uuid).Scan(&subject, &model, &language, &pages, &toc, &gendate)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(toc), w)
	if err != nil {
		return nil, err
	}

	w.Uuid = uuid
	w.Subject = subject
	w.Model = model
	w.Language = language
	w.Pages = pages
	w.GenerationDate = gendate

	return w, nil
}

func (wg *WikibookGen) List(page int64, size int64, language string) ([]*Wikibook, error) {

	if size > 200 {
		size = 200
	}
	if page < 1 {
		page = 1
	}

	limit := size
	offset := (page - 1) * size

	query := `SELECT id, subject, model, language, pages, gen_date
	FROM wikibook ORDER BY gen_date DESC OFFSET $1 LIMIT $2`

	rows, err := wg.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Wikibook
	for rows.Next() {
		b := &Wikibook{}

		err = rows.Scan(&b.Uuid, &b.Subject, &b.Model, &b.Language, &b.Pages, &b.GenerationDate)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (wg *WikibookGen) Complete(value string, language string) ([]string, error) {
	begin := time.Now()

	titles, err := wg.gen.Complete(value, language)
	log.Infof("Complete(%s, %s): %s", value, language, time.Since(begin))
	if err != nil {
		return nil, err
	}
	return titles, nil
}

func (wg *WikibookGen) Download(id string, format string, w http.ResponseWriter) error {
	log.Infof("Download request for %s.%s", id, format)

	_, err := wg.Load(id)
	if err != nil {
		return err
	}

	reader, err := wg.gen.Open(id, format)
	if err != nil {
		return err
	}
	defer reader.Close()

	// TODO: increase download count

	switch format {
	case "epub":
		w.Header().Set("Content-Type", "application/epub+zip")
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}

	_, err = io.Copy(w, reader)
	if err != nil {
		return err
	}

	return nil
}

func (wg *WikibookGen) AvailableFormat(id string) (epub bool, pdf bool, err error) {
	var reader io.ReadCloser

	_, err = wg.Load(id)
	if err != nil {
		return false, false, err
	}

	reader, err = wg.gen.Open(id, `epub`)
	if err == nil {
		epub = true
		reader.Close()
	}

	reader, err = wg.gen.Open(id, `pdf`)
	if err == nil {
		pdf = true
		reader.Close()
	}

	return epub, pdf, nil
}

func (wg *WikibookGen) PrintWikibook(id string, format string) error {
	log.Infof("Print request for %s.%s", id, format)

	w, err := wg.Load(id)
	if err != nil {
		return err
	}

	err = wg.gen.Print(w)
	if err != nil {
		return err
	}

	return nil
}
