package clusterer

/*
import (
	"fmt"
	"testing"
	//	"github.com/proullon/graph"

	"github.com/proullon/wikibookgen/pkg/classifier"
	"github.com/proullon/wikibookgen/pkg/loader"

	. "github.com/proullon/wikibookgen/api/model"
)

func TestClassifyV1(t *testing.T) {
	MathPageID := 3697062

	loader, err := loader.NewFileLoader("../../samples/mathematiques.json")
	if err != nil {
		t.Errorf("NewFileLoader: %s", err)
	}

	cla, err := classifier.NewV1(loader)
	if err != nil {
		t.Errorf("NewV1: %s", err)
	}

	root, vertices, err := cla.LoadGraph(MathPageID)
	if err != nil {
		t.Errorf("LoadGraph: %s", err)
	}

	g := NewGraph(root, vertices)

	// graph should not be higly connected
	if HighlyConnected(g) {
		t.Errorf("Graph should not be highly connected (%d vertices, %d edges)", g.VertexCount(), g.EdgeCount())
	}

	clu := NewV1()

	j := Job{
		Model: string(TOUR),
	}

	clusters, err := clu.Cluster(j, root, vertices)
	if err != nil {
		t.Errorf("Cluster: %s", err)
	}

	_ = clusters
}

func HighlyConnected(g *Graph) bool {
	minimal_degree := g.VertexCount() / 2
	fmt.Printf("minimal_degree: %d\n", minimal_degree)

	vertices := g.Vertices()
	for _, v := range vertices {
		if v.Degree() < minimal_degree {
			return false
		}
	}

	return true
}
*/
