package agouti

type Tap int

const (
	SingleTap Tap = iota
	DoubleTap
	LongTap
)

func (t Tap) String() string {
	switch t {
	case SingleTap:
		return "tap"
	case DoubleTap:
		return "double tap"
	case LongTap:
		return "long tap"
	}
	return "perform tap"
}

type Touch int

const (
	HoldFinger Touch = iota
	ReleaseFinger
	MoveFinger
)

func (t Touch) String() string {
	switch t {
	case HoldFinger:
		return "hold finger down"
	case ReleaseFinger:
		return "release finger"
	case MoveFinger:
		return "move finger"
	}
	return "perform touch"
}

type Button int

const (
	LeftButton Button = iota
	MiddleButton
	RightButton
)

func (b Button) String() string {
	switch b {
	case LeftButton:
		return "left mouse button"
	case MiddleButton:
		return "middle mouse button"
	case RightButton:
		return "right mouse button"
	}
	return "unknown"
}

type Click int

const (
	SingleClick Click = iota
	HoldClick
	ReleaseClick
)

func (c Click) String() string {
	switch c {
	case SingleClick:
		return "single click"
	case HoldClick:
		return "hold"
	case ReleaseClick:
		return "release"
	}
	return "unknown"
}
