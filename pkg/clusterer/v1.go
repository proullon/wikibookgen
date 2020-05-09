package clusterer

import (
	"fmt"
	"golang.org/x/exp/rand"
	"sort"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/simple"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
}

func NewV1() *V1 {
	c := &V1{}

	return c
}

func (c *V1) Version() string {
	return "1"
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *V1) Cluster(j Job, gr graph.Directed) (*Cluster, error) {

	/*
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
	*/

	log.Infof("Graph has %d vertices", gr.Nodes().Len())
	g, ok := gr.(*simple.DirectedGraph)
	if !ok {
		return nil, fmt.Errorf("not a simple directed graph")
	}
	// RemoveNodes with 1 degree
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		degree := g.From(n.ID()).Len()
		if degree <= 1 {
			g.RemoveNode(n.ID())
		}
	}
	log.Infof("Graph has %d vertices after removing vertices with degree 1", g.Nodes().Len())

	// Profile calls Modularize which implements the Louvain modularization algorithm.
	// Since this is a randomized algorithm we use a defined random source to ensure
	// consistency between test runs. In practice, results will not differ greatly
	// between runs with different PRNG seeds.
	src := rand.NewSource(1)

	// Get the profile of internal node weight for resolutions
	// between 0.1 and 10 using logarithmic bisection.
	effort := 2 // instead of 10
	p, err := community.Profile(community.ModularScore(g, community.Weight, effort, src), true, 1e-3, 0.1, 10)
	if err != nil {
		log.Errorf("Profiling: %s", err)
		return nil, err
	}

	// Print out each step with communities ordered.
	for _, d := range p {
		comm := d.Communities()
		for _, c := range comm {
			sort.Sort(ByID(c))
		}
		sort.Sort(BySliceIDs(comm))
		var bigcomm [][]graph.Node
		for _, c := range comm {
			if len(c) > 1 {
				bigcomm = append(bigcomm, c)
			}
		}

		fmt.Printf("Low:%.2v High:%.2v Score:%v Communities:%d Communities with more than 1 member:%d Q=%.3v\n",
			d.Low, d.High, d.Score, len(comm), len(bigcomm), community.Q(g, comm, d.Low))
	}

	return nil, nil
}

// BySliceIDs implements the sort.Interface sorting a slice of
// []graph.Node lexically by the IDs of the []graph.Node.
type BySliceIDs [][]graph.Node

func (c BySliceIDs) Len() int { return len(c) }
func (c BySliceIDs) Less(i, j int) bool {
	a, b := c[i], c[j]
	l := len(a)
	if len(b) < l {
		l = len(b)
	}
	for k, v := range a[:l] {
		if v.ID() < b[k].ID() {
			return true
		}
		if v.ID() > b[k].ID() {
			return false
		}
	}
	return len(a) < len(b)
}
func (c BySliceIDs) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// ByID implements the sort.Interface sorting a slice of graph.Node
// by ID.
type ByID []graph.Node

func (n ByID) Len() int           { return len(n) }
func (n ByID) Less(i, j int) bool { return n[i].ID() < n[j].ID() }
func (n ByID) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
