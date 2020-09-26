package orderer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"

	. "github.com/proullon/wikibookgen/api/model"
)

type V1 struct {
}

func NewV1() *V1 {
	o := &V1{}

	return o
}

func (o *V1) Version() string {
	return "1"
}

// Order v1 will not actually order, just copy clusters into chapters
func (o *V1) Order(l Loader, j Job, g graph.Directed, clusters *Cluster) (*Wikibook, error) {
	data, err := json.Marshal(clusters)
	if err != nil {
		log.Errorf("Cannot marshal clusters: %s", err)
		return nil, err
	}
	err = ioutil.WriteFile(fmt.Sprintf("/tmp/wikibookgen/clusters.%s.json", j.Subject), data, 0666)
	if err != nil {
		log.Errorf("cannot write clusters: %s", err)
		return nil, err
	}

	wikibook, err := o.GenTour(l, j, g, clusters)
	if err != nil {
		return nil, err
	}

	for _, v := range wikibook.Volumes {
		o.OrderChapters(g, v)
		for _, c := range v.Chapters {
			o.OrderChapter(g, c)
		}
	}

	return wikibook, nil
}

func (o *V1) GenTour(l Loader, j Job, g graph.Directed, cluster *Cluster) (*Wikibook, error) {
	wikibook := &Wikibook{
		Language: j.Language,
		Subject:  j.Subject,
		Title:    j.Subject,
		Model:    j.Model,
	}

	v, pages, err := o.VolumeThis(l, j, g, cluster)
	if err != nil {
		return nil, err
	}
	wikibook.Volumes = append(wikibook.Volumes, v)
	wikibook.Pages += int64(pages)

	return wikibook, nil
}

func (o *V1) VolumeThis(l Loader, j Job, g graph.Directed, cluster *Cluster) (*Volume, int, error) {
	v := &Volume{Title: j.Subject}
	var pages int

	for _, cluster := range cluster.Subclusters {
		c, err := o.ChapterThis(l, g, cluster)
		if err != nil {
			log.Errorf(`Cannot create chapter from cluster %+v: %s`, cluster, err)
			continue
		}
		pages += len(c.Articles)

		v.Chapters = append(v.Chapters, c)
	}

	return v, pages, nil
}

func (o *V1) ChapterThis(l Loader, g graph.Directed, cluster *Cluster) (*Chapter, error) {
	c := &Chapter{}

	center := o.Center(g, cluster)
	title, err := l.Title(center)
	if err != nil {
		return nil, fmt.Errorf("ChapterTitle(%d): %s", center, err)
		//title = fmt.Sprintf("ChapterTitle(%d): %s", center, err)
	}
	c.Title = title

	for _, n := range cluster.Members {
		p, err := o.PageThis(l, n)
		if err != nil {
			return nil, err
		}
		c.Articles = append(c.Articles, p)
	}

	return c, nil
}

func (o *V1) PageThis(l Loader, n graph.Node) (*Page, error) {
	p := &Page{Id: n.ID()}
	title, err := l.Title(n.ID())
	if err != nil {
		//return nil, fmt.Errorf("PageTitle(%d): %s", n.ID(), err)
		title = fmt.Sprintf("PageTitle(%d): %s", n.ID(), err)
	}
	p.Title = title
	return p, nil
}

func (o *V1) Center(g graph.Directed, cluster *Cluster) int64 {
	bvalues := cluster.Members.Betweenness(g)
	var center int64
	var centerb float64
	for id, b := range bvalues {
		if b > centerb {
			center = id
			centerb = b
		}
	}
	return center
}

func (o *V1) OrderChapters(g graph.Directed, v *Volume) {
	sort.Sort(BySumReference{g: g, c: v.Chapters})
}

func (o *V1) OrderChapter(g graph.Directed, c *Chapter) {
	sort.Sort(ByReference{g: g, p: c.Articles})
}

type ByReference struct {
	p []*Page
	g graph.Graph
}

func (a ByReference) Len() int      { return len(a.p) }
func (a ByReference) Swap(i, j int) { a.p[i], a.p[j] = a.p[j], a.p[i] }
func (a ByReference) Less(i, j int) bool {

	ni := a.p[i].Id
	nj := a.p[j].Id
	if a.g.Edge(ni, nj) != nil && a.g.Edge(nj, ni) == nil {
		return false
	}

	return true
}

type BySumReference struct {
	c []*Chapter
	g graph.Graph
}

func (a BySumReference) Len() int      { return len(a.c) }
func (a BySumReference) Swap(i, j int) { a.c[i], a.c[j] = a.c[j], a.c[i] }
func (a BySumReference) Less(i, j int) bool {

	var refcountj, refcounti int

	for _, cj := range a.c[j].Articles {
		for _, ci := range a.c[i].Articles {

			ni := a.g.Node(ci.Id)
			nj := a.g.Node(cj.Id)
			if a.g.Edge(ni.ID(), nj.ID()) != nil {
				refcounti++
			}
			if a.g.Edge(nj.ID(), ni.ID()) != nil {
				refcountj++
			}
		}
	}

	log.Infof("Chapter '%s' has %d references to '%s' / '%s' %d reference to '%s'", a.c[i].Title, refcounti, a.c[j].Title, a.c[j].Title, refcountj, a.c[i].Title)
	return refcounti < refcountj
}
