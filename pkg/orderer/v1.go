package orderer

import (
	"database/sql"

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

func (o *V1) Order(j Job, root *Vertex, vertices map[int]*Vertex, clusters *Cluster) (*Wikibook, error) {
	wikibook := &Wikibook{}
	return wikibook, nil
}
