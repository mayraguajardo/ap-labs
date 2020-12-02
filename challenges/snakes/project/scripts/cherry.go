package scripts

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Cherry : Object which snakes eats
type Cherry struct {
	xLimit int
	yLimit int
	xPos   float64
	yPos   float64
	eaten  bool
	game   *Game
	cherry ebiten.Image
}

// CreateCherry Generates a Cherry at random position
func CreateCherry(g *Game) *Cherry {
	c := Cherry{
		game:   g,
		xLimit: 30,
		yLimit: 30,
		eaten:  false,
	}
	cherry, _, _ := ebitenutil.NewImageFromFile("images/cherry1.png", ebiten.FilterDefault)
	c.cherry = *cherry

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	c.xPos = float64(random.Intn(c.xLimit) * 20)
	c.yPos = float64(random.Intn(c.yLimit) * 20)
	return &c
}

// Update the state of a cherry to delete
func (c *Cherry) Update(dotTime int) error {
	if c.eaten == false {
		return nil
	}
	return nil
}

// Draw the cherry
func (c *Cherry) Draw(screen *ebiten.Image, dotTime int) error {
	snakeDE = &ebiten.DrawImageOptions{}
	snakeDE.GeoM.Translate(c.xPos, c.yPos)
	screen.DrawImage(&c.cherry, snakeDE)
	return nil
}
