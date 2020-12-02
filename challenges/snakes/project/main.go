// The MIT License (MIT)
//
// Copyright (c) 2015-2016 Martin Lindhe
// Copyright (c) 2016      Hajime Hoshi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// BACKGROUND OBTAINED FROM -> GAME OF LIFE EBITEN

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"scripts/scripts"

	"github.com/hajimehoshi/ebiten"
)

var gm scripts.Game
var cherryN int
var enemiesN int

// init function default runs at start
func init() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 3 {
		fmt.Println("Wrong number of parameters")
		os.Exit(3)
	}

	var err error
	cherryN, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Cherry must be a number")
		os.Exit(3)
	}

	enemiesN, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Enemies must be number")
		os.Exit(3)
	}

	if cherryN == 0 || enemiesN == 0 {
		fmt.Println("At least one cherry and one enemy")
		os.Exit(3)
	}
	gm = scripts.NewGame(cherryN, enemiesN)
}

// World represents the game state.
type World struct {
	area   []bool
	width  int
	height int
}

// Game implements ebiten.Game interface.
type Game struct {
	world  *World
	pixels []byte
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]bool, width*height),
		width:  width,
		height: height,
	}
	w.Init(maxInitLiveCells)
	return w
}

// Init inits world with a random state for the background
func (w *World) Init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x] = true
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {

	g.world.Update()
	if err := gm.Update(); err != nil {
		return err
	}
	return nil
}

// Update game state by one tick. For the world
func (w *World) Update() {
	width := w.width
	height := w.height
	next := make([]bool, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop := NeighbourCount(w.area, width, height, x, y)
			switch {
			case pop < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[y*width+x] = false

			case (pop == 2 || pop == 3) && w.area[y*width+x]:
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[y*width+x] = true

			case pop > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[y*width+x] = false

			case pop == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[y*width+x] = true
			}
		}
	}
	w.area = next
}

// NeighbourCount calculates the Moore neighborhood of (x, y). Function already done
func NeighbourCount(a []bool, width, height, x, y int) int {
	c := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			if a[y2*width+x2] {
				c++
			}
		}
	}
	return c
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			pix[4*i] = 0xff
			pix[4*i+1] = 0xff
			pix[4*i+2] = 0xff
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

//measurements set as constance for better understanding
const (
	screenWidth  = 1080
	screenHeight = 720
)

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)

	if err := gm.Draw(screen); err != nil {
		fmt.Println(err)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1080, 720
}

//principal function
func main() {
	game := &Game{}

	game.world = NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/20))

	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Snake")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
