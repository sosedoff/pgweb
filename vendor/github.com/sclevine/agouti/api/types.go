package api

type Log struct {
	Message   string
	Level     string
	Timestamp int64
}

// A Cookie defines a web cookie
type Cookie struct {
	// Name is the name of the cookie (required)
	Name string `json:"name"`

	// Value is the value of the cookie (required)
	Value string `json:"value"`

	// Path is the path of the cookie (default: "/")
	Path string `json:"path,omitempty"`

	// Domain is the domain of the cookie (default: current page domain)
	Domain string `json:"domain,omitempty"`

	// Secure is set to true for secure cookies (default: false)
	Secure bool `json:"secure,omitempty"`

	// HTTPOnly is set to true for HTTP-Only cookies (default: false)
	HTTPOnly bool `json:"httpOnly,omitempty"`

	// Expiry is the time when the cookie expires
	Expiry float64 `json:"expiry,omitempty"`
}

type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

type Button int

const (
	LeftButton Button = iota
	MiddleButton
	RightButton
)
