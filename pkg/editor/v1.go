package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

var tmpl = `
% My Book ({{.Title}})
% Sam Smith

This is my book!

# Summary

{{range $chapter := .Chapters}}
	- {{.Title}}{{end}}

# Chapter One

Chapter one has just begun.

{{range $chapter := .Chapters}}
# {{.Title}}

{{range $page := $chapter.Articles}}
{{.Content}}
{{end}}
{{end}}
`

type V1 struct {
}

func NewV1() *V1 {
	return &V1{}
}

func (e *V1) Version() string {
	return "1"
}

func (e *V1) Edit(l Loader, j Job, w *Wikibook) error {
	return nil
}

func (e *V1) Print(l Loader, w *Wikibook, folder string) error {

	if len(w.Volumes) == 1 {
		return e.printSingleVolume(l, w, folder)
	}

	return e.printVolumes(l, w, folder)
}

func (e *V1) printSingleVolume(l Loader, w *Wikibook, folder string) error {

	err := e.printVolume(l, w, 0, folder, w.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func (e *V1) printVolumes(l Loader, w *Wikibook, folder string) error {
	return fmt.Errorf("not implemented")
}

func (e *V1) printVolume(l Loader, w *Wikibook, vol int, folder string, name string) error {

	fmt.Printf("%+v\n", w)

	err := e.loadPageContent(l, w.Volumes[vol])
	if err != nil {
		return err
	}

	txtpath := path.Join(folder, fmt.Sprintf("%s.txt", w.Uuid))
	err = e.printWikitxt(w.Volumes[vol], txtpath)
	if err != nil {
		return err
	}

	epubpath := path.Join(folder, fmt.Sprintf("%s.epub", w.Uuid))
	err = e.printEpub(txtpath, epubpath)
	if err != nil {
		return err
	}

	return nil
}

func (e *V1) printWikitxt(v *Volume, dest string) error {
	// create a new template and parse the wikibook format into it.
	t := template.Must(template.New("wikibook").Parse(tmpl))

	f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	err = t.Execute(f, v)
	if err != nil {
		return err
	}

	return nil
}

func (e *V1) printEpub(src string, dst string) error {

	cmd := exec.Command("pandoc", src, "-o", dst)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (e *V1) loadPageContent(l Loader, v *Volume) error {
	var err error

	for ci := range v.Chapters {
		for ai := range v.Chapters[ci].Articles {
			v.Chapters[ci].Articles[ai].Content, err = l.Content(v.Chapters[ci].Articles[ai].Id)
			if err != nil {
				log.Errorf("loadPageContent: cannot load %v: %s", v.Chapters[ci].Articles[ai], err)
				continue
				//return err
			}
		}
	}

	return nil
}
