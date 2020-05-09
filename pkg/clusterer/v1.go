package clusterer

import (
	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
}

func NewV1() *V1 {
	c := &V1{}

	return c
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *V1) Cluster(j Job, root *Vertex, vertices map[int]*Vertex) (*Cluster, error) {
	var list []*Vertex

	for _, v := range vertices {
		if v.Degree() > 2 {
			list = append(list, v)
		}
	}

	var modified bool
	for {
		modified = false
		for i := range list {
			if i > 0 && list[i].Degree() > list[i-1].Degree() {
				tmp := list[i]
				list[i] = list[i-1]
				list[i-1] = tmp
				modified = true
			}
		}
		if !modified {
			break
		}
	}

	var maxpage int
	switch Model(j.Model) {
	case ABSTRACT:
		maxpage = 100
	case TOUR:
		maxpage = 1000
	case ENCYCLOPEDIA:
		maxpage = 10000
	}

	cluster := &Cluster{}
	for i, _ := range list {
		if i == maxpage {
			break
		}
		//log.Infof("%-15d: %d edges", g.ID, g.Edges())
		cluster.Members = append(cluster.Members, list[i])
	}

	//g := NewGraph(root, vertices)
	//clusters := graph.HCS(g)
	log.Infof("Got %d members", len(cluster.Members))
	return cluster, nil
}
