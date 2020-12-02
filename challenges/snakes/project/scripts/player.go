package scripts

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var snakeDE *ebiten.DrawImageOptions

// Snake : Object which the player controls
type Snake struct {
	game             *Game
	numParts         int
	lastDir          string
	seperntHeadUp    ebiten.Image
	serpentHeadDown  ebiten.Image
	serpentHeadLeft  ebiten.Image
	serpentHeadRight ebiten.Image
	bodyH            ebiten.Image
	bodyV            ebiten.Image
	bodyParts        [][]float64
	pointsWaiting    int
	points           int
	channelMovements chan int
	collision        bool
}

// CreateSnake : Generates a snake
func CreateSnake(g *Game) *Snake {
	s := Snake{
		game:          g,
		numParts:      0,
		lastDir:       "right",
		pointsWaiting: 0,
		collision:     false,
	}
	s.channelMovements = make(chan int)
	s.bodyParts = append(s.bodyParts, []float64{300, 300})
	seperntHeadUp, _, _ := ebitenutil.NewImageFromFile("images/headSerpentUp.png", ebiten.FilterDefault)
	serpentHeadDown, _, _ := ebitenutil.NewImageFromFile("images/headSerpentDown.png", ebiten.FilterDefault)
	serpentHeadLeft, _, _ := ebitenutil.NewImageFromFile("images/headSerpentLeft.png", ebiten.FilterDefault)
	serpentHeadRight, _, _ := ebitenutil.NewImageFromFile("images/headSerpentRight.png", ebiten.FilterDefault)
	bodyH, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentH.png", ebiten.FilterDefault)
	bodyV, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentV.png", ebiten.FilterDefault)
	s.seperntHeadUp = *seperntHeadUp
	s.serpentHeadDown = *serpentHeadDown
	s.serpentHeadLeft = *serpentHeadLeft
	s.serpentHeadRight = *serpentHeadRight
	s.bodyH = *bodyH
	s.bodyV = *bodyV
	return &s
}

//ChannelPipe Pipe movements serpent
func (s *Snake) ChannelPipe() error {
	dotTime := <-s.channelMovements
	for {
		s.Direction(dotTime)
		dotTime = <-s.channelMovements
	}
}

// Direction Logical update of the snake
func (s *Snake) Direction(dotTime int) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) && s.lastDir != "right" { //movs the snake by pressing key
		s.lastDir = "right"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && s.lastDir != "down" {
		s.lastDir = "down"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && s.lastDir != "up" {
		s.lastDir = "up"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && s.lastDir != "left" {
		s.lastDir = "left"
		return nil
	}

	if dotTime == 1 { //snakes collide with the boudings
		xPos, yPos := s.GetSerpentHead()
		if xPos < 0 || xPos > 1050 || yPos < 0 || yPos > 690 || s.CollisionWithHimself() {
			s.collision = true
			s.game.End()
		}
	}
	return nil
}

// Draw the snake
func (s *Snake) Draw(screen *ebiten.Image, dotTime int) error {

	if s.game.playing {
		s.MoveSnake(dotTime)
	}
	snakeDE = &ebiten.DrawImageOptions{}

	xPos, yPos := s.GetSerpentHead()
	snakeDE.GeoM.Translate(xPos, yPos)

	if s.lastDir == "up" {
		screen.DrawImage(&s.seperntHeadUp, snakeDE)
	} else if s.lastDir == "down" {
		screen.DrawImage(&s.serpentHeadDown, snakeDE)
	} else if s.lastDir == "right" {
		screen.DrawImage(&s.serpentHeadLeft, snakeDE)
	} else if s.lastDir == "left" {
		screen.DrawImage(&s.serpentHeadRight, snakeDE)
	}

	for i := 0; i < s.numParts; i++ { //create the snakes parts
		partDO := &ebiten.DrawImageOptions{}
		xPos, yPos := s.GetSerpentBody(i)
		partDO.GeoM.Translate(xPos, yPos)
		if s.lastDir == "up" || s.lastDir == "down" {
			screen.DrawImage(&s.bodyH, partDO)
		} else {
			screen.DrawImage(&s.bodyV, partDO)
		}

	}

	return nil
}

// MoveSnake position values for the snake head
func (s *Snake) MoveSnake(dotTime int) {
	if dotTime == 1 {
		if s.pointsWaiting > 0 {
			s.numParts++
			s.pointsWaiting--
		}
		switch s.lastDir { //method for parts to follow the main snake, controls velocity
		case "up":
			s.TranslateHeadPos(0, -10)
		case "down":
			s.TranslateHeadPos(0, +10)
		case "right":
			s.TranslateHeadPos(10, 0)
		case "left":
			s.TranslateHeadPos(-10, 0)
		}

	}
}

// GetSerpentHead returns position of head
func (s *Snake) GetSerpentHead() (float64, float64) {
	return s.bodyParts[0][0], s.bodyParts[0][1]
}

// GetSerpentBody returns position of last body
func (s *Snake) GetSerpentBody(pos int) (float64, float64) {
	return s.bodyParts[pos+1][0], s.bodyParts[pos+1][1]
}

// ChangeBody adds a body to serpent
func (s *Snake) ChangeBody(newX, newY float64) {
	s.bodyParts = append([][]float64{{newX, newY}}, s.bodyParts...)
	s.bodyParts = s.bodyParts[:s.numParts+1]
}

// CollisionWithHimself evaluates the collision
func (s *Snake) CollisionWithHimself() bool {
	posX, posY := s.GetSerpentHead()
	for i := 1; i < len(s.bodyParts); i++ {
		if posX == s.bodyParts[i][0] && posY == s.bodyParts[i][1] {
			return true
		}
	}
	return false
}

// TranslateHeadPos changes body position in general
func (s *Snake) TranslateHeadPos(newXPos, newYPos float64) {
	newX := s.bodyParts[0][0] + newXPos
	newY := s.bodyParts[0][1] + newYPos
	s.ChangeBody(newX, newY)
}

// AddPoint controls game's points
func (s *Snake) AddPoint() {
	s.points++
	s.pointsWaiting++
}
