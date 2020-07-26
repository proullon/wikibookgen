package classifier

import (
	"sync"

	"github.com/proullon/workerpool"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
}

type Grapher struct {
	loader  Loader
	wp      *workerpool.WorkerPool
	mu      sync.RWMutex
	maxsize int64
	cursize int64

	graph    *simple.DirectedGraph
	vertices map[int64]*Vertex
	//root     *Vertex
}

func NewV1() (*V1, error) {
	c := &V1{}

	return c, nil
}

func (c *V1) Version() string {
	return "1"
}

func (c *V1) LoadGraph(loader Loader, rootID int64, maxSize int64) (graph.Directed, error) {
	var err error

	l := Grapher{loader: loader, maxsize: maxSize}
	l.vertices = make(map[int64]*Vertex)
	l.graph = simple.NewDirectedGraph()

	root, err := l.loadVertex(rootID)
	if err != nil {
		return nil, err
	}
	trail := make([]int64, 1)
	trail[0] = rootID

	l.wp, err = workerpool.New(l.classify,
		workerpool.WithRetry(5),
		workerpool.WithMaxWorker(200),
		workerpool.WithEvaluationTime(1),
		workerpool.WithSizePercentil(workerpool.LogSizesPercentil),
	)
	if err != nil {
		return nil, err
	}

	for _, id := range root.References {
		p := &payload{
			ID:    id,
			Trail: []int64{rootID},
		}
		l.wp.Feed(p)
	}
	for _, id := range root.Referers {
		p := &payload{
			ID:    id,
			Trail: []int64{rootID},
		}
		l.wp.Feed(p)
	}

	l.wp.Wait()
	l.wp.Stop()

	r := l.wp.AvailableResponses()
	if r > 0 {
		log.Errorf("%d: %d errors", rootID, r)
	}

	return l.graph, nil
}

func (l *Grapher) loadVertex(id int64) (*Vertex, error) {
	v := &Vertex{
		ID:     id,
		Loaded: true,
	}

	outgoing, err := l.loader.LoadOutgoingReferences(id)
	if err != nil {
		return nil, err
	}
	v.References = outgoing

	incoming, err := l.loader.LoadIncomingReferences(id)
	if err != nil {
		return nil, err
	}
	v.Referers = incoming

	return v, nil
}

type payload struct {
	ID    int64
	Trail []int64
}

// classify func
// - check if vertex isn't already loaded
// - load given page if necessary
// - add node to node map
// - add node to graph
// - check if loading node references is pertinent
// - feed reference to workerpool
func (g *Grapher) classify(_payload interface{}) (interface{}, error) {
	p := _payload.(*payload)
	id := p.ID
	trail := p.Trail
	//log.Infof("Classify %-10d from '%+v' please", id, trail)

	if g.maxsize > 0 && g.Size() >= g.maxsize {
		log.Debugf("MaxSize %d reached, ignoring %d", g.maxsize, id)
		return nil, nil
	}

	// check if vertex is in vertices map and loaded. If so you are done
	exist, loaded := g.loaded(id)
	if exist && loaded {
		log.Debugf("Stop. %d exists and loaded", id)
		return nil, nil
	}

	// load node
	v, err := g.loadVertex(id)
	if err != nil {
		log.Errorf("LoadVertex(%d): %s", id, err)
		return nil, err
	}

	// add node to node map
	v = g.addToMap(v)

	// add node to graph
	g.addToGraph(v)

	// check if loading node references is pertinent
	// stop at geodesic distance 2 to avoid infinite loading
	maxdistance := 2
	trail = append(trail, v.ID)
	if len(trail) >= maxdistance {
		//log.Infof("Stop. cause of degree %d", maxdistance)
		return nil, nil
	}

	var tofeed []int64

	for _, r := range v.References {
		// stop if id is already in trail, no loop loading
		for _, t := range trail {
			if t == r {
				log.Debugf("Ignoring %d cause loop detected in trail %+v", t, trail)
				continue
			}
			// check if vertex is in vertices map and loaded. If so you are done
			exist, loaded := g.loaded(r)
			if exist && loaded {
				log.Debugf("Ignoring %d exists and loaded", id)
				continue
			}

			tofeed = append(tofeed, r)
		}
	}

	for _, r := range v.Referers {
		// stop if id is already in trail, no loop loading
		for _, t := range trail {
			if t == r {
				log.Debugf("Ignoring %d cause loop detected in trail %+v", t, trail)
				continue
			}
			exist, loaded := g.loaded(r)
			if exist && loaded {
				log.Debugf("Ignoring %d exists and loaded", id)
				continue
			}

			tofeed = append(tofeed, r)
		}
	}

	// Feed ref in workerpool
	for _, id := range tofeed {
		p := &payload{
			ID:    id,
			Trail: trail,
		}
		g.wp.Feed(p)
	}

	return nil, nil
}

func (g *Grapher) Size() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.cursize
}

// addToGraph needs to go through all references and referers, try to find them in vertices map
//  - if found, add Vertex instance to IncomingEdges and OutgoingEdges
//  - if not found, insert unloaded Vertex instance in vertices map to add connection
// in all case, add resulting Vertex instance to own IncomingEdges and OutgoingEdges slices
func (g *Grapher) addToGraph(v *Vertex) {
	g.mu.Lock()
	defer g.mu.Unlock()
	/*
		for _, refID := range v.References {
			ref, exist := g.vertices[refID]
			if !exist {
				ref = &Vertex{ID: refID}
				g.vertices[refID] = ref
			}

			ref.IncomingEdges = append(ref.IncomingEdges, v)
			v.OutgoingEdges = append(v.OutgoingEdges, ref)
		}

		for _, refID := range v.Referers {
			ref, exist := g.vertices[refID]
			if !exist {
				ref = &Vertex{ID: refID}
				g.vertices[refID] = ref
			}

			ref.OutgoingEdges = append(ref.OutgoingEdges, v)
			v.IncomingEdges = append(v.IncomingEdges, ref)
		}
	*/
	//log.Infof("NewNode(%d)", v.ID)
	n := g.getOrNew(v.ID)

	for _, refID := range v.References {
		if refID == v.ID {
			continue // avoid self edge
		}
		if g.maxsize > 0 && g.cursize >= g.maxsize {
			return
		}

		to := g.getOrNew(refID)
		log.Debugf("NewEdge {%d, %d}", n.ID(), to.ID())
		g.graph.SetEdge(g.graph.NewEdge(n, to))
	}

	for _, refID := range v.Referers {
		if refID == v.ID {
			continue // avoid self edge
		}
		if g.maxsize > 0 && g.cursize >= g.maxsize {
			return
		}
		from := g.getOrNew(refID)
		log.Debugf("NewEdge {%d, %d}", from.ID(), n.ID())
		g.graph.SetEdge(g.graph.NewEdge(from, n))
	}
}

func (g *Grapher) getOrNew(id int64) graph.Node {

	n := g.graph.Node(id)
	if n == nil {
		n = NewNode(id)
		g.graph.AddNode(n)
		g.cursize++
	}

	return n
}

// add to map, but if unloaded instance already exist,
// copy References and Referers then set loaded to true
// then return existing instance
func (g *Grapher) addToMap(v *Vertex) *Vertex {
	g.mu.Lock()
	defer g.mu.Unlock()

	/*
		original, exist := g.vertices[v.ID]
		if exist {
			original.Loaded = true
			original.References = original.References
			original.Referers = v.Referers
			return original
		}

	*/
	g.vertices[v.ID] = v
	return v
}

func (g *Grapher) loaded(id int64) (exist, loaded bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, exist = g.vertices[id]
	if !exist {
		return false, false
	}

	return exist, true
}
