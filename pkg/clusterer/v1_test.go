package clusterer

import (
	//	"fmt"
	"testing"
	//	"github.com/proullon/graph"

	"gonum.org/v1/gonum/graph/simple"

	. "github.com/proullon/wikibookgen/api/model"
)

func testGraph() *simple.DirectedGraph {
	g := simple.NewDirectedGraph()

	n1 := NewNode(1)
	n2 := NewNode(2)
	n3 := NewNode(3)
	n4 := NewNode(4)
	n5 := NewNode(5)
	n6 := NewNode(6)
	n7 := NewNode(7)
	n8 := NewNode(8)
	n9 := NewNode(9)
	n10 := NewNode(10)
	n11 := NewNode(11)
	n12 := NewNode(12)

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

func TestClassifyV1(t *testing.T) {

	g := testGraph()

	clu := NewV1()

	j := Job{
		Model: string(TOUR),
	}

	clusters, err := clu.Cluster(j, g)
	if err != nil {
		t.Errorf("Cluster: %s", err)
	}

	_ = clusters
}
