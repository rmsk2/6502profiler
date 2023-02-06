package memory

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type PicProcessor struct {
	image  *image.RGBA
	countX uint16
	countY uint16
}

func NewPicProcessor(w, h uint16) *PicProcessor {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(w), int(h)}

	return &PicProcessor{
		image:  image.NewRGBA(image.Rectangle{upLeft, lowRight}),
		countX: 0,
		countY: 0,
	}
}

func (p *PicProcessor) SetPoint(data uint8) {
	colGreen := color.NRGBA{0, 255, 0, 255}
	colWhite := color.NRGBA{255, 255, 255, 255}

	if data == 24 {
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

func (p *PicProcessor) Save(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() { f.Close() }()

	err = png.Encode(f, p.image)
	if err != nil {
		return err
	}

	return nil
}
