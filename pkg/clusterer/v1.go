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
	return 1000000
}

// Cluster will group given articles into highly connected group
// Now since it's a dumb V1, all it does is remove vertices with low trail count
func (c *V1) Cluster(j Job, rootID int64, g graph.Directed) (*Cluster, error) {

	var layer, maxpages, maxchapters int
	var maxtime time.Duration
	switch Model(j.Model) {
	case ABSTRACT:
		layer = 2
		maxtime = 5 * time.Minute
		maxpages = 100
		maxchapters = 10
	case TOUR:
		layer = 2
		maxtime = 60 * time.Minute
		maxpages = 500
		maxchapters = 50
	case ENCYCLOPEDIA:
		layer = 3
		maxtime = 120 * time.Minute
		maxpages = 10000
		maxchapters = 1000
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

	componentmap, maxidx, err := c.findcomponents(g, root, maxtime)
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

	unique := make(Component)
	for _, components := range componentmap {
		for _, component := range components {
			for _, n := range component {
				unique[n.ID()] = n
			}
		}
	}
	log.Infof("Found %d unique nodes", len(unique))

	componentmap, maxidx = removeDuplicateInClusters(componentmap, maxidx)
	log.Infof("Removed duplicate in clusters")

	for i := maxidx; i > 1; i-- {
		log.Infof("- %d components of size %d", len(componentmap[i]), i)
	}

	cluster.Members = make(Component)
	var chaptercount, pagecount int
	if layer >= 2 {
		for i := maxidx; i > 2; i-- {
			components := componentmap[i]
			for _, component := range components {
				if pagecount > maxpages || chaptercount > maxchapters {
					break
				}
				for id, n := range component {
					cluster.Members[id] = n
				}
				s := &Cluster{Members: component}
				cluster.Subclusters = append(cluster.Subclusters, s)
				pagecount += len(component)
				chaptercount++
			}
		}
	}

	log.Infof("Clustered %d unique nodes", len(cluster.Members))
	return cluster, nil
}

func (c *V1) findcomponents(g graph.Directed, root graph.Node, maxtime time.Duration) (map[int][]Component, int, error) {
	var maxidx int

	componentmap := make(map[int][]Component)
	forbidden := make(map[int64]bool)

	neighbours := createNeighboursPool(g, []graph.Node{root})
	neighbours = prepareNodePool(g, neighbours, 3)
	//sort.Sort(ByDegree(neighbours))

	allocation := len(neighbours)
	if len(neighbours) > 100 {
		allocation = len(neighbours) / 4
	}
	allocated := maxtime / time.Duration(allocation)
	globalTimeout := time.Now().Add(maxtime)

	for i, n := range neighbours {
		// if majority of BiggestBestCandidate search finished late, check for globalTimeout
		if time.Now().After(globalTimeout) {
			log.Infof("Stopping early (%d/%d done)", i, len(neighbours))
			break
		}
		// skip if neighbour was previously added in a component
		if skip, ok := forbidden[n.ID()]; ok && skip {
			continue
		}

		// search for biggest component with base root + neighbour
		cm := make(Component)
		cm[root.ID()] = root
		cm[n.ID()] = n
		log.Debugf("You got %s to find biggest component from %v", maxtime, cm)
		begin := time.Now()
		end := time.Now().Add(allocated)
		component := c.BiggestBestCandidate(end, g, cm, n, forbidden)
		if len(component) <= 2 {
			continue
		}

		log.Infof("Got %v from %v after %s with %s left", component, cm, time.Since(begin), time.Until(end))

		// add new component member to forbidden list
		for id, _ := range component {
			forbidden[id] = true
		}

		// add new component to map
		l := len(component)
		if l > maxidx {
			maxidx = l
		}
		componentmap[l] = append(componentmap[l], component)
	}

	return componentmap, maxidx, nil
}

func (c *V1) BiggestBestCandidate(timeout time.Time, g graph.Directed, roots Component, last graph.Node, forbidden map[int64]bool) Component {
	var biggest Component
	biggest = roots

	if len(biggest) > 5 && !biggest.CanGrow(g) {
		log.Infof("Component %v cannot grow to %d", biggest, len(biggest)+1)
		return biggest
	}

	neighbours := createComponentNeighboursPool(g, roots)
	//	neighbours := createNeighboursPool(g, []graph.Node{last})
	var allowedNeighbours []graph.Node
	for _, n := range neighbours {
		if _, exists := roots[n.ID()]; exists {
			continue
		}
		if _, exists := forbidden[n.ID()]; exists {
			continue
		}
		allowedNeighbours = append(allowedNeighbours, n)
	}
	neighbours = prepareNodePool(g, allowedNeighbours, len(roots)+1)

	//timeleft := time.Until(timeout)
	//maxtime := timeleft / time.Duration(len(neighbours))
	//log.Infof("%s until timeout. You got %s to find biggest component from  %v", timeleft, maxtime, roots)
	if _, ok := roots[2350068]; ok {
		log.Infof("%s to find biggest component from %v", time.Until(timeout), roots)
		timeout = time.Now().Add(1 * time.Minute)
	}

	bestfit, maxidx := biggest.BestCandidates(g, neighbours)

	for i := maxidx; i > 0; i-- {
		components := bestfit[i]
		for _, n := range components {

			if n.ID() == 2350068 {
				log.Infof("TESTING GRAPH with %s", roots)
				timeout = time.Now().Add(1 * time.Minute)
			}
			/*
				if time.Now().After(timeout) {
					return biggest
				}
			*/

			r2 := roots.Copy()
			r2[n.ID()] = n

			if time.Now().After(timeout) {
				return biggest
			}

			if !isComponentHCS(g, r2) {
				continue
			}

			c := c.BiggestBestCandidate(timeout, g, r2, n, forbidden)
			if len(c) > len(biggest) {
				biggest = c
			}
		}
	}

	return biggest
}

func (c *V1) BiggestBreadthFirst(timeout time.Time, g graph.Directed, roots Component, last graph.Node, forbidden map[int64]bool) Component {
	var biggest Component
	biggest = roots

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
	log.Debugf("%s to find biggest component from %v", time.Until(timeout), roots)

	var added []graph.Node

	for _, n := range neighbours {
		if _, exists := biggest[n.ID()]; exists {
			continue
		}
		if _, exists := forbidden[n.ID()]; exists {
			continue
		}
		/*
			if time.Now().After(timeout) {
				return biggest
			}
		*/

		r2 := biggest.Copy()
		r2[n.ID()] = n

		if !isComponentHCS(g, r2) {
			continue
		}

		added = append(added, n)
		biggest = r2
	}

	for _, n := range added {
		if time.Now().After(timeout) {
			return biggest
		}

		c := c.BiggestBreadthFirst(timeout, g, biggest, n, forbidden)
		if len(c) > len(biggest) {
			biggest = c
		}

	}

	return biggest
}

func (c *V1) BiggestDepthFirst(timeout time.Time, g graph.Directed, roots Component, last graph.Node, forbidden map[int64]bool) Component {
	var biggest Component
	biggest = roots

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
	if _, ok := roots[2350068]; ok {
		log.Infof("%s to find biggest component from %v", time.Until(timeout), roots)
		timeout = time.Now().Add(1 * time.Minute)
	}

	for _, n := range neighbours {
		if n.ID() == 2350068 {
			log.Infof("TESTING GRAPH with %s", roots)
		}
		if _, exists := roots[n.ID()]; exists {
			continue
		}
		if _, exists := forbidden[n.ID()]; exists {
			if n.ID() == 2350068 {
				log.Infof("GRAPH is forbidden :(")
			}
			continue
		}
		/*
			if time.Now().After(timeout) {
				return biggest
			}
		*/

		r2 := roots.Copy()
		r2[n.ID()] = n

		if time.Now().After(timeout) {
			return biggest
		}

		if !isComponentHCS(g, r2) {
			continue
		}

		c := c.BiggestDepthFirst(timeout, g, r2, n, forbidden)
		if len(c) > len(biggest) {
			biggest = c
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

func removeDuplicateInClusters(componentmap map[int][]Component, maxidx int) (map[int][]Component, int) {
	cmap := make(map[int][]Component)
	fmap := make(map[int64]bool)
	var newmaxidx int

	for i := maxidx; i > 2; i-- {
		components := componentmap[i]
		for _, component := range components {
			for id, _ := range component {
				// if already found, delete from component, if not add to fmap
				_, found := fmap[id]
				if found {
					delete(component, id)
				} else {
					fmap[id] = true
				}
			}
			l := len(component)
			//if l > 2 {
			cmap[l] = append(cmap[l], component)
			if l > newmaxidx {
				newmaxidx = l
			}
			//}
		}
	}

	return cmap, newmaxidx
}
