package clusterer

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"

	. "github.com/proullon/wikibookgen/api/model"
)

type Tarjan struct {
}

func NewTarjan() *Tarjan {
	c := &Tarjan{}

	return c
}

func (c *Tarjan) Version() string {
	return "2"
}

func (c *Tarjan) MaxSize(j Job) int64 {
	return 30000
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *Tarjan) Cluster(j Job, rootID int64, g graph.Directed) (*Cluster, error) {

	components := topo.TarjanSCC(g)
	fmt.Printf("%d components\n", len(components))

	for _, component := range components {
		for _, n := range component {
			if n.ID() == rootID {
				log.Infof("Found root component (%d): %v", len(component), component)
				cluster := &Cluster{Members: component}
				return cluster, nil
			}
		}
	}

	return nil, fmt.Errorf("did not found root component")
}
