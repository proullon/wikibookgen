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

type Louvain struct {
}

func NewLouvain() *Louvain {
	c := &Louvain{}

	return c
}

func (c *Louvain) Version() string {
	return "1"
}

/*
GraphSize:6000  Duration:13m52.7674639s
GraphSize:7000  Duration:19m10.3990492s
GraphSize:8000  Duration:27m59.5882278s
GraphSize:9000  Duration:32m18.4351384s
GraphSize:10000  Duration:38m6.0134278s
GraphSize:11000  Duration:47m4.4327723s
GraphSize:12000  Duration:55m50.8206221s
GraphSize:13000  Duration:1h6m1.4112657s
GraphSize:14000  Duration:1h14m25.1810241s
GraphSize:15000  Duration:1h25m38.1851993s
GraphSize:16000  Duration:1h38m9.7543215s
GraphSize:17000  Duration:1h51m30.4948495s
GraphSize:18000  Duration:2h6m39.3636176s
*/
func (c *Louvain) MaxSize(j Job) int64 {
	switch Model(j.Model) {
	case ABSTRACT:
		return 3500
	case TOUR:
		return 3500
	case ENCYCLOPEDIA:
		return 3500
	}

	return 3500
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb Louvain, all it does is remove vertices with low trail count
func (c *Louvain) Cluster(j Job, rootID int64, gr graph.Directed) (*Cluster, error) {

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

	var layer int
	switch Model(j.Model) {
	case ABSTRACT:
		layer = 1
	case TOUR:
		layer = 2
	case ENCYCLOPEDIA:
		layer = 3
	}

	log.Infof("Graph has %d vertices", gr.Nodes().Len())
	g, ok := gr.(*simple.DirectedGraph)
	if !ok {
		return nil, fmt.Errorf("not a simple directed graph")
	}

	/*
			// RemoveNodes with 1 degree
			nodes := g.Nodes()
			for nodes.Next() {
				n := nodes.Node()
				degree := g.From(n.ID()).Len()
				degree += g.To(n.ID()).Len()
				if degree <= 1 {
					g.RemoveNode(n.ID())
				}
			}

		log.Infof("Graph has %d vertices after removing vertices with degree 1", g.Nodes().Len())
	*/
	graphsize := g.Nodes().Len()

	// Profile calls Modularize which implements the Louvain modularization algorithm.
	// Since this is a randomized algorithm we use a defined random source to ensure
	// consistency between test runs. In practice, results will not differ greatly
	// between runs with different PRNG seeds.
	src := rand.NewSource(1)

	// Get the profile of internal node weight for resolutions
	// between 0.1 and 10 using logarithmic bisection.
	effort := 10 // instead of 10
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

	// find communities with root for each profile
	for _, d := range p {
		comm := d.Communities()
		for _, c := range comm {
			for _, n := range c {
				if n.ID() == rootID {
					fmt.Printf("Found rootID %d in community of len %d\n", rootID, len(c))
					relativesize := len(c) * 100 / graphsize
					if relativesize > 90 {
						fmt.Printf("Relative size (%d) is %d%% of total size (%d) skipping\n", len(c), relativesize, graphsize)
					}
					break
				}
			}
		}
	}

	// find root cluster, less than 90% of orignal graph
	cluster, index, err := findRootCluster(p, rootID, graphsize)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Got root cluster on index %d\n", index)

	// now fill as many layer as necessary
	populateClusters(cluster, 1, layer, index+1, p)

	return cluster, nil
}

// Find all communities including member of parent cluster
func populateClusters(parent *Cluster, depth int, wanted int, index int, profile []community.Interval) error {
	if depth == wanted {
		return nil
	}

	if index >= len(profile) {
		return fmt.Errorf("cannot use index %d for layer %d, only %d intervals in profile", index, depth, len(profile))
	}

	for _, n := range parent.Members {
		comm := findNodeInInterval(index, n.ID(), profile)
		if comm == nil {
			fmt.Printf("Could not find %d in interval %d", n.ID(), index)
		}
		cluster := &Cluster{Members: comm}
		parent.Subclusters = append(parent.Subclusters, cluster)
	}

	for _, sub := range parent.Subclusters {
		populateClusters(sub, depth+1, wanted, index+1, profile)
	}

	return nil
}

func findRootCluster(profile []community.Interval, rootID int64, graphsize int) (*Cluster, int, error) {
	for i := 0; i < len(profile); i++ {
		comm := findNodeInInterval(i, rootID, profile)
		relativesize := len(comm) * 100 / graphsize
		if relativesize > 90 {
			continue
		}
		cluster := &Cluster{Members: comm}
		return cluster, i, nil
	}

	return nil, 0, fmt.Errorf("no suitable cluster found for root cluster")
}

func findNodeInInterval(index int, nodeID int64, profile []community.Interval) []graph.Node {
	communities := profile[index].Communities()
	for _, comm := range communities {
		for _, n := range comm {
			if n.ID() == nodeID {
				return comm
			}
		}
	}

	return nil
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
