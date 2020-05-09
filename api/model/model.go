package model

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
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
}

// Loader defines objects able to retrieve article references
// Production implementation use database, redis, unit test local file
type Loader interface {
	LoadIncomingReferences(int64) ([]int64, error)
	LoadOutgoingReferences(int64) ([]int64, error)
}

// Classifier interface defines objects able to select
// a coherent graph of Wikipedia articles related to
// given job constraints
type Classifier interface {
	LoadGraph(int64) (graph.Directed, error)
	Version() string
}

// Clusterer interface defines objects able to group
// Wikipedia articles graph into a 1 dimension storyline (chapters)
type Clusterer interface {
	Cluster(Job, graph.Directed) (*Cluster, error)
	Version() string
}

// Orderer interface defines objects able to order
// clusters (chapters) and Wikipedia articles inside cluster
// in coherent reading order
type Orderer interface {
	Order(Job, graph.Directed, *Cluster) (*Wikibook, error)
	Version() string
}

type JobStatus string

const (
	CREATED JobStatus = "created"
	ONGOING JobStatus = "ongoing"
	DONE    JobStatus = "done"
)

type Job struct {
	ID      string
	Subject string
	Model   string
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

	Members     []*Vertex
	Subclusters []*Cluster
}
