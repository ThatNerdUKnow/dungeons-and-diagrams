package editor

import "dungeons-and-diagrams/helpers"

func (m *Model) UpdateTable() {
	helpers.UpdateTable(m.Board, m.table)
}
