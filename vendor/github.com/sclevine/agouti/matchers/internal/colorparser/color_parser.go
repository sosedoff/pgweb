package colorparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Color struct {
	R, G, B uint8
	A       float64
}

func (c Color) String() string {
	return fmt.Sprintf("Color{R:%d, G:%d, B:%d, A:%.2f}", c.R, c.G, c.B, c.A)
}

var (
	shortRGBHexRE    = regexp.MustCompile(`^#([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])$`)
	longRGBHexRE     = regexp.MustCompile(`^#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`)
	rgbIntegerRE     = regexp.MustCompile(`^rgb\(\s*([-0-9]+)\s*,\s*([-0-9]+)\s*,\s*([-0-9]+)\s*\)$`)
	rgbPercentageRE  = regexp.MustCompile(`^rgb\(\s*([-0-9.]+)%\s*,\s*([-0-9.]+)%\s*,\s*([-0-9.]+)%\s*\)$`)
	rgbaIntegerRE    = regexp.MustCompile(`^rgba\(\s*(-?[0-9]+)\s*,\s*(-?[0-9]+)\s*,\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)\s*\)$`)
	rgbaPercentageRE = regexp.MustCompile(`^rgba\(\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)\s*\)$`)
	hslRE            = regexp.MustCompile(`^hsl\(\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*\)$`)
	hslaRE           = regexp.MustCompile(`^hsla\(\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)\s*\)$`)
)

func ParseCSSColor(color string) (Color, error) {
	color = strings.Trim(color, " ")
	rgba, ok := colorLookup[color]
	if ok {
		return rgba, nil
	}
	switch {
	case shortRGBHexRE.MatchString(color):
		return parseShortRGBHex(color)
	case longRGBHexRE.MatchString(color):
		return parseLongRGBHex(color)
	case rgbIntegerRE.MatchString(color):
		return parseRGBInteger(color)
	case rgbPercentageRE.MatchString(color):
		return parseRGBPercentage(color)
	case rgbaIntegerRE.MatchString(color):
		return parseRGBAInteger(color)
	case rgbaPercentageRE.MatchString(color):
		return parseRGBAPercentage(color)
	case hslRE.MatchString(color):
		return parseHSL(color)
	case hslaRE.MatchString(color):
		return parseHSLA(color)
	default:
		return Color{}, errors.New("unparseable color")
	}
}

func parseShortRGBHex(color string) (Color, error) {
	components := shortRGBHexRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb hex")
	}
	r, err := strconv.ParseUint(components[1]+components[1], 16, 32)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseUint(components[2]+components[2], 16, 32)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseUint(components[3]+components[3], 16, 32)
	if err != nil {
		return Color{}, err
	}
	return Color{uint8(r), uint8(g), uint8(b), 1.0}, nil
}

func parseLongRGBHex(color string) (Color, error) {
	components := longRGBHexRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb hex")
	}
	r, err := strconv.ParseUint(components[1], 16, 32)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseUint(components[2], 16, 32)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseUint(components[3], 16, 32)
	if err != nil {
		return Color{}, err
	}
	return Color{uint8(r), uint8(g), uint8(b), 1.0}, nil
}

func parseRGBInteger(color string) (Color, error) {
	components := rgbIntegerRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb")
	}
	r, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseInt(components[3], 10, 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		clamp255(r),
		clamp255(g),
		clamp255(b),
		1.0,
	}, nil
}

func parseRGBPercentage(color string) (Color, error) {
	components := rgbPercentageRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb percentage")
	}
	r, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		round255(r / 100.0 * 255),
		round255(g / 100.0 * 255),
		round255(b / 100.0 * 255),
		1.0,
	}, nil
}

func parseRGBAInteger(color string) (Color, error) {
	components := rgbaIntegerRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid rgba")
	}
	r, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseInt(components[3], 10, 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		clamp255(r),
		clamp255(g),
		clamp255(b),
		clamp1(a),
	}, nil
}

func parseRGBAPercentage(color string) (Color, error) {
	components := rgbaPercentageRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid rgb percentage")
	}
	r, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		round255(r / 100.0 * 255),
		round255(g / 100.0 * 255),
		round255(b / 100.0 * 255),
		clamp1(a),
	}, nil
}

func parseHSL(color string) (Color, error) {
	components := hslRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid hsl percentage")
	}
	h, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	s, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	l, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}

	return colorFromHSL(h, s, l, 1.0), nil
}

func parseHSLA(color string) (Color, error) {
	components := hslaRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid hsl percentage")
	}
	h, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	s, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	l, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}

	return colorFromHSL(h, s, l, clamp1(a)), nil
}

func clamp255(value int64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}

	return uint8(value)
}

func round255(value float64) uint8 {
	value = value + 0.5 //round!
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}

	return uint8(value)
}

func clamp1(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func colorFromHSL(hDegrees int64, s float64, l float64, a float64) Color {
	//see http://www.w3.org/TR/css3-color
	h := float64((((hDegrees % 360) + 360) % 360)) / 360.0
	s = clamp1(s / 100.0)
	l = clamp1(l / 100.0)

	var m2 float64
	if l <= 0.5 {
		m2 = l * (s + 1)
	} else {
		m2 = l + s - l*s
	}
	m1 := l*2 - m2

	return Color{
		R: hueToRGB(m1, m2, h+1.0/3.0),
		G: hueToRGB(m1, m2, h),
		B: hueToRGB(m1, m2, h-1.0/3.0),
		A: a,
	}
}

func hueToRGB(m1 float64, m2 float64, h float64) uint8 {
	if h < 0 {
		h = h + 1
	}
	if h > 1 {
		h = h - 1
	}
	if h*6 < 1 {
		return uint8(255 * (m1 + (m2-m1)*h*6.0))
	}
	if h*2 < 1 {
		return uint8(255 * m2)
	}
	if h*3 < 2 {
		return uint8(255 * (m1 + (m2-m1)*(2.0/3.0-h)*6.0))
	}
	return uint8(255 * m1)
}

var colorLookup = map[string]Color{
	"aliceblue":            {240, 248, 255, 1.0},
	"antiquewhite":         {250, 235, 215, 1.0},
	"aqua":                 {0, 255, 255, 1.0},
	"aquamarine":           {127, 255, 212, 1.0},
	"azure":                {240, 255, 255, 1.0},
	"beige":                {245, 245, 220, 1.0},
	"bisque":               {255, 228, 196, 1.0},
	"black":                {0, 0, 0, 1.0},
	"blanchedalmond":       {255, 235, 205, 1.0},
	"blue":                 {0, 0, 255, 1.0},
	"blueviolet":           {138, 43, 226, 1.0},
	"brown":                {165, 42, 42, 1.0},
	"burlywood":            {222, 184, 135, 1.0},
	"cadetblue":            {95, 158, 160, 1.0},
	"chartreuse":           {127, 255, 0, 1.0},
	"chocolate":            {210, 105, 30, 1.0},
	"coral":                {255, 127, 80, 1.0},
	"cornflowerblue":       {100, 149, 237, 1.0},
	"cornsilk":             {255, 248, 220, 1.0},
	"crimson":              {220, 20, 60, 1.0},
	"cyan":                 {0, 255, 255, 1.0},
	"darkblue":             {0, 0, 139, 1.0},
	"darkcyan":             {0, 139, 139, 1.0},
	"darkgoldenrod":        {184, 134, 11, 1.0},
	"darkgray":             {169, 169, 169, 1.0},
	"darkgreen":            {0, 100, 0, 1.0},
	"darkgrey":             {169, 169, 169, 1.0},
	"darkkhaki":            {189, 183, 107, 1.0},
	"darkmagenta":          {139, 0, 139, 1.0},
	"darkolivegreen":       {85, 107, 47, 1.0},
	"darkorange":           {255, 140, 0, 1.0},
	"darkorchid":           {153, 50, 204, 1.0},
	"darkred":              {139, 0, 0, 1.0},
	"darksalmon":           {233, 150, 122, 1.0},
	"darkseagreen":         {143, 188, 143, 1.0},
	"darkslateblue":        {72, 61, 139, 1.0},
	"darkslategray":        {47, 79, 79, 1.0},
	"darkslategrey":        {47, 79, 79, 1.0},
	"darkturquoise":        {0, 206, 209, 1.0},
	"darkviolet":           {148, 0, 211, 1.0},
	"deeppink":             {255, 20, 147, 1.0},
	"deepskyblue":          {0, 191, 255, 1.0},
	"dimgray":              {105, 105, 105, 1.0},
	"dimgrey":              {105, 105, 105, 1.0},
	"dodgerblue":           {30, 144, 255, 1.0},
	"firebrick":            {178, 34, 34, 1.0},
	"floralwhite":          {255, 250, 240, 1.0},
	"forestgreen":          {34, 139, 34, 1.0},
	"fuchsia":              {255, 0, 255, 1.0},
	"gainsboro":            {220, 220, 220, 1.0},
	"ghostwhite":           {248, 248, 255, 1.0},
	"gold":                 {255, 215, 0, 1.0},
	"goldenrod":            {218, 165, 32, 1.0},
	"gray":                 {128, 128, 128, 1.0},
	"green":                {0, 128, 0, 1.0},
	"greenyellow":          {173, 255, 47, 1.0},
	"grey":                 {128, 128, 128, 1.0},
	"honeydew":             {240, 255, 240, 1.0},
	"hotpink":              {255, 105, 180, 1.0},
	"indianred":            {205, 92, 92, 1.0},
	"indigo":               {75, 0, 130, 1.0},
	"ivory":                {255, 255, 240, 1.0},
	"khaki":                {240, 230, 140, 1.0},
	"lavender":             {230, 230, 250, 1.0},
	"lavenderblush":        {255, 240, 245, 1.0},
	"lawngreen":            {124, 252, 0, 1.0},
	"lemonchiffon":         {255, 250, 205, 1.0},
	"lightblue":            {173, 216, 230, 1.0},
	"lightcoral":           {240, 128, 128, 1.0},
	"lightcyan":            {224, 255, 255, 1.0},
	"lightgoldenrodyellow": {250, 250, 210, 1.0},
	"lightgray":            {211, 211, 211, 1.0},
	"lightgreen":           {144, 238, 144, 1.0},
	"lightgrey":            {211, 211, 211, 1.0},
	"lightpink":            {255, 182, 193, 1.0},
	"lightsalmon":          {255, 160, 122, 1.0},
	"lightseagreen":        {32, 178, 170, 1.0},
	"lightskyblue":         {135, 206, 250, 1.0},
	"lightslategray":       {119, 136, 153, 1.0},
	"lightslategrey":       {119, 136, 153, 1.0},
	"lightsteelblue":       {176, 196, 222, 1.0},
	"lightyellow":          {255, 255, 224, 1.0},
	"lime":                 {0, 255, 0, 1.0},
	"limegreen":            {50, 205, 50, 1.0},
	"linen":                {250, 240, 230, 1.0},
	"magenta":              {255, 0, 255, 1.0},
	"maroon":               {128, 0, 0, 1.0},
	"mediumaquamarine":     {102, 205, 170, 1.0},
	"mediumblue":           {0, 0, 205, 1.0},
	"mediumorchid":         {186, 85, 211, 1.0},
	"mediumpurple":         {147, 112, 219, 1.0},
	"mediumseagreen":       {60, 179, 113, 1.0},
	"mediumslateblue":      {123, 104, 238, 1.0},
	"mediumspringgreen":    {0, 250, 154, 1.0},
	"mediumturquoise":      {72, 209, 204, 1.0},
	"mediumvioletred":      {199, 21, 133, 1.0},
	"midnightblue":         {25, 25, 112, 1.0},
	"mintcream":            {245, 255, 250, 1.0},
	"mistyrose":            {255, 228, 225, 1.0},
	"moccasin":             {255, 228, 181, 1.0},
	"navajowhite":          {255, 222, 173, 1.0},
	"navy":                 {0, 0, 128, 1.0},
	"oldlace":              {253, 245, 230, 1.0},
	"olive":                {128, 128, 0, 1.0},
	"olivedrab":            {107, 142, 35, 1.0},
	"orange":               {255, 165, 0, 1.0},
	"orangered":            {255, 69, 0, 1.0},
	"orchid":               {218, 112, 214, 1.0},
	"palegoldenrod":        {238, 232, 170, 1.0},
	"palegreen":            {152, 251, 152, 1.0},
	"paleturquoise":        {175, 238, 238, 1.0},
	"palevioletred":        {219, 112, 147, 1.0},
	"papayawhip":           {255, 239, 213, 1.0},
	"peachpuff":            {255, 218, 185, 1.0},
	"peru":                 {205, 133, 63, 1.0},
	"pink":                 {255, 192, 203, 1.0},
	"plum":                 {221, 160, 221, 1.0},
	"powderblue":           {176, 224, 230, 1.0},
	"purple":               {128, 0, 128, 1.0},
	"red":                  {255, 0, 0, 1.0},
	"rosybrown":            {188, 143, 143, 1.0},
	"royalblue":            {65, 105, 225, 1.0},
	"saddlebrown":          {139, 69, 19, 1.0},
	"salmon":               {250, 128, 114, 1.0},
	"sandybrown":           {244, 164, 96, 1.0},
	"seagreen":             {46, 139, 87, 1.0},
	"seashell":             {255, 245, 238, 1.0},
	"sienna":               {160, 82, 45, 1.0},
	"silver":               {192, 192, 192, 1.0},
	"skyblue":              {135, 206, 235, 1.0},
	"slateblue":            {106, 90, 205, 1.0},
	"slategray":            {112, 128, 144, 1.0},
	"slategrey":            {112, 128, 144, 1.0},
	"snow":                 {255, 250, 250, 1.0},
	"springgreen":          {0, 255, 127, 1.0},
	"steelblue":            {70, 130, 180, 1.0},
	"tan":                  {210, 180, 140, 1.0},
	"teal":                 {0, 128, 128, 1.0},
	"thistle":              {216, 191, 216, 1.0},
	"tomato":               {255, 99, 71, 1.0},
	"turquoise":            {64, 224, 208, 1.0},
	"violet":               {238, 130, 238, 1.0},
	"wheat":                {245, 222, 179, 1.0},
	"white":                {255, 255, 255, 1.0},
	"whitesmoke":           {245, 245, 245, 1.0},
	"yellow":               {255, 255, 0, 1.0},
	"yellowgreen":          {154, 205, 50, 1.0},
}
