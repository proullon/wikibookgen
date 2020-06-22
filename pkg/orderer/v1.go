package orderer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"

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

	/*
		wikibook := &Wikibook{
			Subject:  j.Subject,
			Title:    j.Subject,
			Language: j.Language,
			Pages:    int64(len(clusters.Members)),
		}
		volume := &Volume{
			Title: j.Subject,
		}
		wikibook.Volumes = []*Volume{volume}

		for _, cluster := range clusters.Subclusters {
			chapter := &Chapter{}
			bvalues := cluster.Members.Betweenness(g)
			var center int64
			var centerb float64
			for id, b := range bvalues {
				if b > centerb {
					center = id
					centerb = b
				}
			}
			chapter.Title = "???"
			if center > 0 {
				chapter.Title, err = l.Title(center)
				if err != nil {
					return nil, fmt.Errorf("Order Chapter Title(%d): %s", center, err)
				}
				log.Infof("Component center is %d with %f !", center, centerb)
			}
			for _, n := range cluster.Members {
				article, err := o.PageThis(l, n)
				if err != nil {
					return nil, err
				}
				chapter.Articles = append(chapter.Articles, article)
			}

			volume.Chapters = append(volume.Chapters, chapter)
		}
	*/

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
		Subject: j.Subject,
		Title:   j.Subject,
	}

	v, err := o.VolumeThis(l, j, g, cluster)
	if err != nil {
		return nil, err
	}
	wikibook.Volumes = append(wikibook.Volumes, v)

	return wikibook, nil
}

func (o *V1) VolumeThis(l Loader, j Job, g graph.Directed, cluster *Cluster) (*Volume, error) {
	v := &Volume{Title: j.Subject}

	for _, cluster := range cluster.Subclusters {
		c, err := o.ChapterThis(l, g, cluster)
		if err != nil {
			return nil, err
		}

		v.Chapters = append(v.Chapters, c)
	}

	return v, nil
}

func (o *V1) ChapterThis(l Loader, g graph.Directed, cluster *Cluster) (*Chapter, error) {
	c := &Chapter{}

	center := o.Center(g, cluster)
	title, err := l.Title(center)
	if err != nil {
		//return nil, fmt.Errorf("ChapterTitle(%d): %s", center, err)
		title = fmt.Sprintf("ChapterTitle(%d): %s", center, err)
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
}

func (o *V1) OrderChapter(g graph.Directed, c *Chapter) {
	for {
		modified := false

		for i := range c.Articles {
			if i == 0 {
				continue
			}

			g1 := g.Node(c.Articles[i].Id)
			g2 := g.Node(c.Articles[i-1].Id)
			ptoi := path.YenKShortestPaths(g, 10, g1, g2)
			ptoim := path.YenKShortestPaths(g, 10, g2, g1)
			// if there is more path to article, put it first in chapter
			if len(ptoi) < len(ptoim) {
				log.Infof("Swapping %v (%d) and %v (%d)", c.Articles[i], len(ptoi), c.Articles[i-1], len(ptoim))
				s := c.Articles[i]
				c.Articles[i] = c.Articles[i-1]
				c.Articles[i-1] = s
				modified = true
			}
		}

		if !modified {
			return
		}
	}
}
