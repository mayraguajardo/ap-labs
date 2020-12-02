package scripts

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font/basicfont"
)

// Window GUI
type Window struct {
	game        *Game
	score       int
	points      int
	totalPoints int
	cherrys     int
}

// CreateWindow initialices window
func CreateWindow(g *Game, max int) *Window {
	w := Window{
		game:        g,
		points:      0,
		totalPoints: max,
		score:       0,
	}

	return &w
}

// AddPoint updates state GUI
func (w *Window) AddPoint() {
	w.points++
}

//TextFormat gives format to text
func TextFormat(text string) (w int, h int) {
	return 10 * len(text), 16
}

// EndGame show results and detect ends
func (w *Window) EndGame(screen *ebiten.Image) {

	if w.game.snake.collision == true {
		goText := "GAME OVER -> COLLISION DETECTED"
		textW, textH := TextFormat(goText)
		screenW := screen.Bounds().Dx()
		screenH := screen.Bounds().Dy()
		text.Draw(screen, goText, basicfont.Face7x13, screenW/2-textW/2, screenH/2+textH/2, color.White)
		text.Draw(screen, "Score: "+strconv.Itoa(w.points), basicfont.Face7x13, 30, 30, color.White)

	} else if w.points == w.score {
		goText := "GREAT YOU WIN!"
		textW, textH := TextFormat(goText)
		screenW := screen.Bounds().Dx()
		screenH := screen.Bounds().Dy()
		text.Draw(screen, goText, basicfont.Face7x13, screenW/2-textW/2, screenH/2+textH/2, color.White)
		text.Draw(screen, "Score: "+strconv.Itoa(w.points), basicfont.Face7x13, 30, 30, color.White)
	} else {
		goText := "GAME OVER -> ENEMY ATE MORE"
		textW, textH := TextFormat(goText)
		screenW := screen.Bounds().Dx()
		screenH := screen.Bounds().Dy()
		text.Draw(screen, goText, basicfont.Face7x13, screenW/2-textW/2, screenH/2+textH/2, color.White)
		text.Draw(screen, "Score: "+strconv.Itoa(w.points), basicfont.Face7x13, 30, 30, color.White)
	}
}

// Draw text points
func (w *Window) Draw(screen *ebiten.Image) error {
	if !w.game.playing {
		eatedCherrys := 0
		max := 0
		for i := 0; i < len(w.game.enemies); i++ {
			eatedCherrys += w.game.enemies[i].points
			if max < w.game.enemies[i].points {
				max = w.game.enemies[i].points
			}
		}

		eatedCherrys += w.game.snake.points
		if max < w.game.snake.points {
			max = w.game.snake.points
		}
		w.score = max
		w.cherrys = eatedCherrys
		w.EndGame(screen)
	}

	return nil
}

//EndAux calls end principal and prepares function
func (w *Window) EndAux(screen *ebiten.Image) {

	eatedCherrys := 0
	max := 0
	for i := 0; i < len(w.game.enemies); i++ {
		eatedCherrys += w.game.enemies[i].points
		if max < w.game.enemies[i].points {
			max = w.game.enemies[i].points
		}
	}

	eatedCherrys += w.game.snake.points
	if max < w.game.snake.points {
		max = w.game.snake.points
	}
	w.score = max
	w.cherrys = eatedCherrys
	w.EndGame(screen)
}
