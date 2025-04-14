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

func (c Cursor2D) Coords() [2]int {
	return [2]int{int(c.X.i), int(c.Y.i)}
}
