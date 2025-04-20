package model

import (
	"dungeons-and-diagrams/editor"
	"dungeons-and-diagrams/solved"
)

const (
	editmode int = iota
	solvedmode
)

type Model struct {
	mode   int
	edit   editor.Model
	solved *solved.SolvedModel
}

func New() Model {
	mode := editmode
	editor := editor.Default()
	return Model{mode: mode, edit: editor}
}
