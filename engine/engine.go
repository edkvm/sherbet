package engine

import (
	"image"
	"log"
	"time"

	"fmt"
)

type cmd struct {
	op    *Op
	param string
}

func Execute(cmdList []cmd) ([]byte, error) {

	var img image.Image

	// TODO(ekiselman): Convert this to an op as well
	// TODO(ekiselman): Add struct to hold []byte
	// inroder not to pass the memory around so much

	for _, cmd := range cmdList {

		result, err := opRunner(cmd.op, cmd.param, img)
		if err != nil {
			return nil, err
		}
		if result == nil {
			log.Println("(execute) status: ", cmd)
		}

		img = result
	}

	if img == nil {
		return nil, fmt.Errorf("Image is nil")
	}

	// TODO(ekiselman): Move this outside of the engine
	encode := JPEG

	result, err := Encode(encode, img)
	if err != nil {
		return nil, err
	}

	return result, nil
}



func opRunner(op *Op, param string, src image.Image) (image.Image, error) {
	defer func(then time.Time) {
		log.Println("(opRunner) op: ", op.Name, " duration: ", time.Since(then))
	}(time.Now())

	// TODO(ekiselman): Allocate buffer from a pool

	dst, err := op.Perform(src, param)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
