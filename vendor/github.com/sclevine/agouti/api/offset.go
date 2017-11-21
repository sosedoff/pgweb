package api

type Offset interface {
	x() (x int, present bool)
	y() (y int, present bool)
	position() (x int, y int)
}

type XYOffset struct {
	X int
	Y int
}

func (o XYOffset) x() (x int, present bool) {
	return o.X, true
}

func (o XYOffset) y() (y int, present bool) {
	return o.Y, true
}

func (o XYOffset) position() (x int, y int) {
	return o.X, o.Y
}

type XOffset int

func (o XOffset) x() (x int, present bool) {
	return int(o), true
}

func (XOffset) y() (y int, present bool) {
	return 0, false
}

func (o XOffset) position() (x int, y int) {
	return int(o), 0
}

type YOffset int

func (YOffset) x() (x int, present bool) {
	return 0, false
}

func (o YOffset) y() (y int, present bool) {
	return int(o), true
}

func (o YOffset) position() (x int, y int) {
	return 0, int(o)
}
