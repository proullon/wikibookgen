package clusterer

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"

	. "github.com/proullon/wikibookgen/api/model"
)

type Betweenness struct {
}

func NewBetweenness() *Betweenness {
	c := &Betweenness{}

	return c
}

func (c *Betweenness) Version() string {
	return "3"
}

func (c *Betweenness) MaxSize(j Job) int64 {
	return 10000
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *Betweenness) Cluster(j Job, rootID int64, g graph.Directed) (*Cluster, error) {

	bvalues := network.Betweenness(g)
	fmt.Printf("%d betweenness values\n", len(bvalues))

	if len(bvalues) == 0 {
		return nil, fmt.Errorf("Betweenness: no values")
	}

	for k, v := range bvalues {
		log.Infof("%v %v", k, v)
	}

	return nil, fmt.Errorf("not implemented")
}
