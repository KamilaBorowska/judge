package main

import (
	"errors"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/veandco/go-sdl2/sdl"
)

type board struct {
	size    int
	fields  []field
	surface *sdl.Surface
}

func newBoard(size int, surface *sdl.Surface) board {
	if surface != nil {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				rect := rectangleForPosition(size, x, y)
				surface.FillRect(&rect, 0xFF707070)
			}
		}
	}
	return board{
		size:    size,
		fields:  make([]field, size*size),
		surface: surface,
	}
}

func (b *board) at(x int, y int) (*field, error) {
	if x < 0 || y < 0 {
		return nil, errors.New("Ujemna wartość pola")
	}
	if x >= b.size || y >= b.size {
		return nil, errors.New("Pole poza tablicą")
	}
	return &b.fields[b.size*y : b.size*(y+1)][x], nil
}

func getColor(player uint8, x float64, y float64) uint32 {
	r, g, b := colorful.Hsl(hueForPlayer(player), 0.5+x*0.3, 0.5+y*0.2).RGB255()
	return 0xFF<<24 + uint32(r)<<16 + uint32(g)<<8 + uint32(b)
}

func hueForPlayer(player uint8) float64 {
	hues := [2]float64{
		// Red
		0,
		// Blue
		240,
	}
	return hues[player]
}

func (b *board) write(player uint8, x int, y int) error {
	f, err := b.at(x, y)
	if err != nil {
		return err
	}
	if f.exists {
		return errors.New("Wybrano wypełnione już pole")
	}
	*f = field{
		exists: true,
		player: player,
	}
	rect := rectangleForPosition(b.size, x, y)
	color := getColor(player, float64(x)/float64(b.size), float64(y)/float64(b.size))
	if b.surface != nil {
		b.surface.FillRect(&rect, color)
	}
	return nil
}

type field struct {
	exists bool
	player uint8
}

func rectangleForPosition(size int, x int, y int) sdl.Rect {
	paddingSize := int32(2)
	boxHorizontalSize := (width - paddingSize*(int32(size)-1)) / int32(size)
	boxVerticalSize := (height - paddingSize*(int32(size)-1)) / int32(size)
	return sdl.Rect{
		X: int32(x) * (boxHorizontalSize + paddingSize),
		Y: int32(y) * (boxVerticalSize + paddingSize),
		W: boxHorizontalSize,
		H: boxVerticalSize,
	}
}
