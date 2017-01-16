package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

var exitStatus int

func main() {
	delay := flag.Duration("delay", 0, "opóźnienie wyświetlania kroków")
	noGraphics := flag.Bool("no-graphics", false, "wyłącz grafikę")
	flag.Parse()

	commandQueue := make(chan func(*board), 512)
	size, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Println("Pierwszy argument to rozmiar planszy")
		return
	}
	go func() {
		game, err := newGame(commandQueue, size, [2]string{flag.Arg(1), flag.Arg(2)})
		if err != nil {
			panic(err)
		}
		err = game.doGame()
		if err != nil {
			close(commandQueue)
			//fmt.Println("Przyczyna (porażki innego programu):", err)
		}
	}()

	if *noGraphics {
		if *delay != 0 {
			fmt.Println("Nie można używać opóźnienia bez grafiki, ignorowanie opóźnienia")
		}
		board := newBoard(size, nil)
		for command := range commandQueue {
			command(&board)
		}
	} else {
		// SDL2 spodziewa się że wszystkie jego wywołania będą z tego samego
		// wątku. Brak tej instrukcji powoduje dziwne braki wypełnień niektórych
		// pól.
		runtime.LockOSThread()
		startGraphics(commandQueue, size, *delay)
	}
	os.Exit(exitStatus)
}
