// Package golearning/font
package font

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
)

var (
	Face      *truetype.Font
	SmallFace *truetype.Font
)

func Init() {

	// Import font bytes
	fontBytes, err := ioutil.ReadFile("assets/arial-bold.ttf")
	if err != nil {
		panic("ioutil.ReadFile()")
	}

	Face, err = truetype.Parse(fontBytes)

	if err != nil {
		panic("truetype.Parse()")
	}

}
