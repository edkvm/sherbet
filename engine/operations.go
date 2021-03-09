package engine

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

type Op struct {
	Name    string
	Params  string
	Perform func(image.Image, string) (image.Image, error)
}

var ops map[string]Op = make(map[string]Op, 4)

func init() {
	generateOps()
}

// TODO(edkvm): Convert to autogenerate
func generateOps() {

	ops["resize"] = Op{
		Name: "resize",
		Perform: func(src image.Image, param string) (image.Image, error) {
			// TODO(edkvm): Handle parsing error
			rect, _ := parseRelativeSize(param, image.Rect(0, 0, 1, 1))
			return resize(src, rect), nil
		},
	}

	ops["fill"] = Op{
		Name: "fill",
		Perform: func(src image.Image, param string) (image.Image, error) {
			// TODO(edkvm): Handle parsing error
			tuple := strings.Split(param, ",")

			rect, _ := parseRelativeSize(tuple[0], image.Rect(0, 0, 1, 1))
			anchor := parseAnchor(tuple[1])

			log.Println("Fill params: ", rect, ", ", anchor)
			return fill(src, rect, anchor), nil
		},
	}

	ops["crop"] = Op{
		Name: "crop",
		Perform: func(src image.Image, param string) (image.Image, error) {
			// TODO(edkvm): Handle parsing error

			rect, _ := parseRelativeBox(param, src.Bounds())

			log.Println("Crop params: ", rect)
			return crop(src, rect), nil
		},
	}

	ops["fetch"] = Op{
		Name: "fetch",
		Perform: func(src image.Image, param string) (image.Image, error) {
			log.Println("Fetch: ", param)
			return fetch(nil, param)
		},
	}

	ops["format"] = Op{
		Name: "format",
		Perform: func(src image.Image, param string) (image.Image, error) {
			log.Println("Format: ", param)
			return fetch(nil, param)
		},
	}

	ops["noop"] = Op{
		Name: "noop",
		Perform: func(src image.Image, param string) (image.Image, error) {
			return src, nil
		},
	}

}

func GetOp(name string) Op {
	if op, ok := ops[name]; ok {
		return op
	}
	return ops["noop"]
}

func fetch(_ image.Image, imagepath string) (image.Image, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", imagepath, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting the resource, code: %s", res.StatusCode)
	}

	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func resize(src image.Image, box image.Rectangle) image.Image {
	box = scaleRect(box, src.Bounds())
	return imaging.Resize(src, box.Dx(), box.Dy(), imaging.Lanczos)
}

func fill(src image.Image, box image.Rectangle, anchor imaging.Anchor) image.Image {
	box = scaleRect(box, src.Bounds())
	return imaging.Fill(src, box.Dx(), box.Dy(), anchor, imaging.Box)
}

func crop(src image.Image, box image.Rectangle) image.Image {
	return imaging.Crop(src, box)
}

func scaleRect(input image.Rectangle, orig image.Rectangle) image.Rectangle {

	var height = input.Dy()

	if input.Dy() == 0 && input.Dx() == 0 {
		return orig
	}

	if input.Dy() == 0 {
		height = (input.Dx() * orig.Dy()) / orig.Dx()
		return image.Rect(0, 0, input.Dx(), height)
	}

	if input.Dx() == 0 {
		height = (input.Dy() * orig.Dx()) / orig.Dy()
		return image.Rect(0, 0, height, input.Dy())
	}

	return input

}
