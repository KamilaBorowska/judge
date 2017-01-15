package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const width = 800
const height = 600

func startGraphics(commands <-chan func(*board), boardSize int, delay time.Duration) {
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

	board := newBoard(boardSize, surface)

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
		select {
		case callback := <-commands:
			// Koniec listy
			if callback == nil {
				// Wyłącz kanał
				commands = nil
				break
			}
			callback(&board)
			window.UpdateSurface()
			time.Sleep(delay)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
