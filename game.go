package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type game struct {
	players      [2]program
	commandQueue chan<- func(*board)
	size         int
}

func newGame(commandQueue chan<- func(*board), size int, programs [2]string) (game game, err error) {
	game.commandQueue = commandQueue
	game.size = size
	for i, value := range programs {
		var program program
		program, err = newProgram(value, uint8(i))
		if err != nil {
			return
		}
		game.players[i] = program
	}
	return
}

func (g *game) doGame() error {
	err := g.start()
	if err != nil {
		return err
	}
	err = g.retrievePings()
	if err != nil {
		return err
	}
	player := uint8(0)
	for {
		err := g.doGameStep(player)
		if err != nil {
			g.lose(player)
			return err
		}
		switch player {
		case 0:
			player = 1
		case 1:
			player = 0
		}
	}
}

func (g *game) start() error {
	for i := range g.players {
		err := g.players[i].start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *game) lose(player uint8) {
	process := g.players[player].executable.Process
	process.Kill()
	// Get rid of zombie process by reading its exit value
	process.Wait()

	var winner uint8
	switch player {
	case 0:
		winner = 1
	case 1:
		winner = 0
	}
	fmt.Printf("Program %d wygrał!\n", winner+1)
	exitStatus = int(winner) + 1
}

func (g *game) retrievePings() error {
	for _, player := range g.players {
		player.writeLine("PING")
	}
	type playerStatus struct {
		err    error
		player uint8
	}

	statuses := make(chan playerStatus, 2)
	for i := range g.players {
		i := uint8(i)
		player := &g.players[i]
		go func() {
			line, err := player.readLine()
			if err == nil {
				if strings.EqualFold(line, "PONG") {
					player.writeLine(strconv.Itoa(g.size))
				} else {
					err = errors.New("PONG nie był odpowiedzią")
				}
			}
			statuses <- playerStatus{err: err, player: i}
		}()
	}
	for range g.players {
		status := <-statuses
		if status.err != nil {
			g.lose(status.player)
			return status.err
		}
	}
	g.players[0].writeLine("ZACZYNAJ")
	return nil
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (g *game) doGameStep(playerNumber uint8) error {
	p := &g.players[playerNumber]
	line, err := p.readLine()
	if err != nil {
		return err
	}
	stringParts := strings.SplitN(line, " ", 4)
	if len(stringParts) < 4 {
		return fmt.Errorf("Oczekiwano czterech części odpowiedzi, otrzymano %d", len(stringParts))
	}
	var numericParts [4]int
	for i, part := range stringParts {
		numericParts[i], err = strconv.Atoi(part)
		if err != nil {
			return err
		}
	}

	if abs(numericParts[0]-numericParts[2])+abs(numericParts[1]-numericParts[3]) != 1 {
		return errors.New("Niepoprawny ruch")
	}

	result := make(chan error)
	g.commandQueue <- func(board *board) {
		var err error
		points := [2][2]int{
			{numericParts[0], numericParts[1]},
			{numericParts[2], numericParts[3]},
		}
		for _, point := range points {
			// Indeksowanie od jednego
			err = board.write(playerNumber, point[0]-1, point[1]-1)
		}
		result <- err
	}
	err = <-result
	if err == nil {
		var otherPlayer uint8
		switch playerNumber {
		case 0:
			otherPlayer = 1
		case 1:
			otherPlayer = 0
		}
		// Aby uniknąć ataków polegających na przekazywaniu formatów które ten
		// skrypt rozumie, ale przeciwnik nie.
		g.players[otherPlayer].writeLine(fmt.Sprintf("%d %d %d %d", numericParts[0], numericParts[1], numericParts[2], numericParts[3]))
	}
	return err
}

type programState int

type program struct {
	executable exec.Cmd
	stdout     bufio.Reader
	stdin      io.Writer
	state      programState
	playerID   uint8
}

func newProgram(name string, playerID uint8) (program program, err error) {
	program.playerID = playerID
	parts := strings.Split(name, " ")
	command := *exec.Command(parts[0], parts[1:]...)
	program.executable = command
	stdout, err := program.executable.StdoutPipe()
	program.stdout = *bufio.NewReader(stdout)
	if err != nil {
		return
	}
	program.stdin, err = program.executable.StdinPipe()
	return
}

func (p *program) start() error {
	return p.executable.Start()
}

func (p *program) rawReadLine() (string, error) {
	type readOutput struct {
		line string
		err  error
	}

	lineReceiver := make(chan readOutput, 1)
	go func() {
		line, err := p.stdout.ReadString('\n')
		line = strings.TrimSpace(line)
		//fmt.Printf("ODEBRANO #%d: %s\n", p.playerID+1, line)
		lineReceiver <- readOutput{
			line: line,
			err:  err,
		}
	}()

	select {
	case <-time.After(time.Second):
		return "", errors.New("Zbyt powolna odpowiedź")
	case line := <-lineReceiver:
		return line.line, line.err
	}
}

func (p *program) readLine() (string, error) {
	line, err := p.rawReadLine()
	return line, err
}

func (p *program) writeLine(line string) error {
	//fmt.Printf("WYSŁANO #%d: %s\n", p.playerID+1, line)
	_, err := p.stdin.Write([]byte(line + "\n"))
	return err
}
