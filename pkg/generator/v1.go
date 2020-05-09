package generator

import (
	"database/sql"
	"encoding/json"
	//	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
	"github.com/proullon/wikibookgen/pkg/parsing"
)

type V1 struct {
	version string
	db      *sql.DB

	classifier Classifier
	clusterer  Clusterer
	orderer    Orderer
}

func NewV1(db *sql.DB, classifier Classifier, clusterer Clusterer, orderer Orderer) *V1 {
	g := &V1{
		version:    "1",
		db:         db,
		classifier: classifier,
		clusterer:  clusterer,
		orderer:    orderer,
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
	id, err := g.loadID(j.Subject)
	if err != nil {
		return err
	}

	begin := time.Now()
	graph, err := g.classifier.LoadGraph(id)
	if err != nil {
		return err
	}
	classificationDuration := time.Since(begin)
	log.Infof("%+v: %d articles", j, graph.Nodes().Len())
	//stat(graph, articles)

	/*
		data, err := json.Marshal(articles)
		if err != nil {
			return err
		}
		fmt.Printf("\n%s\n", string(data))
	*/

	begin = time.Now()
	clusters, err := g.clusterer.Cluster(j, graph)
	if err != nil {
		return err
	}
	clusteringDuration := time.Since(begin)

	begin = time.Now()
	wikibook, err := g.orderer.Order(j, graph, clusters)
	if err != nil {
		return err
	}
	orderingDuration := time.Since(begin)

	log.Infof("Generate: %+v inserting wikibook", j)

	err = g.insertWikibook(j, wikibook)
	if err != nil {
		return err
	}

	log.Infof("Generate: %+v done ! (classification: %s, clustering: %s, ordering: %s", j, classificationDuration, clusteringDuration, orderingDuration)

	return nil
}

func (g *V1) Version() string {
	return g.version // todo concat version of classifier, clusterer and orderer
}

// insertWikibook insertWikibook inserts newly generated wikibook and
// set job status to DONE
func (g *V1) insertWikibook(j Job, wikibook *Wikibook) error {

	toc, err := json.Marshal(wikibook)
	if err != nil {
		return err
	}

	tx, err := g.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO wikibook (id, subject, generator_version, gen_date, model, pages, table_of_content)
													VALUES ($1, $2, $3, NOW(), $4, $5, $6)`
	_, err = tx.Exec(query, j.ID, j.Subject, g.Version(), j.Model, wikibook.Pages, toc)
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

func (g *V1) loadID(s string) (int64, error) {
	query := `SELECT page_id FROM page WHERE lower_title = $1`

	var id int64
	err := g.db.QueryRow(query, parsing.CleanupTitle(s)).Scan(&id)
	return id, err
}

func stat(graph *Vertex, vertices map[int]*Vertex) {

	var inedge, outedge int
	for _, v := range vertices {
		inedge += len(v.IncomingEdges)
		outedge += len(v.OutgoingEdges)
	}

	log.Printf("%d vertices, %d incoming edges, %d outgoing edges", len(vertices), inedge, outedge)
}
