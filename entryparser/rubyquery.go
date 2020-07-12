package entryparser

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"strings"
)


type RubyParser struct {

}

// TODO(edkvm): Move to file
func NewHTTPRubyParserMiddelware() http.Handler {
	router := httprouter.New()
	parser := RubyParser{}
	router.GET("/process", parser.handleProcess())
	router.GET("/resize", parser.handleResize())
	router.GET("/crop", parser.handleCrop())




	return router
}

func (pr RubyParser) handleProcess() func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rawCmdChain, err := pr.ParseProcessQuery(r)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		log.Println("parser=ruby", "cmd=", rawCmdChain)
		output(rawCmdChain, w)
	}
}

func (pr RubyParser) handleResize() func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {



		rawCmdChain, err := pr.ParseResizeQuery(r)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		log.Println("parser=ruby", "cmd=", rawCmdChain)
		output(rawCmdChain, w)
	}
}

func (pr RubyParser) handleCrop() func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		rawCmdChain, err := pr.ParseCropQuery(r)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		log.Println("parser=ruby", "cmd=", rawCmdChain)
		output(rawCmdChain, w)
	}
}


func fnZero(s string) string {
	if s == "" {
		return "0"
	}

	return s
}

func (rp RubyParser) ParseCropQuery(r *http.Request) ([]string, error) {
	q := r.URL.Query()

	paramURL := q.Get("url")

	if paramURL == "" {
		return nil, fmt.Errorf("d")
	}

	rawCmd := make([]string, 0)

	// URL
	rawCmd = append(rawCmd, fmt.Sprintf("fetch(%s)", paramURL))

	// Crop always relative
	rawCmd = append(rawCmd,
		fmt.Sprintf("crop(%sx%s:%sx%s)",
			fnZero(q.Get("offset_x")),
			fnZero(q.Get("offset_y")),
			fnZero(q.Get("width")),
			fnZero(q.Get("height")),
		),
	)



	// Image Format
	accept := r.Header.Get("Accept")

	if accept != "" && strings.Contains(accept, "image/webp") {
		//rawCmd = append(rawCmd, "format(webp)")
	}

	return rawCmd, nil
}

func (rp RubyParser) ParseResizeQuery(r *http.Request) ([]string, error) {
	q := r.URL.Query()

	paramURL := q.Get("url")

	if paramURL == "" {
		return nil, fmt.Errorf("d")
	}

	rawCmd := make([]string, 0)

	// URL
	rawCmd = append(rawCmd, fmt.Sprintf("fetch(%s)", paramURL))

	// Crop always relative
	rawCmd = append(rawCmd,
		fmt.Sprintf("resize(%sx%s:%sx%s)",
			fnZero(q.Get("width")),
			fnZero(q.Get("height")),
		),
	)



	// Image Format
	accept := r.Header.Get("Accept")

	if accept != "" && strings.Contains(accept, "image/webp") {
		//rawCmd = append(rawCmd, "format(webp)")
	}

	return rawCmd, nil
}

func (rp RubyParser) ParseProcessQuery(r *http.Request) ([]string, error) {
	q := r.URL.Query()

	paramURL := q.Get("url")

	if paramURL == "" {
		return nil, fmt.Errorf("d")
	}
	
	rawCmd := make([]string, 0)

	// URL
	rawCmd = append(rawCmd, fmt.Sprintf("fetch(%s)", paramURL))

	uq, _ := url.QueryUnescape(r.URL.RawQuery)

	// Crop always relative
	if strings.Contains(uq, "filters[crop]") {
		rawCmd = append(rawCmd,
			fmt.Sprintf("crop(%sx%s:%sx%s)",
				fnZero(q.Get("filters[crop][o_x]")),
				fnZero(q.Get("filters[crop][o_y]")),
				fnZero(q.Get("filters[crop][w]")),
				fnZero(q.Get("filters[crop][h]")),
			),
		)

	}

	// Resize
	if strings.Contains(uq, "filters[resize]") {
		if strings.Contains(uq, "filters[resize][gravity]") {
			rawCmd = append(rawCmd,
				fmt.Sprintf("fill(%sx%s,%s)",
					fnZero(q.Get("filters[resize][w]")),
					fnZero(q.Get("filters[resize][h]")),
					q.Get("filters[resize][gravity]")),
			)
		} else {
			rawCmd = append(rawCmd,
				fmt.Sprintf("fill(%sx%s,Center)",
					fnZero(q.Get("filters[resize][w]")),
					fnZero(q.Get("filters[resize][h]"))),
			)
		}
	}
	// Quality

	// Image Format
	accept := r.Header.Get("Accept")

	if accept != "" && strings.Contains(accept, "image/webp") {
		//rawCmd = append(rawCmd, "format(webp)")
	}

	return rawCmd, nil
}