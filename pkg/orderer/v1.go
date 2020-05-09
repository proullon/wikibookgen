package orderer

import (
	"database/sql"

	"gonum.org/v1/gonum/graph"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
	db *sql.DB
}

func NewV1(db *sql.DB) *V1 {
	o := &V1{
		db: db,
	}

	return o
}

func (o *V1) Version() string {
	return "1"
}

func (o *V1) Order(j Job, g graph.Directed, clusters *Cluster) (*Wikibook, error) {
	wikibook := &Wikibook{}
	return wikibook, nil
}
