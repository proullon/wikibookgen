package orderer

import (
	"encoding/json"
	"io/ioutil"
	"testing"

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

func LoadCluster() (*Cluster, error) {
	data, err := ioutil.ReadFile("../../samples/clusters.mathematiques.json")
	if err != nil {
		return nil, err
	}

	var clusters *Cluster
	err = json.Unmarshal(data, &clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func TestOrdering(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	j := Job{
		Model:   string(TOUR),
		Subject: "Mathematiques",
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

	gm, err := cla.LoadGraph(MathPageID, 0)
	if err != nil {
		t.Fatalf("LoadGraph: %s", err)
	}

	clusters, err := LoadCluster()
	if err != nil {
		t.Fatalf("LoadCluster: %s", err)
	}

	ord := NewV1()
	wikibook, err := ord.Order(loader, j, gm, clusters)
	if err != nil {
		t.Fatalf("Order: %s", err)
	}

	t.Logf("%v", wikibook)
}
