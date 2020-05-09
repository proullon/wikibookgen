package classifier

import (
	"fmt"
	"io/ioutil"
	"testing"

	//	"github.com/proullon/graph"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"

	"github.com/proullon/wikibookgen/pkg/loader"
	//	. "github.com/proullon/wikibookgen/api/model"
)

func TestClassifyV1(t *testing.T) {
	var MathPageID int64 = 3697062

	loader, err := loader.NewFileLoader("../../samples/mathematiques.json")
	if err != nil {
		t.Errorf("NewFileLoader: %s", err)
	}

	cla, err := NewV1(loader)
	if err != nil {
		t.Errorf("NewV1: %s", err)
	}

	g, err := cla.LoadGraph(MathPageID)
	if err != nil {
		t.Errorf("LoadGraph: %s", err)
	}
	t.Logf("Graph has %d vertices", g.Nodes().Len())

	// graph should not be higly connected
	if HighlyConnected(g) {
		t.Errorf("Graph should not be highly connected (%d vertices)", g.Nodes().Len())
	}

	data, err := dot.Marshal(g, "mathematiques", "", "")
	if err != nil {
		t.Errorf("cannot marshal graph: %s", err)
	}

	err = ioutil.WriteFile("./math.graph.dot", data, 0666)
	if err != nil {
		t.Errorf("cannot write graph: %s", err)
	}

}

// HighlyConnected determine if graph G is highly connected
// With Let n=|G|. If G is highly connected, every vertex has degree >= n/2
func HighlyConnected(g graph.Graph) bool {
	//	minimal_degree := g.VertexCount() / 2
	//nodes := g.Nodes()
	// edges
	minimal_degree := g.Nodes().Len() / 2
	fmt.Printf("minimal_degree: %d\n", minimal_degree)

	vertices := g.Nodes()

	for vertices.Next() {
		n := vertices.Node()
		degree := g.From(n.ID()).Len()
		if degree < minimal_degree {
			fmt.Printf("Node %d has only a degree %d\n", n.ID(), degree)
			return false
		}
	}

	return true
}

/*
func HighlyConnected(g *Graph) bool {
	minimal_degree := g.VertexCount() / 2

	vertices := g.Vertices()
	for _, v := range vertices {
		if v.Degree() < minimal_degree {
			return false
		}
	}

	return true
}
*/
