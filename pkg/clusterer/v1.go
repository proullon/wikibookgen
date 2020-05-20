package clusterer

import (
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
}

func NewV1() *V1 {
	c := &V1{}

	return c
}

func (c *V1) Version() string {
	return "1"
}

func (c *V1) MaxSize(j Job) int64 {
	return 0 // No limit
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *V1) Cluster(j Job, rootID int64, g graph.Directed) (*Cluster, error) {

	var layer int
	switch Model(j.Model) {
	case ABSTRACT:
		layer = 1
	case TOUR:
		layer = 2
	case ENCYCLOPEDIA:
		layer = 3
	}

	graphsize := g.Nodes().Len()
	log.Infof("Graph has %d vertices", graphsize)

	root := g.Node(rootID)
	if root == nil {
		return nil, fmt.Errorf("root node (%d) not found", rootID)
	}
	log.Infof("Ach so, Root %d exists in graph", rootID)

	_ = layer
	cluster := &Cluster{}

	componentmap, maxidx, err := c.findcomponents(g, root)
	if err != nil {
		return nil, err
	}

	for i := maxidx; i > 1; i-- {
		for _, cl := range componentmap[i] {
			for _, n := range cl {
				if n.ID() == 818568 {
					log.Infof("Found graph page id ! (degree %d) in %v", g.From(n.ID()).Len()+g.To(n.ID()).Len(), cl)
				}
				if n.ID() == 2350068 {
					log.Infof("Found graph theory page id ! (degree %d) in %v", g.From(n.ID()).Len()+g.To(n.ID()).Len(), cl)
				}
				if n.ID() == 14993 {
					log.Infof("Found probability theory page id ! (degree %d) in %v", g.From(n.ID()).Len()+g.To(n.ID()).Len(), cl)
				}
				if n.ID() == 1218 {
					log.Infof("Found geometry page id ! (degree %d) in %v", g.From(n.ID()).Len()+g.To(n.ID()).Len(), cl)
				}
				if n.ID() == 2057625 {
					log.Infof("Found arithmetique page id !(degree %d) in %v", g.From(n.ID()).Len()+g.To(n.ID()).Len(), cl)
				}
			}
		}
	}

	for i := maxidx; i > 1; i-- {
		log.Infof("- %d components of size %d", len(componentmap[i]), i)
	}
	return cluster, nil
}

func (c *V1) findcomponents(g graph.Directed, root graph.Node) (map[int][]Component, int, error) {
	var maxidx int

	componentmap := make(map[int][]Component)

	neighbours := createNeighboursPool(g, []graph.Node{root})
	neighbours = prepareNodePool(g, neighbours, 3)

	maxtime := 15 * time.Minute / time.Duration(len(neighbours))

	for _, n := range neighbours {
		cm := make(Component)
		cm[root.ID()] = root
		cm[n.ID()] = n
		//log.Infof("You got %s to find biggest component from %v", maxtime, cm)
		begin := time.Now()
		component := c.Biggest(time.Now().Add(maxtime), g, cm, n)
		log.Infof("Got %v from %v after %s", component, cm, time.Since(begin))
		l := len(component)
		if l > maxidx {
			maxidx = l
		}
		componentmap[l] = append(componentmap[l], component)
	}

	return componentmap, maxidx, nil
}

func (c *V1) Biggest(timeout time.Time, g graph.Directed, roots Component, last graph.Node) Component {
	var biggest Component
	biggest = roots

	if !isComponentHCS(g, roots) {
		return nil
	}

	/*
		if !biggest.CanGrow(g) {
			log.Infof("Component %v cannot grow to %d", biggest, len(biggest))
			return biggest
		}
	*/

	//	neighbours := createComponentNeighboursPool(g, roots)
	neighbours := createNeighboursPool(g, []graph.Node{last})
	neighbours = prepareNodePool(g, neighbours, len(roots)+1)

	//timeleft := time.Until(timeout)
	//maxtime := timeleft / time.Duration(len(neighbours))
	//log.Infof("%s until timeout. You got %s to find biggest component from  %v", timeleft, maxtime, roots)
	//log.Infof("%s to find biggest component from %v", timeleft, roots)

	for _, n := range neighbours {
		if _, exists := roots[n.ID()]; exists {
			continue
		}
		/*
			if time.Now().After(timeout) {
				return biggest
			}
		*/

		r2 := roots.Copy()
		r2[n.ID()] = n

		c := c.Biggest(timeout, g, r2, n)
		if len(c) > len(biggest) {
			biggest = c
		}
		if time.Now().After(timeout) {
			return biggest
		}
	}

	return biggest
}

func lookForHCS(g graph.Directed, root graph.Node, neigboors []graph.Node) [][]graph.Node {
	max := 100
	min := 3

	max = 6
	min = 6

	for i := max; i >= min; i-- {
		begin := time.Now()
		// trim neigboors having degree < target |G| / 2, we know we won't be able to add them
		pool := prepareNodePool(g, neigboors, i)

		// do we still have enough nodes to form at least 1 cluster ?
		if len(pool) < i-1 {
			log.Infof("Only %d node in pool for target of size %d, giving up", len(pool), i)
			continue
		}

		log.Infof("Looking for HCS of size %d with pool of %d nodes", i, len(pool))
		tested := make(map[int][][]graph.Node)
		clusters := kHCSrec(g, pool, i, []graph.Node{root}, tested)
		log.Infof("Looked for HCS of size %d with pool of %d nodes: took %s", i, len(pool), time.Since(begin))

		if len(clusters) == 0 {
			log.Infof("Found 0 clusters of size %d, trying smaller", i)
			continue
		}

		/*
			if len(clusters) >= 60 {
				log.Infof("Found at least %d clusters of size %d, trying bigger", len(clusters), i)
				continue
			}
		*/

		return clusters
	}

	return nil
}

/*
func kHCS(g graph.Graph, pool []graph.Node, size int) [][]graph.Node {
	var clusters [][]graph.Node

	for _, n := range pool {
		var prospect []graph.Node
		prospect = append(prospect, n)
		hcs := kHCSrec(g, pool, size, prospect)
		clusters = append(clusters, hcs...)

		if len(clusters) >= 60 {
			return clusters
		}
	}

	return clusters
}
*/

func kHCSrec(g graph.Graph, pool []graph.Node, targetsize int, prospect []graph.Node, tested map[int][][]graph.Node) [][]graph.Node {

	if len(prospect) == targetsize {
		sort.Sort(ByID(prospect))
		if alreadyTested(tested, prospect) {
			log.Debugf("%v alreadyTested", prospect)
			return nil
		}
		tested[int(prospect[0].ID())] = append(tested[int(prospect[0].ID())], prospect)

		if isHCS(g, prospect) {
			log.Infof("Found %v !", prospect)
			return [][]graph.Node{prospect}
		}
		return nil
	}

	var clusters [][]graph.Node

	for _, n := range pool {
		// do not repeat node already in prospect
		var found bool = false
		for _, pn := range prospect {
			if pn.ID() == n.ID() {
				found = true
				break
			}
		}
		if found {
			continue
		}

		var p []graph.Node = append(prospect, n)
		hcs := kHCSrec(g, pool, targetsize, p, tested)

		//log.Infof("HCS %d", len(hcs))

		if len(hcs) > 0 {
			//return hcs
			for _, s := range hcs {
				//log.Infof("Rec returned %-3d cluster %v", i, s)
				clusters = append(clusters, s)
			}
		}

		//for _, s := range hcs {
		//log.Infof("Rec returned %-3d cluster %v", i, s)
		//clusters = append(clusters, s)
		//}
		//		clusters = append(clusters, hcs...)

		/*
			for i, c := range clusters {
				log.Infof("CLUSTER %-3d : %v", i, c)
			}
		*/

		if len(clusters) >= 60 {
			return clusters
		}
	}

	return clusters
}
