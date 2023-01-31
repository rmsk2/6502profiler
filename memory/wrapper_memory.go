package memory

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type DataProcFunc func(data uint8)

type PictureWrappingMemory struct {
	mem              Memory
	specialAddresses map[uint16]DataProcFunc
	image            *image.RGBA
	countX           uint16
	countY           uint16
}

func NewPicWrapper(m Memory, w, h uint16) *PictureWrappingMemory {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(w), int(h)}

	res := &PictureWrappingMemory{
		mem:              m,
		image:            image.NewRGBA(image.Rectangle{upLeft, lowRight}),
		specialAddresses: make(map[uint16]DataProcFunc),
		countX:           0,
		countY:           0,
	}

	res.AddSpecialAddress(0xDDDD, res.setPoint)

	return res
}

func (p *PictureWrappingMemory) setPoint(data uint8) {
	col1 := color.NRGBA{0, 255, 0, 255}

	if data == 24 {
		p.image.Set(int(p.countX), int(p.countY), col1)
	} else {
		if (data & 1) != 0 {
			p.image.Set(int(p.countX), int(p.countY), col1)
		} else {
			p.image.Set(int(p.countX), int(p.countY), color.Black)
		}
	}

	p.countX++
	if p.countX >= uint16(p.image.Rect.Dx()) {
		p.countX = 0
		p.countY++
	}
}

func (p *PictureWrappingMemory) AddSpecialAddress(addr uint16, f DataProcFunc) {
	p.specialAddresses[addr] = f
}

func (p *PictureWrappingMemory) Save(fileName string) error {
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

func (p *PictureWrappingMemory) Load(address uint16) uint8 {
	return p.mem.Load(address)
}

func (p *PictureWrappingMemory) Store(address uint16, b uint8) {
	procFunc, ok := p.specialAddresses[address]
	if !ok {
		p.mem.Store(address, b)
		return
	}

	procFunc(b)
}

func (p *PictureWrappingMemory) GetStatistics() []BankAccessStatistics {
	return p.mem.GetStatistics()
}
