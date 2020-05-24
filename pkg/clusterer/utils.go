package clusterer

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"

	. "github.com/proullon/wikibookgen/api/model"
)

func Degree(g graph.Directed, n graph.Node) int {
	return g.From(n.ID()).Len() + g.To(n.ID()).Len()
}

// With Let n=|G|. If G is highly connected, every vertex has degree >= n/2
func isHCS(g graph.Graph, nodes []graph.Node) bool {
	var minimal_degree float64 = float64(float64(len(nodes)) / 2.0)
	log.Debugf("minimal_degree is %.2f", minimal_degree)

	for _, n1 := range nodes {
		var count = 0
		for _, n2 := range nodes {
			if n1.ID() == n2.ID() {
				continue
			}
			if g.HasEdgeBetween(n1.ID(), n2.ID()) {
				count++
			}
		}
		if float64(count) < minimal_degree {
			log.Debugf("Node %d has only %d edges (< %.1f) between prospect HCS", n1.ID(), count, minimal_degree)
			return false
		}
	}

	return true
}

// With Let n=|G|. If G is highly connected, every vertex has degree >= n/2
func isComponentHCS(g graph.Directed, c Component) bool {
	var minimal_degree float64 = float64(float64(len(c)) / 2.0)
	log.Debugf("minimal_degree is %.2f", minimal_degree)

	for _, n1 := range c {
		var count = 0
		if float64(g.From(n1.ID()).Len()+g.To(n1.ID()).Len()) < minimal_degree {
			return false
		}
		for _, n2 := range c {
			if n1.ID() == n2.ID() {
				continue
			}
			if g.HasEdgeBetween(n1.ID(), n2.ID()) {
				log.Debugf("%d - %d", n1.ID(), n2.ID())
				count++
			}
		}
		if float64(count) < minimal_degree {
			log.Debugf("Node %d has only %d edges (< %.1f) between prospect HCS", n1.ID(), count, minimal_degree)
			return false
		}
	}

	return true
}

func createComponentNeighboursPool(g graph.Directed, roots Component) []graph.Node {
	var neighbours []graph.Node
	var total int
	unique := make(map[int64]graph.Node)

	for _, n := range roots {
		nodes := g.From(n.ID())
		for nodes.Next() {
			unique[nodes.Node().ID()] = nodes.Node()
			total++
		}
		nodes = g.To(n.ID())
		for nodes.Next() {
			unique[nodes.Node().ID()] = nodes.Node()
			total++
		}
	}

	for _, n := range unique {
		neighbours = append(neighbours, n)
	}

	sort.Sort(ByID(neighbours))

	return neighbours
}

func createNeighboursPool(g graph.Directed, roots []graph.Node) []graph.Node {
	var neighbours []graph.Node
	var total int
	unique := make(map[int64]graph.Node)

	for _, n := range roots {
		nodes := g.From(n.ID())
		for nodes.Next() {
			unique[nodes.Node().ID()] = nodes.Node()
			total++
		}
		nodes = g.To(n.ID())
		for nodes.Next() {
			unique[nodes.Node().ID()] = nodes.Node()
			total++
		}
	}

	for _, n := range unique {
		neighbours = append(neighbours, n)
	}

	//sort.Sort(ByID(neighbours))

	return neighbours
}

func prepareNodePool(g graph.Directed, neighbours []graph.Node, targetsize int) []graph.Node {
	var minimal_degree float64 = float64(float64(targetsize) / 2.0)
	var prospect []graph.Node

	for _, n := range neighbours {
		degree := Degree(g, n)
		if float64(degree) >= minimal_degree {
			prospect = append(prospect, n)
		}
	}

	return prospect
}

func alreadyTested(tested map[int][][]graph.Node, prospect []graph.Node) bool {

	list := tested[int(prospect[0].ID())]

	for _, c := range list {
		same := true
		for i := range c {
			if c[i] != prospect[i] {
				same = false
				break
			}
		}
		if same {
			return true
		}
	}

	return false
}

// expect sorted components
func componentcmp(c1 []graph.Node, c2 []graph.Node) bool {
	if len(c1) != len(c2) {
		return false
	}

	for i := range c1 {
		if c1[i].ID() != c2[i].ID() {
			return false
		}
	}

	return true
}

func removeduplicate(components [][]graph.Node) [][]graph.Node {
	var nodup [][]graph.Node

	begin := time.Now()

	for _, c := range components {
		found := false
		for _, dc := range nodup {
			if componentcmp(c, dc) {
				found = true
				break
			}
		}

		if !found {
			nodup = append(nodup, c)
		}
	}

	log.Infof("duplicate removal took %s. Removed %d components", time.Since(begin), len(components)-len(nodup))
	return nodup
}

func removeDuplicateComponents(components []Component) []Component {
	var nodup []Component

	begin := time.Now()

	for _, c := range components {
		found := false
		for _, dc := range nodup {
			if c.Equal(dc) {
				found = true
				break
			}
		}

		if !found {
			nodup = append(nodup, c)
		}
	}

	log.Infof("duplicate removal took %s. Removed %d components", time.Since(begin), len(components)-len(nodup))
	return nodup
}

func exists(c []graph.Node, components [][]graph.Node) bool {
	for _, comp := range components {
		if componentcmp(c, comp) {
			return true
		}
	}

	return false
}

func mergecomponentold(g graph.Graph, c1 []graph.Node, c2 []graph.Node) ([]graph.Node, bool) {
	var merged []graph.Node

	mmap := make(map[int64]graph.Node)

	for _, n := range c1 {
		mmap[n.ID()] = n
	}
	for _, n := range c2 {
		mmap[n.ID()] = n
	}

	for _, v := range mmap {
		merged = append(merged, v)
	}

	hcs := isHCS(g, merged)
	return merged, hcs
}

func mergeComponent(g graph.Directed, c1 Component, c2 Component) (Component, bool) {

	mmap := make(map[int64]graph.Node)

	for _, n := range c1 {
		mmap[n.ID()] = n
	}
	for _, n := range c2 {
		mmap[n.ID()] = n
	}

	hcs := isComponentHCS(g, mmap)
	return mmap, hcs
}

func mergecomponent(g graph.Directed, c1 []graph.Node, c2 []graph.Node) (Component, bool) {

	mmap := make(map[int64]graph.Node)

	for _, n := range c1 {
		mmap[n.ID()] = n
	}
	for _, n := range c2 {
		mmap[n.ID()] = n
	}

	hcs := isComponentHCS(g, mmap)
	return mmap, hcs
}

func mergeComponents(g graph.Directed, components []Component) ([]Component, int) {
	var mergedcomps []Component
	mergedmap := make(map[int]bool)

	log.Infof("mergecomponents: Merging %d components", len(components))
	begin := time.Now()

	c := components[0]

	for i := 1; i < len(components); i++ {
		c2 := components[i]
		if mc, ok := mergeComponent(g, c, c2); ok {
			//if !exists(mc, mergedcomps) {
			mergedcomps = append(mergedcomps, mc)
			//log.Infof("Adding %v", nc)
			mergedmap[i] = true
			mergedmap[0] = true
			//}
		}
	}

	for i, c := range components {
		_, merged := mergedmap[i]
		if !merged {
			log.Debugf("Not merged: %v", c)
			mergedcomps = append(mergedcomps, c)
		}
	}

	log.Infof("Merging components took %s. diff: %d components", time.Since(begin), len(mergedcomps)-len(components))
	return mergedcomps, len(mergedcomps) - len(components)
}

func mergecomponents(g graph.Directed, components [][]graph.Node) ([][]graph.Node, int) {
	var mergedcomps [][]graph.Node
	mergedmap := make(map[int]bool)

	log.Infof("mergecomponents: Merging %d components", len(components))
	begin := time.Now()

	c := components[0]

	for i := 1; i < len(components); i++ {
		c2 := components[i]
		if mc, ok := mergecomponent(g, c, c2); ok {
			//if !exists(mc, mergedcomps) {
			var nc []graph.Node
			for _, n := range mc {
				nc = append(nc, n)
			}
			sort.Sort(ByID(nc))
			mergedcomps = append(mergedcomps, nc)
			//log.Infof("Adding %v", nc)
			mergedmap[i] = true
			mergedmap[0] = true
			//}
		}
	}

	for i, c := range components {
		merged, ok := mergedmap[i]
		if !ok || !merged {
			log.Debugf("Not merged: %v", c)
			mergedcomps = append(mergedcomps, c)
		}
	}

	log.Infof("Merging components took %s. diff: %d components", time.Since(begin), len(mergedcomps)-len(components))
	return mergedcomps, len(mergedcomps) - len(components)
}

func mergecomponentsold(g graph.Graph, components [][]graph.Node) ([][]graph.Node, int) {
	var mergedcomps [][]graph.Node

	begin := time.Now()

	for _, c := range components {
		merged := false
		for i, m := range mergedcomps {
			if mc, ok := mergecomponentold(g, c, m); ok {
				mergedcomps[i] = mc
				merged = true
				break
			}
		}

		if !merged {
			mergedcomps = append(mergedcomps, c)
		}
	}

	log.Infof("Merging components took %s. %d less components", time.Since(begin), len(components)-len(mergedcomps))
	return mergedcomps, len(components) - len(mergedcomps)
}

func mergecomponentmap(g graph.Directed, componentmap map[int][]Component) (map[int][]Component, int, int) {
	var components []Component

	begin := time.Now()

	for _, v := range componentmap {
		components = append(components, v...)
	}

	before := len(components)

	var removed int
	for {
		components = removeDuplicateComponents(components)
		components, removed = mergeComponents(g, components)
		if removed == 0 {
			break
		}
	}

	after := len(components)

	newcomponentmap := make(map[int][]Component)
	var maxidx int

	for _, comp := range components {
		l := len(comp)
		if l > maxidx {
			maxidx = l
		}
		newcomponentmap[l] = append(newcomponentmap[l], comp)
	}

	log.Infof("Merging componentmap took %s. %d less components", time.Since(begin), before-after)
	return newcomponentmap, maxidx, before - after
}
