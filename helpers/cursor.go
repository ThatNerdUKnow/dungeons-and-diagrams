package helpers

type Cursor1D struct {
	i     int
	i_max int
}

func NewCursor1D(max int) Cursor1D {
	return Cursor1D{i: 0, i_max: max}
}

func (c *Cursor1D) Inc() {
	tmp := (c.i + 1) % c.i_max
	if tmp < 0 {
		tmp += c.i_max
	}
	c.i = tmp
}

func (c *Cursor1D) Dec() {
	tmp := (c.i - 1) % c.i_max
	if tmp < 0 {
		tmp += c.i_max
	}
	c.i = tmp
}

type Cursor2D struct {
	X Cursor1D
	Y Cursor1D
}

func NewCursor2D(x_max int, y_max int) Cursor2D {
	return Cursor2D{X: NewCursor1D(x_max), Y: NewCursor1D(y_max)}
}

func (c Cursor2D) Coords() (int, int) {
	return c.CoordsOffset(0, 0)
}

func (c Cursor2D) CoordsOffset(x int, y int) (int, int) {
	return int(c.X.i) + x, c.Y.i + y
}
