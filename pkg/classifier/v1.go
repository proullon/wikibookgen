package classifier

import (
	"sync"

	"github.com/proullon/workerpool"
	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
	loader Loader
}

type Grapher struct {
	loader Loader
	wp     *workerpool.WorkerPool
	mu     sync.RWMutex

	vertices map[int]*Vertex
	root     *Vertex
}

func NewV1(l Loader) (*V1, error) {
	c := &V1{
		loader: l,
	}

	return c, nil
}

func (c *V1) LoadGraph(rootID int) (*Vertex, map[int]*Vertex, error) {
	var err error

	l := Grapher{loader: c.loader}

	l.vertices = make(map[int]*Vertex)
	l.root, err = l.loadVertex(rootID)
	if err != nil {
		return nil, nil, err
	}
	trail := make([]int, 1)
	trail[0] = rootID

	l.wp, err = workerpool.New(l.classify,
		workerpool.WithRetry(5),
		workerpool.WithMaxWorker(100),
		workerpool.WithEvaluationTime(1),
		workerpool.WithSizePercentil(workerpool.LogSizesPercentil),
	)
	if err != nil {
		return nil, nil, err
	}

	for _, id := range l.root.References {
		p := &payload{
			ID:    id,
			Trail: []int{rootID},
		}
		l.wp.Feed(p)
	}
	for _, id := range l.root.Referers {
		p := &payload{
			ID:    id,
			Trail: []int{rootID},
		}
		l.wp.Feed(p)
	}

	l.wp.Wait()
	l.wp.Stop()

	r := l.wp.AvailableResponses()
	if r > 0 {
		log.Errorf("%d: %d errors", rootID, r)
	}

	return l.root, l.vertices, nil
}

func (l *Grapher) loadVertex(id int) (*Vertex, error) {
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
	ID    int
	Trail []int
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

	// check if vertex is in vertices map and loaded. If so you are done
	exist, loaded := g.exist(id)
	if exist && loaded {
		log.Debugf("Stop. %d exists and loaded", id)
		return nil, nil
	}

	// load node
	v, err := g.loadVertex(id)
	if err != nil {
		return nil, err
	}

	// add node to node map
	v = g.addToMap(v)

	// add node to graph
	g.addToGraph(v)

	// check if loading node references is pertinent
	// stop at 2nd degree to avoid infinite loading
	maxdegree := 2
	if len(trail) == maxdegree {
		//log.Infof("Stop. cause of degree %d", maxdegree)
		return nil, nil
	}
	trail = append(trail, v.ID)

	var tofeed []int

	for _, r := range v.References {
		// stop if id is already in trail, no loop loading
		for _, t := range trail {
			if t == r {
				log.Infof("Ignoring %d cause loop detected in trail %+v", t, trail)
				continue
			}

			tofeed = append(tofeed, r)
		}
	}

	for _, r := range v.Referers {
		// stop if id is already in trail, no loop loading
		for _, t := range trail {
			if t == r {
				log.Infof("Ignoring %d cause loop detected in trail %+v", t, trail)
				continue
			}

			tofeed = append(tofeed, r)
		}
	}

	// Feed ref in workerpool
	go func(trail []int, tofeed []int) {
		for _, id := range tofeed {
			p := &payload{
				ID:    id,
				Trail: trail,
			}

			g.wp.Feed(p)
		}
	}(trail, tofeed)

	return nil, nil
}

// addToGraph needs to go through all references and referers, try to find them in vertices map
//  - if found, add Vertex instance to IncomingEdges and OutgoingEdges
//  - if not found, insert unloaded Vertex instance in vertices map to add connection
// in all case, add resulting Vertex instance to own IncomingEdges and OutgoingEdges slices
func (g *Grapher) addToGraph(v *Vertex) {
	g.mu.Lock()
	defer g.mu.Unlock()

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
}

// add to map, but if unloaded instance already exist,
// copy References and Referers then set loaded to true
// then return existing instance
func (g *Grapher) addToMap(v *Vertex) *Vertex {
	g.mu.Lock()
	defer g.mu.Unlock()

	original, exist := g.vertices[v.ID]
	if exist {
		original.Loaded = true
		original.References = original.References
		original.Referers = v.Referers
		return original
	}

	g.vertices[v.ID] = v
	return v
}

func (g *Grapher) exist(id int) (exist, loaded bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	v, exist := g.vertices[id]
	if !exist {
		return false, false
	}

	return exist, v.Loaded
}
