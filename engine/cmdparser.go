package engine

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"regexp"
	"strconv"
	"strings"
)

func BuildCmdChain(inputCmdList []string) ([]cmd, error) {

	cmds := make([]cmd, 0)
	for _, rawCmd := range inputCmdList {

		op, params, err := parseOp(rawCmd)
		if err != nil {
			panic(err)
		}

		cmds = append(cmds, cmd{
			op:    op,
			param: params,
		})

	}

	return cmds, nil
}

// Parses string from the form: 0.0x0.08285714285714285:0.9988392857142857x0.3632142857142857
// Returns a rectangle
func parseRelativeBox(raw string, mul image.Rectangle) (image.Rectangle, error) {
	corners := strings.Split(raw, ":")

	r1, _ := parseRelativeSize(corners[0], mul)
	r2, _ := parseRelativeSize(corners[1], mul)

	r := image.Rectangle{
		Min: image.Point{
			X: r1.Max.X,
			Y: r1.Max.Y,
		},
		Max: image.Point{
			X: r2.Max.X,
			Y: r2.Max.Y,
		},
	}

	return r, nil
}

// Parses string from the form: 0.0x0.08285714285714285
// Returns a rectangle at (0, 0)
func parseRelativeSize(raw string, orig image.Rectangle) (image.Rectangle, error) {
	rawValues := strings.Split(raw, "x")
	var values [2]float64

	for i, v := range rawValues {
		temp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		values[i] = float64(temp)
	}

	oW := 1.0
	oH := 1.0

	if values[0] <= 1 {
		oW = float64(orig.Dx())
	}

	if values[1] <= 1 {
		oH = float64(orig.Dy())
	}

	x := int(values[0] * oW)
	y := int(values[1] * oH)
	r := image.Rect(0, 0, x, y)

	return r, nil
}

func parseAnchor(raw string) (imaging.Anchor) {
	anchor := imaging.TopLeft

	switch raw {
	case "Center":
		anchor = imaging.Center
		break
	case "TopLeft":
		anchor = imaging.TopLeft
		break
	case "Top":
		anchor = imaging.Top
		break
	case "Bottom":
		anchor = imaging.Bottom
		break
	}


	return anchor
}
func parseOp(raw string) (*Op, string, error) {
	reNumber := `[|\d]*[|\.]{0,1}\d{0,32}`
	opMatchers := map[string]*regexp.Regexp{
		"fetch": regexp.MustCompile(`^fetch\((?P<url>.*)\)$`),
		"resize": regexp.MustCompile(fmt.Sprintf(`^resize\((?P<size>%sx%s)\)$`, reNumber, reNumber)),
		"fill": regexp.MustCompile(fmt.Sprintf(`^fill\((?P<size>%sx%s,\w*)\)$`, reNumber, reNumber)),
		"crop":   regexp.MustCompile(fmt.Sprintf(`^crop\((?P<crop>%sx%s:%sx%s)\)$`, reNumber, reNumber, reNumber, reNumber)),
	}

	params := ""
	var op *Op
	for opName, opMatch := range opMatchers {
		match := opMatch.FindAllStringSubmatch(raw, -1)

		if match != nil {
			tmp := GetOp(opName)
			op = &tmp
			params = match[0][1]

			break
		}

	}

	if op == nil {
		return nil, "", fmt.Errorf("could not parse command")
	}

	return op, params, nil
}