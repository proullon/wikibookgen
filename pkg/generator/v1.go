package generator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
	version string
	db      *sql.DB
	workdir string

	loaders    map[string]Loader
	classifier Classifier
	clusterer  Clusterer
	orderer    Orderer
	editor     Editor
}

func NewV1(db *sql.DB, classifier Classifier, clusterer Clusterer, orderer Orderer, loaders map[string]Loader, editor Editor, workdir string) *V1 {
	g := &V1{
		version:    "1",
		db:         db,
		workdir:    workdir,
		classifier: classifier,
		clusterer:  clusterer,
		orderer:    orderer,
		loaders:    loaders,
		editor:     editor,
	}

	return g
}

// Generate table of content for given job
//
// Run Classifier, Clusterer and Orderer then save result in database.
// If an error is encountered, set job back to Created
func (g *V1) Generate(j Job) {
	err := g.generate(j)
	if err == nil {
		return
	}
	log.Errorf("Generate: generation failed: %s", err)

	query := `UPDATE job SET status = $1 WHERE id = $2`
	_, err = g.db.Exec(query, CREATED, j.ID)
	if err != nil {
		log.Errorf("Generate: cannot reset job %s", j.ID)
	}
}

func (g *V1) generate(j Job) error {
	loader, ok := g.loaders[j.Language]
	if !ok {
		return fmt.Errorf("Loader for lang=%s not available", j.Language)
	}

	id, err := g.Find(j.Subject, j.Language)
	if err != nil {
		return err
	}

	begin := time.Now()
	graph, err := g.classifier.LoadGraph(loader, id, g.clusterer.MaxSize(j))
	if err != nil {
		return err
	}
	classificationDuration := time.Since(begin)
	log.Infof("%+v: %d articles", j, graph.Nodes().Len())

	begin = time.Now()
	clusters, err := g.clusterer.Cluster(j, id, graph)
	if err != nil {
		return err
	}
	clusteringDuration := time.Since(begin)

	begin = time.Now()
	wikibook, err := g.orderer.Order(loader, j, graph, clusters)
	if err != nil {
		return err
	}
	orderingDuration := time.Since(begin)

	log.Infof("Generate: %+v inserting wikibook", j)

	begin = time.Now()
	err = g.editor.Edit(loader, j, wikibook)
	if err != nil {
		return err
	}
	editingDuration := time.Since(begin)

	for i := 0; i < 3; i++ {
		err = g.insertWikibook(j, wikibook)
		if err == nil {
			break
		}
		if i == 2 {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	begin = time.Now()
	g.editor.Print(loader, wikibook, g.workdir)
	if err != nil {
		log.Errorf(`Printing failed for job %+v: %s`, j, err)
	}
	printingDuration := time.Since(begin)

	log.Infof("Generate: %+v done ! (classification: %s, clustering: %s, ordering: %s, editing: %s, printing: %s", j, classificationDuration, clusteringDuration, orderingDuration, editingDuration, printingDuration)

	return nil
}

func (g *V1) Version() string {
	return fmt.Sprintf("%s-%s-%s", g.version, g.classifier.Version(), g.clusterer.Version())
}

// insertWikibook insertWikibook inserts newly generated wikibook and
// set job status to DONE
func (g *V1) insertWikibook(j Job, wikibook *Wikibook) error {

	toc, err := json.Marshal(wikibook)
	if err != nil {
		return err
	}

	log.Infof("Inserting wikibook:\n%s", toc)

	tx, err := g.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO wikibook (id, subject, generator_version, gen_date, model, pages, language, table_of_content)
													VALUES ($1, $2, $3, NOW(), $4, $5, $6, $7)`
	_, err = tx.Exec(query, j.ID, j.Subject, g.Version(), j.Model, wikibook.Pages, j.Language, toc)
	if err != nil {
		return err
	}

	query = `UPDATE job SET book_id = $1, status = $2 WHERE id = $3`
	_, err = tx.Exec(query, j.ID, DONE, j.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (g *V1) Find(s string, lang string) (int64, error) {
	l, ok := g.loaders[lang]
	if !ok {
		return 0, fmt.Errorf("no loader available for lang '%s'", lang)
	}

	return l.ID(s)
}

func (g *V1) Complete(value string, language string) ([]string, error) {

	l, ok := g.loaders[language]
	if !ok {
		return nil, fmt.Errorf("no loader available for lang '%s'", language)
	}

	return l.Search(value)
}

func (g *V1) Print(w *Wikibook) error {

	loader, ok := g.loaders[w.Language]
	if !ok {
		return fmt.Errorf("Loader for lang=%s not available", w.Language)
	}

	err := g.editor.Print(loader, w, g.workdir)
	if err != nil {
		return err
	}

	return nil
}

func (g *V1) Open(id string, format string) (io.ReadCloser, error) {

	p := path.Join(g.workdir, fmt.Sprintf("%s.%s", id, format))
	reader, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
