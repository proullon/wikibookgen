package model

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/simple"
)

/*
type Wikibook struct {
	Title   string
	Subject string
	Volumes []Volume
}

type Volume struct {
	Title    string
	Chapters []Chapter
}

type Chapter struct {
	Title    string
	Articles []Page
}
*/

// Generator interface defines objects able to generate
// a Wikibook table of content given job constraints
type Generator interface {
	Generate(Job)
	Find(string, string) (int64, error)
	Complete(string, string) ([]string, error)
}

// Loader defines objects able to retrieve article references
// Production implementation use database, redis, unit test local file
type Loader interface {
	LoadIncomingReferences(int64) ([]int64, error)
	LoadOutgoingReferences(int64) ([]int64, error)
	ID(string) (int64, error)
	Title(int64) (string, error)
	Search(string) ([]string, error)
}

// Classifier interface defines objects able to select
// a coherent graph of Wikipedia articles related to
// given job constraints
type Classifier interface {
	LoadGraph(l Loader, rootID int64, maxSize int64) (graph.Directed, error)
	Version() string
}

// Clusterer interface defines objects able to group
// Wikipedia articles graph into a 1 dimension storyline (chapters)
type Clusterer interface {
	Cluster(Job, int64, graph.Directed) (*Cluster, error)
	MaxSize(Job) int64
	Version() string
}

// Orderer interface defines objects able to order
// clusters (chapters) and Wikipedia articles inside cluster
// in coherent reading order
type Orderer interface {
	Order(Loader, Job, graph.Directed, *Cluster) (*Wikibook, error)
	Version() string
}

// Editor interface defines objects able to edit graph clusters into
// a humain readable table of content
type Editor interface {
	Edit(Job, *Cluster) (*Wikibook, error)
}

type JobStatus string

const (
	CREATED JobStatus = "created"
	ONGOING JobStatus = "ongoing"
	DONE    JobStatus = "done"
)

type Job struct {
	ID       string
	Subject  string
	Model    string
	Language string
}

type Model string

const (
	ABSTRACT     Model = "abstract"
	TOUR         Model = "tour"
	ENCYCLOPEDIA Model = "encyclopedia"
)

type Node struct {
	id int64
}

func (n *Node) ID() int64 {
	return n.id
}

func (n *Node) String() string {
	return fmt.Sprintf("%d", n.id)
}

func NewNode(id int64) *Node {
	n := &Node{
		id: id,
	}

	return n
}

/*
type Graph struct {
	root     *Vertex
	vertices map[int]*Vertex
}

func NewGraph(root *Vertex, vertices map[int]*Vertex) *Graph {

	g := &Graph{
		root:     root,
		vertices: vertices,
	}
	return g
}

func (g *Graph) EdgeCount() int {
	var c int
	for _, v := range g.vertices {
		c += len(v.OutgoingEdges)
	}
	return c
}

func (g *Graph) VertexCount() int {
	return len(g.vertices)
}

func (g *Graph) Vertices() []*Vertex {
	var vertices []*Vertex

	for _, v := range g.vertices {
		vertices = append(vertices, v)
	}

	return vertices
}
*/
/*

func (g *Graph) Duplicate() graph.Graph {
	return g
}

func (g *Graph) Merge(ge graph.Edge) {
	e := ge.(*Edge)

	// create new cluster
	c := &Cluster{}

	// add all edges from 2 vertices to cluster
	for _, e := range e.src.IncomingEdges {
		c.IncomingEdges = append(c.IncomingEdges, e)
	}
	for _, e := range e.src.OutgoingEdges {
		c.OutgoingEdges = append(c.OutgoingEdges, e)
	}
	for _, e := range e.sink.IncomingEdges {
		c.IncomingEdges = append(c.IncomingEdges, e)
	}
	for _, e := range e.sink.OutgoingEdges {
		c.OutgoingEdges = append(c.OutgoingEdges, e)
	}

	// remove all reference to merged vertices in connected vertices

	// add reference to cluster in all connected vertices

}

func (g *Graph) String() string {
	return fmt.Sprintf("G(V, E) = G(%d, %d)", g.VertexCount(), g.EdgeCount())
}

func (g *Graph) RandomEdge() graph.Edge {
	// go map access is randomized by runtime on range, should do the trick
	for _, v := range g.vertices {
		if len(v.OutgoingEdges) > 0 {
			e := &Edge{
				src:  v,
				sink: v.OutgoingEdges[0],
			}
			return e
		}
	}

	return nil
}

type Edge struct {
	src  *Vertex
	sink *Vertex
}
*/

type Vertex struct {
	ID            int64
	Loaded        bool
	References    []int64
	Referers      []int64
	IncomingEdges []*Vertex `json:"-"`
	OutgoingEdges []*Vertex `json:"-"`
}

func (v *Vertex) Degree() int {
	return len(v.IncomingEdges) + len(v.OutgoingEdges)
}

/*
func (v *Vertex) Edges() []graph.Edge {
	return nil
}

func (v *Vertex) SubGraph() graph.Graph {
	return nil
}
*/

/*
// Graph wraps articles page to help generation
type Graph struct {
	TrailCount int
	Page       *Page
	Nodes      []*Graph
}

// Page holds article metadata
type Page struct {
	ID         int
	Title      string
	LowerTitle string
	References []*Reference
}

// Reference between 2 articles
type Reference struct {
	PageID      int
	ReferedPage int
	Occurence   int
	Index       int
	LowerTitle  string
}

type Cluster struct {
	Members     []*Graph
	Subclusters []*Cluster
}
*/

type Cluster struct {
	IncomingEdges []*Vertex `json:"-"`
	OutgoingEdges []*Vertex `json:"-"`

	Members     Component
	Subclusters []*Cluster
}

type TestCluster struct {
	Members     map[int64]*Node
	Subclusters []*TestCluster
}

func (c *Cluster) Depth() int {
	return c.depth(c, 1)
}

func (c *Cluster) depth(cl *Cluster, depth int) int {
	fmt.Printf("Cluster d%d : %v\n", depth, cl.Members)

	var max int = depth
	for _, sub := range cl.Subclusters {
		d := c.depth(sub, depth+1)
		if d > max {
			max = d
		}
	}

	return max
}

type Component map[int64]graph.Node

func NewComponent(nodes []graph.Node) Component {
	c := make(Component)

	for _, n := range nodes {
		c[n.ID()] = n
	}

	return c
}

func (c Component) Copy() Component {
	dup := make(map[int64]graph.Node)

	for k, v := range c {
		dup[k] = v
	}

	return dup
}

func (c Component) String() string {
	var s string

	size := len(c)
	s = "["
	var i int
	for k, _ := range c {
		s = fmt.Sprintf("%s%d", s, k)
		i++
		if i < size {
			s += ", "
		}
	}
	s += "]"

	return s
}

func (c Component) CanGrow(g graph.Directed) bool {
	minimal_degree := float64(float64(len(c)+1) / 2.0)

	for _, n := range c {
		if float64(c.Degree(g, n)) < minimal_degree {
			return false
		}
	}

	return true
}

func (c Component) Graph(g graph.Directed) graph.Directed {
	cg := simple.NewDirectedGraph()

	for k, n1 := range c {
		from := cg.Node(k)
		if from == nil {
			from = NewNode(k)
			cg.AddNode(from)
		}

		for _, n2 := range c {
			if g.HasEdgeFromTo(n1.ID(), n2.ID()) {
				to := cg.Node(n2.ID())
				if to == nil {
					to = NewNode(n2.ID())
				}
				cg.SetEdge(cg.NewEdge(from, to))
			}
		}
	}

	return cg
}

func (c Component) Betweenness(g graph.Directed) map[int64]float64 {
	// create graph containing only component
	cg := c.Graph(g)

	// compute betweenness for component only
	bvalues := network.Betweenness(cg)
	for k, v := range bvalues {
		log.Infof("%v %v", k, v)
	}

	return bvalues
}

func (c Component) Degree(g graph.Directed, n graph.Node) int {
	return g.From(n.ID()).Len() + g.To(n.ID()).Len()
}

func (c Component) Equal(other Component) bool {
	if len(c) != len(other) {
		return false
	}

	for i := range c {
		_, ok := other[i]
		if !ok {
			return false
		}
		if c[i].ID() != other[i].ID() {
			return false
		}
	}

	return true
}

func (c Component) CanJoin(g graph.Directed, n graph.Node) (int, bool) {
	minimal_degree := float64(float64(len(c)+1) / 2.0)

	if float64(g.From(n.ID()).Len()+g.To(n.ID()).Len()) < minimal_degree {
		return 0, false
	}

	var count int
	for id, _ := range c {
		if g.HasEdgeBetween(id, n.ID()) {
			count++
		}
	}

	if float64(count) >= minimal_degree {
		return count, true
	}

	return count, false
}

func (c Component) BestCandidates(g graph.Directed, candidates []graph.Node) (map[int][]graph.Node, int) {
	var maxidx int
	mm := make(map[int][]graph.Node)
	begin := time.Now()

	for _, candidate := range candidates {
		count, valid := c.CanJoin(g, candidate)
		if valid {
			mm[count] = append(mm[count], candidate)
			if count > maxidx {
				maxidx = count
			}
		}
	}

	if maxidx > 2 {
		log.Debugf("BestCandidates (%s) for cluster %s have %d edges (%d/%d)", time.Since(begin), c, maxidx, len(mm[maxidx]), len(candidates))
	}
	return mm, maxidx
}
