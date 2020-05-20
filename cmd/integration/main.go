package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/proullon/wikibookgen/api/model"
	"github.com/proullon/wikibookgen/pkg/classifier"
	"github.com/proullon/wikibookgen/pkg/clusterer"
	"github.com/proullon/wikibookgen/pkg/loader"

	log "github.com/sirupsen/logrus"

	wikibookgen "github.com/proullon/wikibookgen/api/sdk"
)

type test struct {
	Name string
	Test func() error
}

var tt = []test{
	{
		Name: "Order",
		Test: testOrder,
	},
	{
		Name: "Clustering",
		Test: testClustering,
	},
}

func main() {
	var endpoint string

	endpoint = os.Getenv("WIKIBOOKGEN_API")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8090"
	}

	wikibookgen.Initialize(endpoint, "")

	for _, t := range tt {
		err := t.Test()
		if err != nil {
			log.Errorf("Test %s: Failed (%s)\n", t.Name, err)
		} else {
			log.Infof("Test %s: OK\n", t.Name)
		}
	}

}

func testOrder() error {
	subject := `Mathématiques`
	model := `tour`

	// test invalid model
	orderID, err := wikibookgen.Order(subject, `invalidmodel`)
	if err == nil {
		return fmt.Errorf("Order(%s, invalidmodel): should have error", subject)
	}

	// test invalid subject
	orderID, err = wikibookgen.Order(`invalidsubject`, model)
	if err == nil {
		return fmt.Errorf("Order(invalidsubject, %s): should have error", model)
	}

	// test nominal
	orderID, err = wikibookgen.Order(subject, model)
	if err != nil {
		return fmt.Errorf("Order(%s): %s", subject, err)
	}
	if orderID == "" {
		return fmt.Errorf("Order(%s): invalid uuid", subject)
	}

	return nil
}

func ClusterDepth(c *Cluster) int {
	return clusterDepth(c, 1)
}

func clusterDepth(c *Cluster, depth int) int {
	fmt.Printf("Cluster d%d : %v\n", depth, c.Members)

	var max int = depth
	for _, sub := range c.Subclusters {
		d := clusterDepth(sub, depth+1)
		if d > max {
			max = d
		}
	}

	return max
}

func testLouvain() error {

	fmt.Printf("Initialize loader...\n")
	loader, err := loader.NewFileLoader("./samples/mathematiques.json")
	if err != nil {
		return fmt.Errorf("NewFileLoader: %s", err)
	}
	fmt.Printf("Initialize loader...DONE\n")

	var graphsize int64 = 5000

	for {
		graphsize += 1000
		fmt.Printf("Timing classify with size:%d...\n", graphsize)
		d, err := timeClassify(graphsize, loader)
		fmt.Printf("GraphSize:%-3d  Duration:%-15s Error:%s\n", graphsize, d, err)

		if d > 20*time.Hour {
			return nil
		}

	}

}

func timeClassify(graphsize int64, loader Loader) (time.Duration, error) {
	var MathPageID int64 = 3697062
	var d time.Duration

	cla, err := classifier.NewV1(loader)
	if err != nil {
		return d, fmt.Errorf("NewV1: %s", err)
	}

	fmt.Printf("Loading graph...\n")
	g, err := cla.LoadGraph(MathPageID, graphsize)
	if err != nil {
		return d, fmt.Errorf("LoadGraph: %s", err)
	}
	fmt.Printf("Loading graph...DONE\n")

	clu := clusterer.NewLouvain()

	j := Job{
		Model: string(TOUR),
	}

	begin := time.Now()
	clusters, err := clu.Cluster(j, MathPageID, g)
	d = time.Since(begin)
	if err != nil {
		return d, fmt.Errorf("Cluster: %s", err)
	}

	return d, fmt.Errorf("SUCCESS: d%d cluster", ClusterDepth(clusters))
}

func testClustering() error {
	j := Job{
		Model: string(TOUR),
	}

	var MathPageID int64 = 3697062

	loader, err := loader.NewFileLoader("./samples/mathematiques.dump.json")
	if err != nil {
		return fmt.Errorf("NewFileLoader: %s", err)
	}

	cla, err := classifier.NewV1(loader)
	if err != nil {
		return fmt.Errorf("NewV1: %s", err)
	}

	gm, err := cla.LoadGraph(MathPageID, 0)
	if err != nil {
		return fmt.Errorf("LoadGraph: %s", err)
	}

	clu := clusterer.NewV1()

	clusters, err := clu.Cluster(j, MathPageID, gm)
	if err != nil {
		return fmt.Errorf("Cluster: %s", err)
	}

	// shoud have 2 layer
	d := clusters.Depth()
	if d != 2 {
		return fmt.Errorf("Expected depth 2, got %d", d)
	}

	return nil
}
