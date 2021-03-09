package engine

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/jpeg"
	"log"
)

type Encoding uint8

const (
	JPEG Encoding = iota
	WebP
)

func Encode(enc Encoding, src image.Image) ([]byte, error) {

	switch enc {
	case JPEG:
		return jpegEncode(src)
	case WebP:
		return webpEncode(src)
	}

	return nil, fmt.Errorf("Format %s does not exists", string(enc))
}

func jpegEncode(src image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := jpeg.Encode(buf, src, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func webpEncode(src image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}

	log.Println("Encode webp")
	err := webp.Encode(buf, src, &webp.Options{Lossless: false, Quality: 80})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
