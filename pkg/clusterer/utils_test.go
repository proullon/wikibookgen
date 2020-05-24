package clusterer

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"

	. "github.com/proullon/wikibookgen/api/model"
	"github.com/proullon/wikibookgen/pkg/classifier"
	"github.com/proullon/wikibookgen/pkg/loader"
)

var (
	MathPageID              int64 = 3697062
	GraphPageID             int64 = 818568
	GraphMathPageID         int64 = 2642295
	GraphTheoryPageID       int64 = 2998
	GeometryPageID          int64 = 1218
	ArithmetiquePageID      int64 = 2057625
	ProbabilityTheoryPageID int64 = 14993
)

func testGraph() *simple.DirectedGraph {
	g := simple.NewDirectedGraph()

	n1 := NewNode(1)
	g.AddNode(n1)
	n2 := NewNode(2)
	g.AddNode(n2)
	n3 := NewNode(3)
	g.AddNode(n3)
	n4 := NewNode(4)
	g.AddNode(n4)
	n5 := NewNode(5)
	g.AddNode(n5)
	n6 := NewNode(6)
	g.AddNode(n6)
	n7 := NewNode(7)
	g.AddNode(n7)
	n8 := NewNode(8)
	g.AddNode(n8)
	n9 := NewNode(9)
	g.AddNode(n9)
	n10 := NewNode(10)
	g.AddNode(n10)
	n11 := NewNode(11)
	g.AddNode(n11)
	n12 := NewNode(12)
	g.AddNode(n12)

	g.SetEdge(g.NewEdge(n1, n2))
	g.SetEdge(g.NewEdge(n1, n11))
	g.SetEdge(g.NewEdge(n1, n12))

	g.SetEdge(g.NewEdge(n2, n3))
	g.SetEdge(g.NewEdge(n2, n12))

	g.SetEdge(g.NewEdge(n3, n4))
	g.SetEdge(g.NewEdge(n3, n11))
	g.SetEdge(g.NewEdge(n3, n12))

	g.SetEdge(g.NewEdge(n4, n5))
	g.SetEdge(g.NewEdge(n4, n6))
	g.SetEdge(g.NewEdge(n4, n10))

	g.SetEdge(g.NewEdge(n5, n6))

	g.SetEdge(g.NewEdge(n6, n7))

	g.SetEdge(g.NewEdge(n7, n9))
	g.SetEdge(g.NewEdge(n7, n8))

	g.SetEdge(g.NewEdge(n8, n9))

	g.SetEdge(g.NewEdge(n10, n11))
	g.SetEdge(g.NewEdge(n10, n7))
	g.SetEdge(g.NewEdge(n10, n8))
	g.SetEdge(g.NewEdge(n10, n9))

	return g
}

func ClusteringShort(t *testing.T, clu Clusterer) {

	g := testGraph()

	j := Job{
		Model: string(TOUR),
	}

	clusters, err := clu.Cluster(j, 1, g)
	if err != nil {
		t.Errorf("Cluster: %s", err)
	}
	if clusters == nil {
		t.Fatalf("Error was nil but cluster is also nil")
	}

	// shoud have 2 layer
	d := clusters.Depth()
	if d != 2 {
		t.Errorf("Expected depth 2, got %d", d)
	}

}

func ClusteringMath(t *testing.T, clu Clusterer) {
	if testing.Short() {
		t.Skip()
	}

	j := Job{
		Model: string(ABSTRACT),
	}

	var MathPageID int64 = 3697062

	loader, err := loader.NewFileLoader("../../samples/mathematiques.dump.json")
	if err != nil {
		t.Fatalf("NewFileLoader: %s", err)
	}

	cla, err := classifier.NewV1(loader)
	if err != nil {
		t.Fatalf("NewV1: %s", err)
	}

	gm, err := cla.LoadGraph(MathPageID, clu.MaxSize(j))
	if err != nil {
		t.Fatalf("LoadGraph: %s", err)
	}

	clusters, err := clu.Cluster(j, MathPageID, gm)
	if err != nil {
		t.Fatalf("Cluster: %s", err)
	}
	if clusters == nil {
		t.Fatalf("Error was nil but cluster is also nil")
	}

	// shoud have 2 layer
	d := clusters.Depth()
	if d != 2 {
		t.Fatalf("Expected depth 2, got %d", d)
	}

	if !Find(MathPageID, clusters.Members) {
		t.Errorf("Expected Math in d1 members")
	}
	if !Find(GraphPageID, clusters.Members) {
		t.Errorf("Expected Graph in d1 members")
	}
	if !Find(GraphTheoryPageID, clusters.Members) {
		t.Errorf("Expected GraphTheory in d1 members")
	}
	if !Find(GeometryPageID, clusters.Members) {
		t.Errorf("Expected Geometry in d1 members")
	}
	if !Find(ArithmetiquePageID, clusters.Members) {
		t.Errorf("Expected Arithmetique in d1 members")
	}
	if !Find(ProbabilityTheoryPageID, clusters.Members) {
		t.Errorf("Expected ProbabilityTheory in d1 members")
	}

	if HasDuplicate(clusters) {
		t.Errorf("Found duplicate in clusters")
	}
}

func HasDuplicate(cluster *Cluster) bool {
	for id, _ := range cluster.Members {
		count := 0
		for _, c := range cluster.Subclusters {
			_, exists := c.Members[id]
			if exists {
				count++
			}
			if count > 1 {
				fmt.Printf("Found %d in 2 clusters", id)
				return false
			}
		}

	}

	return false
}

func Find(id int64, c Component) bool {
	for k, _ := range c {
		if k == id {
			return true
		}
	}

	return false
}

func TestIsHCS(t *testing.T) {
	var MathPageID int64 = 3697062
	loader, err := loader.NewFileLoader("../../samples/mathematiques.dump.json")
	if err != nil {
		t.Fatalf("NewFileLoader: %s", err)
	}

	cla, err := classifier.NewV1(loader)
	if err != nil {
		t.Fatalf("NewV1: %s", err)
	}

	gm, err := cla.LoadGraph(MathPageID, 0)
	if err != nil {
		t.Fatalf("LoadGraph: %s", err)
	}

	n1 := gm.Node(24899)
	n2 := gm.Node(3033966)
	n3 := gm.Node(3697062)

	m := make(map[int64]graph.Node)
	m[n1.ID()] = n1
	m[n2.ID()] = n2
	m[n3.ID()] = n3

	if isComponentHCS(gm, m) == false {
		t.Fatalf("Expected %v to be HCS", m)
	}

	/*
		2224
		6116
		8191
		10744
		12592
		17109
		17124
		30437
		50424
		3697062
		4084457
		4138953
	*/

}
