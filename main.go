package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/veandco/go-sdl2/sdl"
)

const width = 800
const height = 600

type board struct {
	size    int
	fields  []field
	surface *sdl.Surface
}

func newBoard(size int, surface *sdl.Surface) board {
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			rect := rectangleForPosition(size, x, y)
			surface.FillRect(&rect, 0xFF707070)
		}
	}
	return board{
		size:    size,
		fields:  make([]field, size*size),
		surface: surface,
	}
}

func (b *board) at(x int, y int) *field {
	return &b.fields[b.size*y : b.size*(y+1)][x]
}

func getColor(player uint8, x float64, y float64) uint32 {
	r, g, b := colorful.Hsl(hueForPlayer(player), 0.5+x*0.3, 0.5+y*0.2).RGB255()
	return 0xFF000000 + uint32(r)<<16 + uint32(g)<<8 + uint32(b)
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

func (b *board) write(player uint8, x int, y int) bool {
	f := b.at(x, y)
	if f.exists {
		return false
	}
	*f = field{
		exists: true,
		player: player,
	}
	rect := rectangleForPosition(b.size, x, y)
	color := getColor(player, float64(x)/float64(b.size), float64(y)/float64(b.size))
	b.surface.FillRect(&rect, color)
	return true
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

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Program sędziowski",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_HIDDEN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	background := sdl.Rect{W: width, H: height, X: 0, Y: 0}
	surface.FillRect(&background, 0xFF000000)

	board := newBoard(3, surface)
	board.write(0, 0, 0)
	board.write(0, 0, 1)

	window.Show()

	window.UpdateSurface()

END:
	for {
		for {
			event := sdl.PollEvent()
			if event == nil {
				break
			}
			switch event.(type) {
			case *sdl.QuitEvent:
				break END
			}
		}
	}
}
