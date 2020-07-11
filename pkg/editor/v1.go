package editor

import (
	. "github.com/proullon/wikibookgen/api/model"
)

var tmpl = `
% My Book
% Sam Smith

This is my book!

# Chapter One

Chapter one is over.

# Chapter Two

Chapter two has just begun.
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

func (e *V1) Print(l Loader, j Job, w *Wikibook, dest string) error {

	return nil
}
