package clusterer

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/proullon/workerpool"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"

	. "github.com/proullon/wikibookgen/api/model"
)

type V2 struct {
}

func NewV2() *V2 {
	c := &V2{}

	return c
}

func (c *V2) Version() string {
	return "2"
}

func (c *V2) MaxSize(j Job) int64 {
	return 0 // No limit
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V2, all it does is remove vertices with low trail count
func (c *V2) Cluster(j Job, rootID int64, g graph.Directed) (*Cluster, error) {

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

	cluster := &Cluster{}

	root := g.Node(rootID)
	if root == nil {
		return nil, fmt.Errorf("root node (%d) not found", rootID)
	}
	log.Infof("Ach so, Root %d exists in graph", rootID)

	componentmap, maxidx, err := c.cluster(g, root)
	if err != nil {
		return nil, err
	}

	unique := make(map[int64]graph.Node)
	for _, components := range componentmap {
		for _, component := range components {
			if len(component) > 2 {
				cluster.Members = append(cluster.Members, component...)
			}
			for _, n := range component {
				unique[n.ID()] = n
			}
		}
	}

	log.Infof("Clustered %d unique nodes", len(unique))

	var printed int
	for i := maxidx; i > 0; i-- {
		components := componentmap[i]
		for _, component := range components {
			if printed > 50 {
				break
			}
			log.Infof("- %v", component)
			printed++
		}
	}

	if layer >= 2 {
		printed = 0
		for i := maxidx; i > 0; i-- {
			components := componentmap[i]
			for _, component := range components {
				if printed > 50 {
					break
				}
				cluster.Subclusters = append(cluster.Subclusters, &Cluster{Members: component})

				printed++
			}
		}

	}

	_ = componentmap
	_ = maxidx

	/*
		// create direct neigboor slice
		var neighbours []graph.Node
		nodes := g.From(rootID)
		for nodes.Next() {
			neighbours = append(neighbours, nodes.Node())
		}
		nodes = g.To(rootID)
		for nodes.Next() {
			neighbours = append(neighbours, nodes.Node())
		}
		sort.Sort(ByID(neighbours))
		log.Infof("root has %d neighbours nodes", len(neighbours))

		clusters, err := c.lookForHCS(g, []graph.Node{root}, neighbours)
		if err != nil {
			return nil, err
		}
		log.Infof("Got %d clusters !", len(clusters))

		members := make(map[int]graph.Node)
		for _, cl := range clusters {
			for _, n := range cl {
				members[int(n.ID())] = n
				if n.ID() == 818568 {
					log.Infof("Found graph page id ! (degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
				if n.ID() == 1218 {
					log.Infof("Found geometry page id ! (degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
				if n.ID() == 1566943 {
					log.Infof("Found arthmetique page id !(degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
			}
		}
		log.Infof("Got %d unique clustered nodes from %d neighbours", len(members), len(neighbours))

		// reclustering := make(map[int]int)
		var reclustering []graph.Node
		for _, v := range members {
			reclustering = append(reclustering, v)
			cluster.Members = append(cluster.Members, v)
		}

		log.Infof("Trying to grow %d existing clusters", len(clusters))

		var clusters4 [][]graph.Node
		for _, roots := range clusters {
			var neighbours4 []graph.Node
			for _, n := range roots {
				nodes := g.From(n.ID())
				for nodes.Next() {
					neighbours4 = append(neighbours4, nodes.Node())
				}
				nodes = g.To(n.ID())
				for nodes.Next() {
					neighbours4 = append(neighbours4, nodes.Node())
				}
				sort.Sort(ByID(neighbours4))
			}
			c4, err := c.ComputeHCS(g, roots, 4, neighbours4)
			if err != nil {
				log.Errorf("Could not ComputeHCS of size 4 for root %v", roots)
				continue
			}

			if len(c4) > 0 {
				log.Infof("Found %d clusters of len 4 from root %v", len(c4), roots)
				clusters4 = append(clusters4, c4...)
			}

		}
		log.Infof("Got %d C4 clusters", len(clusters4))
		for _, cl := range clusters4 {
			for _, n := range cl {
				if n.ID() == 818568 {
					log.Infof("Found graph page id in cluster 4! (degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
				if n.ID() == 1218 {
					log.Infof("Found geometry page id in cluster 4! (degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
				if n.ID() == 1566943 {
					log.Infof("Found arthmetique page id in cluster 4!(degree %d)", g.From(n.ID()).Len()+g.To(n.ID()).Len())
				}
			}
		}

		var clusters5 [][]graph.Node
		for _, roots := range clusters4 {
			var neighbours5 []graph.Node
			for _, n := range roots {
				nodes := g.From(n.ID())
				for nodes.Next() {
					neighbours5 = append(neighbours5, nodes.Node())
				}
				nodes = g.To(n.ID())
				for nodes.Next() {
					neighbours5 = append(neighbours5, nodes.Node())
				}
				sort.Sort(ByID(neighbours5))
			}
			c5, err := c.ComputeHCS(g, roots, 5, neighbours5)
			if err != nil {
				log.Errorf("Could not ComputeHCS of size 5 for root %v", roots)
				continue
			}

			if len(c5) > 0 {
				log.Infof("Found %d clusters of len 4 from root %v", len(c5), roots)
				clusters5 = append(clusters5, c5...)
			}

		}
		log.Infof("Got %d C5 clusters", len(clusters5))
	*/
	/*
		var cc, tc int
		for _, n1 := range neighbours {
			for _, n2 := range neighbours {
				if n1.ID() == n2.ID() {
					continue
				}
				log.Debugf("Is {%d, %d, %d} highly connected ?", root.ID(), n1.ID(), n2.ID())
				var prospect []graph.Node = []graph.Node{root, n1, n2}
				tc++
				if HCS(g, prospect) {
					log.Warnf("{%d, %d, %d} is a HCS of G containing root", root.ID(), n1.ID(), n2.ID())
					cc++
				}
				if cc == 60 {
					log.Infof("Already 60 cluster found on %d test, stop", tc)
					return cluster, nil
				}
			}
		}
	*/
	return cluster, nil
}

func (c *V2) cluster(g graph.Directed, root graph.Node) (map[int][][]graph.Node, int, error) {
	log.Infof("v2.cluster")
	cbegin := time.Now()

	//testn := g.Node(2350068)
	//	testn := g.Node(2057625)

	componentmap := make(map[int][][]graph.Node)
	//previous := [][]graph.Node{[]graph.Node{root, testn}}
	//	previous := [][]graph.Node{[]graph.Node{root}}

	// create 2 node components with root and each of its neighbour
	var previous [][]graph.Node
	nodes := g.From(root.ID())
	for nodes.Next() {
		p := []graph.Node{root, nodes.Node()}
		sort.Sort(ByID(p))
		previous = append(previous, p)
	}
	nodes = g.To(root.ID())
	for nodes.Next() {
		p := []graph.Node{root, nodes.Node()}
		sort.Sort(ByID(p))
		previous = append(previous, p)
	}

	componentmap[2] = previous

	var maxidx, removed int
	componentmap, maxidx, removed = mergecomponentmap(g, componentmap)
	log.Infof("Removed %d components in map after cleanup. Biggest cluster %d", removed, maxidx)

	previous = componentmap[2]

	target := 3
	for {
		for i := 2; i <= maxidx; i++ {
			log.Infof("- %d components of size %d", len(componentmap[i]), i)
		}

		log.Infof("Looking for connected components of size %d from %d root component", target, len(previous))
		begin := time.Now()

		components, err := c.Grow(g, previous, target)
		log.Infof("Looking for connected components of size %d done. Got %d components after %s", target, len(components), time.Since(begin))
		if err != nil {
			return nil, 0, err
		}

		log.Infof("Adding %d components in componentmap", len(components))
		componentmap[target] = append(componentmap[target], components...)

		var removed, m int
		componentmap, m, removed = mergecomponentmap(g, componentmap)
		if m > maxidx {
			maxidx = m
		}
		log.Infof("Removed %d components in map after cleanup. Biggest cluster %d", removed, maxidx)

		/*
				var removed int
				for {
					components = removeduplicate(components)
					components, removed = mergecomponents(g, components)
					if removed == 0 {
						break
					}
				}
			log.Infof("Got %d components of size %d after cleanup", len(components), target)
		*/

		/*
			for _, com := range components {
				found := false
				for _, n := range com {
					if n.ID() == root.ID() {
						found = true
						break
					}
				}
				if !found {
					log.Errorf("Root node not found in component %v", com)
				}
			}
		*/

		/*
			if len(components) == 0 {
				break
			}
		*/

		/*
			for _, comp := range components {
				l := len(comp)
				if l > maxidx {
					maxidx = l
				}
				componentmap[l] = append(componentmap[l], comp)
			}
		*/

		if target > maxidx {
			break
		}

		if time.Since(cbegin) > 120*time.Minute {
			break
		}

		//componentmap[target] = components
		//previous = components
		previous = componentmap[target]
		//maxidx = target
		target++
	}

	log.Infof("v2.cluster done:")
	for i := 2; i <= maxidx; i++ {
		log.Infof("- %d components of size %d", len(componentmap[i]), i)
	}

	log.Infof("v2.cluster total time: %s", time.Since(cbegin))
	return componentmap, maxidx, nil
}

func (c *V2) Grow(g graph.Directed, roots [][]graph.Node, target int) ([][]graph.Node, error) {
	var components [][]graph.Node

	for _, root := range roots {
		neighbours := createNeighboursPool(g, root)

		s, err := c.ComputeHCS(g, root, target, neighbours)
		if err != nil {
			return nil, err
		}

		components = append(components, s...)
	}
	/*
		for _, cl := range components {
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
	*/
	return components, nil
}

func (c *V2) ComputeHCS(g graph.Directed, roots []graph.Node, size int, neighbours []graph.Node) ([][]graph.Node, error) {
	var clusters [][]graph.Node

	begin := time.Now()
	// trim neighbours having degree < target |G| / 2, we know we won't be able to add them
	pool := prepareNodePool(g, neighbours, size)

	// do we still have enough nodes to form at least 1 cluster ?
	if len(pool) < size-len(roots) {
		log.Infof("Only %d node in pool for target of size %d, giving up", len(pool), size)
		return clusters, nil
	}

	log.Debugf("Looking for HCS of size %d with pool of %d nodes with roots %v", size, len(pool), roots)
	//		tested := make(map[int][][]graph.Node)
	//		clusters := c.kHCSrec(g, pool, i, []graph.Node{root}, tested)
	hcs, err := NewHCS()
	if err != nil {
		return nil, err
	}

	hcs.Find(g, pool, size, roots)

	clusters = hcs.Clusters()
	log.Debugf("Looked for HCS of size %d with pool of %d nodes with roots %v: took %s got %d clusters", size, len(pool), time.Since(begin), roots, len(clusters))

	return clusters, nil
}

/*
func (c *V2) lookForHCS(g graph.Directed, roots []graph.Node, neighbours []graph.Node) ([][]graph.Node, error) {
	max := 3
	min := 3

	//	max = 7
	//	min = 6

	for i := max; i >= min; i-- {
		begin := time.Now()
		// trim neighbours having degree < target |G| / 2, we know we won't be able to add them
		pool := prepareNodePool(g, neighbours, i)

		// do we still have enough nodes to form at least 1 cluster ?
		if len(pool) < i-len(roots) {
			log.Infof("Only %d node in pool for target of size %d, giving up", len(pool), i)
			continue
		}

		log.Debugf("Looking for HCS of size %d with pool of %d nodes", i, len(pool))
		//		tested := make(map[int][][]graph.Node)
		//		clusters := c.kHCSrec(g, pool, i, []graph.Node{root}, tested)
		hcs, err := NewHCS()
		if err != nil {
			return nil, err
		}

		hcs.Find(g, pool, i, roots)
		log.Debugf("Looked for HCS of size %d with pool of %d nodes: took %s", i, len(pool), time.Since(begin))

		clusters := hcs.Clusters()

		if len(clusters) == 0 {
			log.Infof("Found 0 clusters of size %d", i)
			continue
		}


		return clusters, nil
	}

	return nil, nil
}
*/

type HCS struct {
	wp     *workerpool.WorkerPool
	tested map[int][][]graph.Node
	mu     sync.Mutex
}

type HCSPayload struct {
	g        graph.Directed
	prospect []graph.Node
}

func NewHCS() (*HCS, error) {
	hcs := &HCS{}

	var err error
	hcs.wp, err = workerpool.New(hcs.check,
		workerpool.WithMaxWorker(runtime.NumCPU()),
		workerpool.WithSizePercentil(workerpool.AllInSizesPercentil),
		workerpool.WithMaxQueue(1000),
	)
	if err != nil {
		return nil, err
	}

	hcs.tested = make(map[int][][]graph.Node)

	return hcs, nil
}

func (hcs *HCS) check(payload interface{}) (interface{}, error) {
	p := payload.(*HCSPayload)
	var prospect []graph.Node
	prospect = append(prospect, p.prospect...)

	f := false
	for _, n := range prospect {
		if n.ID() == 3697062 {
			f = true
			break
		}
	}
	if !f {
		log.Errorf("HCS.check: Root not found in prospect %v", prospect)
		return nil, fmt.Errorf("root not found oO")
	}

	sort.Sort(ByID(prospect))

	hcs.mu.Lock()
	tested := alreadyTested(hcs.tested, prospect)
	hcs.mu.Unlock()
	if tested {
		log.Debugf("%v alreadyTested", prospect)
		return nil, nil
	}

	hcs.mu.Lock()
	hcs.tested[int(prospect[0].ID())] = append(hcs.tested[int(prospect[0].ID())], prospect)
	hcs.mu.Unlock()

	if isHCS(p.g, p.prospect) {
		log.Debugf("HCS.check: Found %v", prospect)
		return p.prospect, nil
	}

	return nil, nil
}

func (hcs *HCS) Find(g graph.Directed, pool []graph.Node, targetsize int, prospect []graph.Node) {

	if len(prospect) == targetsize {
		/*
			sort.Sort(ByID(prospect))
			if alreadyTested(tested, prospect) {
				log.Debugf("%v alreadyTested", prospect)
				return nil
			}
			tested[int(prospect[0].ID())] = append(tested[int(prospect[0].ID())], prospect)

			if HCS(g, prospect) {
				log.Infof("Found %v !", prospect)
				return [][]graph.Node{prospect}
			}
			return nil
		*/

		payload := &HCSPayload{g: g, prospect: prospect}
		hcs.wp.Feed(payload)
		return
	}

	f := false
	for _, n := range prospect {
		if n.ID() == 3697062 {
			f = true
			break
		}
	}
	if !f {
		log.Errorf("Root not found in prospect %v", prospect)
		return
	}

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
		hcs.Find(g, pool, targetsize, p)

		//log.Infof("HCS %d", len(hcs))

		/*
			if len(hcs) > 0 {
				//return hcs
				for _, s := range hcs {
					//log.Infof("Rec returned %-3d cluster %v", i, s)
					clusters = append(clusters, s)
				}
			}
		*/

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

		/*
			if len(clusters) >= 60 {
				return clusters
			}
		*/
	}

	//	return clusters
}

func (hcs *HCS) Clusters() [][]graph.Node {
	hcs.wp.Wait()
	hcs.wp.Stop()

	var clusters [][]graph.Node

	for r := range hcs.wp.ReturnChannel {
		c := r.Body.([]graph.Node)
		log.Debugf("Clusters: got %v", c)
		clusters = append(clusters, c)
	}

	return clusters
}
