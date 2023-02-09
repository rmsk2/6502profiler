package memory

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PicProcessor struct {
	image    *image.RGBA
	countX   uint16
	countY   uint16
	fileName string
	maxIter  uint16
}

// NewPicProcFomConfig parses a config string of the form "pic:apfel2.png:320x200@24" and
// creates the corresponding PicProcessor struct.
func NewPicProcFromConfig(conf string) (MemWrapper, bool) {
	components := strings.Split(conf, ":")

	if len(components) != 3 {
		return nil, false
	}

	if components[0] != "pic" {
		return nil, false
	}

	r := regexp.MustCompile(`^([0-9]+)x([0-9]+)@([0-9]+)$`)

	matches := r.FindStringSubmatch(components[2])
	if matches == nil {
		return nil, false
	}

	xRes, _ := strconv.ParseUint(matches[1], 10, 16)
	yRes, _ := strconv.ParseUint(matches[2], 10, 16)
	maxIter, _ := strconv.ParseUint(matches[3], 10, 16)

	return NewPicProcessor(components[1], uint16(xRes), uint16(yRes), uint16(maxIter)), true
}

func NewPicProcessor(f string, w, h, i uint16) *PicProcessor {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(w), int(h)}

	return &PicProcessor{
		image:    image.NewRGBA(image.Rectangle{upLeft, lowRight}),
		countX:   0,
		countY:   0,
		fileName: f,
		maxIter:  i,
	}
}

func (p *PicProcessor) Write(data uint8) {
	colGreen := color.NRGBA{0, 255, 0, 255}
	colWhite := color.NRGBA{255, 255, 255, 255}

	if data == uint8(p.maxIter) {
		p.image.Set(int(p.countX), int(p.countY), colGreen)
	} else {
		if (data & 1) != 0 {
			p.image.Set(int(p.countX), int(p.countY), colGreen)
		} else {
			p.image.Set(int(p.countX), int(p.countY), colWhite)
		}
	}

	p.countX++
	if p.countX >= uint16(p.image.Rect.Dx()) {
		p.countX = 0
		p.countY++
	}
}

func (p *PicProcessor) Close() {
	f, err := os.Create(p.fileName)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file for pic writer: %v", err))
	}
	defer func() { f.Close() }()

	err = png.Encode(f, p.image)
	if err != nil {
		panic(fmt.Sprintf("Unable to encode picture: %v", err))
	}
}
