package scripts

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Game object that contains everything
type Game struct {
	snake       *Snake
	snakeChan   chan int
	window      *Window
	cherries    []*Cherry
	numCherries int
	numEnemies  int
	enemies     []*EnemySnake
	enemiesChan []chan int
	playing     bool
	points      int
	dotTime     int
}

// NewGame starts a new game
func NewGame(cherrys int, enemies int) Game {
	g := Game{
		playing:     true,
		points:      0,
		dotTime:     0,
		numCherries: cherrys,
		numEnemies:  enemies,
	}
	arrayCherrys := make([]*Cherry, g.numCherries)
	for i := 0; i < g.numCherries; i++ {
		arrayCherrys[i] = CreateCherry(&g)
		time.Sleep(20)
	}
	arrayEnemies := make([]*EnemySnake, g.numEnemies)
	for i := 0; i < len(arrayEnemies); i++ {
		arrayEnemies[i] = CreateEnemySnake(&g)
		time.Sleep(20)
	}
	enemiesChan := make([]chan int, g.numEnemies)
	for i := 0; i < len(enemiesChan); i++ {
		enemiesChan[i] = make(chan int)
		arrayEnemies[i].channelMovements = enemiesChan[i]
		go arrayEnemies[i].ChannelPipe()
		time.Sleep(20)
	}
	g.enemiesChan = enemiesChan
	g.cherries = arrayCherrys
	g.enemies = arrayEnemies
	g.snake = CreateSnake(&g)
	g.snakeChan = make(chan int)
	go g.snake.ChannelPipe()
	g.window = CreateWindow(&g, cherrys)
	return g
}

// End the game
func (g *Game) End() {
	g.playing = false
}

// Update general program
func (g *Game) Update() error {
	if g.playing {
		if g.numCherries == 0 {
			g.playing = false
		}
		g.dotTime = (g.dotTime + 1) % 20
		if err := g.snake.Direction(g.dotTime); err != nil {
			g.snakeChan <- g.dotTime
		}
		for i := 0; i < len(g.enemiesChan); i++ {
			g.enemiesChan[i] <- g.dotTime
		}
		xPos, yPos := g.snake.GetSerpentHead()
		for i := 0; i < len(g.cherries); i++ {
			if xPos == g.cherries[i].xPos && yPos == g.cherries[i].yPos {
				g.cherries[i].yPos = -20
				g.cherries[i].xPos = -20
				g.window.AddPoint()
				g.numCherries--
				g.snake.AddPoint()
				break
			}
		}
		for j := 0; j < len(g.enemies); j++ {
			xPos, yPos := g.enemies[j].GetSerpentHead()
			for i := 0; i < len(g.cherries); i++ {
				if xPos == g.cherries[i].xPos && yPos == g.cherries[i].yPos {
					g.cherries[i].yPos = -20
					g.cherries[i].xPos = -20
					g.numCherries--
					g.enemies[j].AddPoint()
					break
				}
			}
		}
	}
	for i := 0; i < g.numCherries; i++ {
		if err := g.cherries[i].Update(g.dotTime); err != nil {
			return err
		}
	}
	return nil
}

// Draw the game
func (g *Game) Draw(screen *ebiten.Image) error {
	if err := g.snake.Draw(screen, g.dotTime); err != nil {
		return err
	}
	for _, enemy := range g.enemies {
		if err := enemy.Draw(screen, g.dotTime); err != nil {
			return err
		}
	}
	if err := g.window.Draw(screen); err != nil {
		return err
	}
	for i := 0; i < len(g.cherries); i++ {
		if err := g.cherries[i].Draw(screen, g.dotTime); err != nil {
			return err
		}
	}
	if g.numCherries == 0 {
		g.window.EndAux(screen)
	}
	return nil
}
