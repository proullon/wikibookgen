package classifier

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/proullon/wikibookgen/pkg/loader"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
)

var (
	MathPageID              int64 = 3697062
	GraphPageID             int64 = 818568
	GraphMathPageID         int64 = 2642295
	GraphTheoryPageID       int64 = 2350068
	GeometryPageID          int64 = 1218
	ArithmetiquePageID      int64 = 2057625
	ProbabilityTheoryPageID int64 = 14993
)

func TestClassifyV1(t *testing.T) {

	loader, err := loader.NewFileLoader("../../samples/mathematiques.dump.json")
	if err != nil {
		t.Errorf("NewFileLoader: %s", err)
	}

	cla, err := NewV1()
	if err != nil {
		t.Errorf("NewV1: %s", err)
	}

	g, err := cla.LoadGraph(loader, MathPageID, 0)
	if err != nil {
		t.Errorf("LoadGraph: %s", err)
	}
	t.Logf("Graph has %d vertices", g.Nodes().Len())

	// Graph should be present in math graph
	n := g.Node(GraphPageID)
	if n == nil {
		t.Errorf("Expected GraphPage in Math graph")
	}
	if degree := Degree(g, n); degree < 20 || degree > 200 {
		t.Errorf("Expected Graph page to have degree  80 < d < 200, got %d", degree)
	}

	n = g.Node(GeometryPageID)
	if n == nil {
		t.Errorf("Expected Geometry in Math graph")
	}
	if degree := Degree(g, n); degree < 100 {
		t.Errorf("Expected Geometry page to have degree > 100, got %d", degree)
	}

	n = g.Node(ProbabilityTheoryPageID)
	if n == nil {
		t.Errorf("Expected ProbabilityTheory in Math graph")
	}
	if degree := Degree(g, n); degree < 100 {
		t.Errorf("Expected ProbabilityTheory page to have degree > 100, got %d", degree)
	}

	if g.Node(ArithmetiquePageID) == nil {
		t.Errorf("Expected Arithmetique in Math graph")
	}

	//DisplayPath(g)
}

func DisplayPath(g graph.Graph) {
	p, expanded := path.AStar(g.Node(MathPageID), g.Node(GraphPageID), g, nil)
	t, w := p.To(GraphPageID)
	fmt.Printf("Math->Graph: %s %f %d\n", t, w, expanded)

	p, expanded = path.AStar(g.Node(GraphPageID), g.Node(MathPageID), g, nil)
	t, w = p.To(MathPageID)
	fmt.Printf("Graph->Math: %s %f %d\n", t, w, expanded)

	p, expanded = path.AStar(g.Node(MathPageID), g.Node(GraphMathPageID), g, nil)
	t, w = p.To(GraphMathPageID)
	fmt.Printf("Math->GraphMath: %s %f %d\n", t, w, expanded)

	p, expanded = path.AStar(g.Node(GraphMathPageID), g.Node(MathPageID), g, nil)
	t, w = p.To(MathPageID)
	fmt.Printf("GraphMath->Math: %s %f %d\n", t, w, expanded)

	p, expanded = path.AStar(g.Node(GeometryPageID), g.Node(MathPageID), g, nil)
	t, w = p.To(MathPageID)
	fmt.Printf("Geometry->Math: %s %f %d\n", t, w, expanded)
	p, expanded = path.AStar(g.Node(MathPageID), g.Node(ArithmetiquePageID), g, nil)
	t, w = p.To(ArithmetiquePageID)
	fmt.Printf("Math->Arithmetique: %s %f %d\n", t, w, expanded)
}

func DumpGraph(g graph.Graph) error {
	data, err := dot.Marshal(g, "mathematiques", "", "")
	if err != nil {
		return fmt.Errorf("cannot marshal graph: %s", err)
	}

	err = ioutil.WriteFile("./math.graph.dot", data, 0666)
	if err != nil {
		return fmt.Errorf("cannot write graph: %s", err)
	}

	return nil
}

func Degree(g graph.Directed, n graph.Node) int {
	return g.From(n.ID()).Len() + g.To(n.ID()).Len()
}
