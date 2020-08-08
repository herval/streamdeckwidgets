package main

import (
	"bytes"
	"encoding/base64"
	"github.com/valyala/fastjson"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"meow.tf/streamdeck/sdk"
	"time"
)

var context string // set after plugin registration
var board [][]int
var w, h = 72, 72

func main() {
	// Initialize handlers for events
	sdk.RegisterAction("us.hervalicio.gameoflife", handleGameEvents)
	sdk.AddHandler(func(e *sdk.WillAppearEvent) {
		context = e.Context
	})

	// Open and connect the SDK
	err := sdk.Open()
	if err != nil {
		log.Fatalln(err)
	}

	go simulate()

	// Wait until the socket is closed, or SIGTERM/SIGINT is received
	sdk.Wait()
}

func render() image.Image {
	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{w, h},
		},
	)

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == 1 {
				img.Set(i, j, color.White)
			}
		}
	}

	return img
}

func seed() {
	board = newBoard()

	// seed initial population
	for i := rand.Intn(2000); i >= 0; i-- {
		board[rand.Intn(w)][rand.Intn(h)] = 1
	}
}

func newBoard() [][]int {
	board := make([][]int, w)
	for i := range board {
		board[i] = make([]int, h)
	}
	return board
}

func simulate() {
	seed()

	for {
		if context == "" {
			continue
		}

		img := render()

		// encode and send png to display
		var buff bytes.Buffer
		_ = png.Encode(&buff, img)

		str := base64.StdEncoding.EncodeToString(buff.Bytes())

		sdk.SetImage(
			context,
			"data:image/png;base64,"+str,
			0,
		)

		step()
		time.Sleep(time.Millisecond * 50)
	}
}

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
// Based on Java code from https://www.geeksforgeeks.org/program-for-conways-game-of-life/
func step() {
	newBoard := newBoard()

	for l := 1; l < w-1; l++ {
		for m := 1; m < h-1; m++ {
			// finding no Of Neighbours that are alive
			aliveNeighbours := 0

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					// mark everything as alive
					aliveNeighbours += board[l+i][m+j]
				}
			}

			// remove self
			aliveNeighbours -= board[l][m]

			if (board[l][m] == 1) && (aliveNeighbours < 2) { // Cell is lonely and dies
				newBoard[l][m] = 0
			} else if (board[l][m] == 1) && (aliveNeighbours > 3) { // Cell dies due to over population
				newBoard[l][m] = 0
			} else if (board[l][m] == 0) && (aliveNeighbours == 3) { // A new cell is born
				newBoard[l][m] = 1
			} else { // Remains the same
				newBoard[l][m] = board[l][m]
			}
		}
	}
	board = newBoard
}

func handleGameEvents(action, context string, payload *fastjson.Value, deviceId string) {
	seed()
}
