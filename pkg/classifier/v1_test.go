package classifier

import (
	"testing"

	//	"github.com/proullon/graph"

	"github.com/proullon/wikibookgen/pkg/loader"

	. "github.com/proullon/wikibookgen/api/model"
)

func TestClassifyV1(t *testing.T) {
	MathPageID := 3697062

	loader, err := loader.NewFileLoader("../../samples/mathematiques.json")
	if err != nil {
		t.Errorf("NewFileLoader: %s", err)
	}

	cla, err := NewV1(loader)
	if err != nil {
		t.Errorf("NewV1: %s", err)
	}

	root, vertices, err := cla.LoadGraph(MathPageID)
	if err != nil {
		t.Errorf("LoadGraph: %s", err)
	}

	var inedge, outedge int
	for _, v := range vertices {
		inedge += len(v.IncomingEdges)
		outedge += len(v.OutgoingEdges)
	}

	if inedge != outedge {
		t.Errorf("Expected same number of incoming and outgoing edges, got %d incoming and %d outgoing", inedge, outedge)
	}

	g := NewGraph(root, vertices)

	// graph should not be higly connected
	if HighlyConnected(g) {
		t.Errorf("Graph should not be highly connected (%d vertices, %d edges)", g.VertexCount(), g.EdgeCount())
	}
}

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
