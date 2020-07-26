package editor

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

var tmpl = `
% My Book ({{.Title}})
% Sam Smith

# Chapter One

Chapter one has just begun.

{{range $chapter := .Chapters}}
# {{.Title}}

{{range $page := $chapter.Articles}}
{{.Content}}
{{end}}
{{end}}
`

var titletmpl = `---
title:
- type: main
  text: {{.Title}}
- type: subtitle
	text: Wikipedia article collection
creator:
- role: author
	text: wikibookgen.org
- role: editor
	text: wikibookgen.org
publisher:  wikibookgen.org
rights: © 2020 wikibookgen, CC BY-NC
ibooks:
	version: 1.3.4
`

var chaptertmpl = `
= {{.Title}} =

This is a chapter beginning.

Some text lol.

{{range $page := .Articles}}
== {{.Title}} ==

Some text about {{.Title}}.
{{end}}

Some more text.
`

/*
{{range $page := .Articles}}
{{.Content}}
{{end}}
`
*/

type i18n struct {
	VolumeTitleFmt     string
	ChapterTitleFmt    string
	File               string
	SeeAlso            string
	NotesAndReferences string
}

var langmap map[string]i18n = map[string]i18n{
	"en": {
		VolumeTitleFmt:     "Tour of %s",
		ChapterTitleFmt:    "Chapter %d: %s",
		File:               "File",
		SeeAlso:            "See also",
		NotesAndReferences: "References",
	},
	"fr": {
		VolumeTitleFmt:     "%s en long, en large et en travers",
		ChapterTitleFmt:    "Chapter %d: %s",
		File:               "Fichier",
		SeeAlso:            "Voir aussi",
		NotesAndReferences: "Notes et références",
	},
}

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

	err := e.loadPageContent(l, w.Volumes[vol], w.Language)
	if err != nil {
		return err
	}

	//	txtpath := path.Join(folder, fmt.Sprintf("%s.txt", w.Uuid))
	err = e.printWikitxt(w.Volumes[vol], w.Language, folder, w.Uuid)
	if err != nil {
		return err
	}
	/*
		epubpath := path.Join(folder, fmt.Sprintf("%s.epub", w.Uuid))
		err = e.printEpub(txtpath, epubpath)
		if err != nil {
			return err
		}

		pdfpath := path.Join(folder, fmt.Sprintf("%s.pdf", w.Uuid))
		err = e.printPDF(txtpath, pdfpath)
		if err != nil {
			return err
		}
	*/
	return nil
}

func (e *V1) printWikitxt(v *Volume, lang, folder, id string) error {
	// create a new template and parse the wikibook format into it.
	title := template.Must(template.New("wikibook-title").Parse(titletmpl))
	chapter := template.Must(template.New("wikibook-chapter").Parse(chaptertmpl))

	titlename := fmt.Sprintf("%s.txt", id)
	titlepath := path.Join(folder, titlename)
	f, err := os.OpenFile(titlepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	err = title.Execute(f, v)
	if err != nil {
		return err
	}
	log.Infof("Wikibook volume '%s' written in %s", v.Title, titlepath)

	var texfiles []string
	var dest string
	for i, c := range v.Chapters {

		dest = path.Join(folder, fmt.Sprintf("%s-%d.txt", id, i+1))
		f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Errorf("cannot open file %s: %s", dest, err)
			continue
		}
		defer f.Close()
		err = chapter.Execute(f, c)
		if err != nil {
			log.Errorf("cannot execute chapter template for %s: %s", c.Title, err)
			continue
		}
		log.Infof("Wikibook volume chapter written in %s", dest)

		texname := fmt.Sprintf("%s-%d.tex", id, i+1)
		tex := path.Join(folder, texname)
		err = e.convertTex(dest, tex)
		if err != nil {
			log.Errorf("cannot convert %s to tex: %s", dest, err)
			continue
		}

		err = e.downloadFiles(c, lang, folder)
		if err != nil {
			log.Errorf("cannot download files for %s: %s", dest, err)
			continue
		}

		texfiles = append(texfiles, texname)
		/*
			if len(texfiles) == 2 {
				break
			}
		*/
	}

	dst := path.Join(folder, fmt.Sprintf("%s.pdf", id))
	log.Infof("generating pdf %s", dst)

	args := []string{
		`--toc`,
		`--latex-engine=xelatex`,
		`-o`,
		dst,
		fmt.Sprintf("--epub-metadata=%s", titlename),
	}
	args = append(args, texfiles...)
	log.Info(args)

	cmd := exec.Command("pandoc", args...)
	cmd.Dir = folder
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (e *V1) convertTex(src, dst string) error {
	log.Infof("Writing tex chapter in %s", dst)

	cmd := exec.Command("pandoc", `-f`, `mediawiki`, src, "-o", dst)
	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Infof("Wikibook tex chapter written in %s", dst)
	return nil
}

func (e *V1) printEpub(src string, dst string) error {
	log.Infof("Writing epub in %s", dst)

	cmd := exec.Command("pandoc", "--toc", src, "-o", dst)
	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Infof("Wikibook epub written in %s", dst)
	return nil
}

func (e *V1) printPDF(src string, dst string) error {

	log.Infof("Writing pdf in %s", dst)

	cmd := exec.Command("pandoc", "--latex-engine=xelatex", "--toc", src, "-o", dst)
	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Infof("Wikibook pdf written in %s", dst)
	return nil
}

func (e *V1) loadPageContent(l Loader, v *Volume, lang string) error {
	var err error

	for ci := range v.Chapters {
		for ai := range v.Chapters[ci].Articles {
			v.Chapters[ci].Articles[ai].Content, err = l.Content(v.Chapters[ci].Articles[ai].Id)
			if err != nil {
				log.Errorf("loadPageContent: cannot load %v: %s", v.Chapters[ci].Articles[ai], err)
				continue
				//return err
			}
			v.Chapters[ci].Articles[ai].Content = removeSeeAlso(v.Chapters[ci].Articles[ai].Content)
			v.Chapters[ci].Articles[ai].Content = removeVideo(v.Chapters[ci].Articles[ai].Content, lang)
			v.Chapters[ci].Articles[ai].Content = removeLink(v.Chapters[ci].Articles[ai].Content, lang)
			v.Chapters[ci].Articles[ai].Content = removeNotesAndReferences(v.Chapters[ci].Articles[ai].Content, lang)
		}
	}

	return nil
}

// See: https://commons.wikimedia.org/wiki/Commons:FAQ#What_are_the_strangely_named_components_in_file_paths.3F
func (e *V1) downloadFiles(c *Chapter, lang, folder string) error {

	for _, a := range c.Articles {

		re := regexp.MustCompile(`\[\[(?s)(.*?)\]\]`)
		sub := re.FindAllString(a.Content, -1)
		for _, s := range sub {
			// [[Fichier:Nombrepremier 2017.png|vignette|Le nombre [[7 (nombre)|7]]
			//entry := s
			s = strings.TrimPrefix(strings.TrimSuffix(s, "]]"), "[[")
			s = strings.Split(s, "|")[0]
			if !strings.HasPrefix(s, langmap[lang].File+":") && !strings.HasPrefix(s, "Image:") {
				continue
			}
			s = strings.Split(s, ":")[1]
			filename := s
			log.Infof("need to download %s", s)
			/*
				if strings.HasSuffix(s, "webm") {
					log.Infof("Removing from content '%s'", entry)
					strings.ReplaceAll(c.Articles[i].Content, entry, "")
					continue
				}
			*/

			s = strings.TrimSpace(s)
			s = strings.Replace(s, " ", "_", -1)

			sum := md5.Sum([]byte(s))
			hash := hex.EncodeToString(sum[:])

			url := fmt.Sprintf("https://upload.wikimedia.org/wikipedia/commons/%s/%s/%s", hash[:1], hash[:2], s)

			if u, err := LocateWikiFile(s, lang); err == nil {
				url = u
			} else {
				log.Infof("LocateWikiFile: %s", err)
			}

			dst := path.Join(folder, filename)
			cmd := exec.Command("wget", url, "-O", dst)
			log.Infof("running %v", cmd)
			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeLink(c, lang string) string {
	re := regexp.MustCompile(`\[\[(?s)(.*?)\]\]`)
	sub := re.FindAllString(c, -1)
	for _, s := range sub {
		entry := s
		s = strings.TrimPrefix(strings.TrimSuffix(s, "]]"), "[[")

		if strings.HasPrefix(s, FilePrefix(lang)+":") || strings.HasPrefix(s, "Image:") {
			continue
		}

		split := strings.Split(s, "|")
		if len(split) == 1 {
			c = strings.ReplaceAll(c, entry, split[0])
		}
		if len(split) >= 2 {
			c = strings.ReplaceAll(c, entry, split[1])
			//log.Infof("ReplaceAll '%s' with '%s'", entry, split[1])
		}
	}

	return c
}

func removeVideo(c, lang string) string {
	//re := regexp.MustCompile(`\[\[.(.*?)\]\]`)
	re := regexp.MustCompile(`\[\[(?s)(.*?)\]\]`)
	sub := re.FindAllString(c, -1)
	for _, s := range sub {
		entry := s
		if strings.Contains(entry, "lason") {
			log.Warnf("Found '%s' !!!", entry)
		}
		s = strings.TrimPrefix(strings.TrimSuffix(s, "]]"), "[[")
		s = strings.Split(s, "|")[0]
		if !strings.HasPrefix(s, FilePrefix(lang)+":") {
			continue
		}
		s = strings.Split(s, ":")[1]
		//		if strings.HasSuffix(s, "webm") {
		//		log.Infof("Removing from content '%s'", entry)
		c = strings.ReplaceAll(c, entry, "")
		//		}
	}

	return c
}

func removeSeeAlso(c string) string {
	possibletext := []string{
		`== Voir aussi ==`,
		`==Voir aussi==`,
		`==See also==`,
		`== See also ==`,
	}

	for _, t := range possibletext {
		i := strings.Index(c, t)
		if i > 0 {
			c = c[:i]
			return c
		}
	}

	return c
}

func removeNotesAndReferences(c string, lang string) string {
	i := strings.Index(c, fmt.Sprintf("== %s ==", NotesAndReferences(lang)))
	if i > 0 {
		c = c[:i]
		return c
	}

	i = strings.Index(c, fmt.Sprintf("==%s==", NotesAndReferences(lang)))
	if i > 0 {
		c = c[:i]
		return c
	}

	return c
}

func FilePrefix(lang string) string {
	return langmap[lang].File
}

func NotesAndReferences(lang string) string {
	return langmap[lang].NotesAndReferences
}
