package editor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/proullon/wikibookgen/pkg/loader"

	. "github.com/proullon/wikibookgen/api/model"
)

func loadWikibook(p string) (*Wikibook, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var w *Wikibook
	err = json.Unmarshal(data, &w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func editorTest(e Editor, t *testing.T) {

	l, err := loader.NewFileLoader("../../samples/mathematiques.dump.json")
	if err != nil {
		t.Errorf("NewFileLoader: %s", err)
	}

	w, err := loadWikibook("../../samples/mathematiques.wikibook.json")
	if err != nil {
		t.Fatalf("cannot load wikibook: %s", err)
	}

	err = e.Edit(l, w)
	if err != nil {
		t.Fatalf("Edit: %s", err)
	}
	fmt.Printf(`%+v\n`, w)

	err = e.Print(l, w, "../../samples/output")
	if err != nil {
		t.Fatalf("Print: %s", err)
	}

}

func TestEditorV1(t *testing.T) {
	e := NewV1()
	editorTest(e, t)
}
