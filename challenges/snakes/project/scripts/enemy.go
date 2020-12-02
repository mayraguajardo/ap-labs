package scripts

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// EnemySnake object
type EnemySnake struct {
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
	seed             rand.Source
	pointsWaiting    int
	points           int
	channelMovements chan int
	collision        bool
}

// CreateEnemySnake initialize enemy
func CreateEnemySnake(g *Game) *EnemySnake {
	e := EnemySnake{
		game:          g,
		numParts:      0,
		lastDir:       "right",
		pointsWaiting: 0,
		collision:     false,
	}

	e.channelMovements = make(chan int)
	e.seed = rand.NewSource(time.Now().UnixNano())
	random := rand.New(e.seed)

	e.bodyParts = append(e.bodyParts, []float64{float64(random.Intn(30) * 20), float64(random.Intn(30) * 20)})

	seperntHeadUp, _, _ := ebitenutil.NewImageFromFile("images/headSerpentDownEnemy.png", ebiten.FilterDefault)
	serpentHeadDown, _, _ := ebitenutil.NewImageFromFile("images/headSerpentUpEnemy.png", ebiten.FilterDefault)
	serpentHeadLeft, _, _ := ebitenutil.NewImageFromFile("images/headSerpentLeftEnemy.png", ebiten.FilterDefault)
	serpentHeadRight, _, _ := ebitenutil.NewImageFromFile("images/headSerpentRightEnemy.png", ebiten.FilterDefault)
	bodyH, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentHEnemy.png", ebiten.FilterDefault)
	bodyV, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentVEnemy.png", ebiten.FilterDefault)
	e.seperntHeadUp = *seperntHeadUp
	e.serpentHeadDown = *serpentHeadDown
	e.serpentHeadLeft = *serpentHeadLeft
	e.serpentHeadRight = *serpentHeadRight
	e.bodyH = *bodyH
	e.bodyV = *bodyV

	return &e
}

//ChannelPipe Pipe movements enemy
func (s *EnemySnake) ChannelPipe() error {
	for {
		dotTime := <-s.channelMovements
		s.Direction(dotTime)
	}
}

// Direction Logical update of the enemy
func (s *EnemySnake) Direction(dotTime int) error {
	if dotTime == 1 {
		random := rand.New(s.seed)
		action := random.Intn(4)
		changingDirection := random.Intn(3)
		posX, posY := s.GetSerpentHead()
		if changingDirection == 0 {
			switch action {
			case 0:
				if posX < 1040 && s.lastDir != "left" {
					s.lastDir = "right"
				} else {
					s.lastDir = "left"
				}
				return nil
			case 1:
				if posY < 680 && s.lastDir != "up" {
					s.lastDir = "down"
				} else {
					s.lastDir = "up"
				}
				return nil
			case 2:
				if posY > 20 && s.lastDir != "down" {
					s.lastDir = "up"
				} else {
					s.lastDir = "down"
				}
				return nil
			case 3:
				if posX > 20 && s.lastDir != "right" {
					s.lastDir = "left"
				} else {
					s.lastDir = "right"
				}
				return nil
			}
		}
		//bounds collision
		if posX >= 1040 {
			s.lastDir = "left"
			return nil
		}
		if posX == 20 {
			s.lastDir = "right"
			return nil
		}
		if posY == 680 {
			s.lastDir = "up"
			return nil
		}
		if posY == 20 {
			s.lastDir = "down"
			return nil
		}
	}

	if dotTime == 1 { //checks collision with enemy snake
		xPos, yPos := s.game.snake.GetSerpentHead()
		if s.CollisionWithPlayer(xPos, yPos) {
			s.game.snake.collision = true
			s.game.End()
		}
	}
	return nil
}

// Draw the snake
func (s *EnemySnake) Draw(screen *ebiten.Image, dotTime int) error {
	if s.game.playing {
		s.UpdatePos(dotTime)
	}
	enemyDO := &ebiten.DrawImageOptions{}
	xPos, yPos := s.GetSerpentHead()
	enemyDO.GeoM.Translate(xPos, yPos)

	if s.lastDir == "up" {
		screen.DrawImage(&s.seperntHeadUp, enemyDO)
	} else if s.lastDir == "down" {
		screen.DrawImage(&s.serpentHeadDown, enemyDO)
	} else if s.lastDir == "right" {
		screen.DrawImage(&s.serpentHeadRight, enemyDO)
	} else if s.lastDir == "left" {
		screen.DrawImage(&s.serpentHeadLeft, enemyDO)
	}

	for i := 0; i < s.numParts; i++ {
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

// UpdatePos changes position values for the snake head
func (s *EnemySnake) UpdatePos(dotTime int) {
	if dotTime == 1 {
		if s.pointsWaiting > 0 {
			s.numParts++
			s.pointsWaiting--
		}
		switch s.lastDir {
		case "up":
			s.TranslateHeadPos(0, -20)
		case "down":
			s.TranslateHeadPos(0, +20)
		case "right":
			s.TranslateHeadPos(20, 0)
		case "left":
			s.TranslateHeadPos(-20, 0)
		}

	}
}

// AddPoint controls game's points
func (s *EnemySnake) AddPoint() {
	s.points++
	s.pointsWaiting++
}

// GetSerpentHead returns position of head
func (s *EnemySnake) GetSerpentHead() (float64, float64) {
	return s.bodyParts[0][0], s.bodyParts[0][1]
}

// GetSerpentBody returns position of last body
func (s *EnemySnake) GetSerpentBody(pos int) (float64, float64) {
	return s.bodyParts[pos+1][0], s.bodyParts[pos+1][1]
}

// CollisionWithPlayer evaluates the collision
func (s *EnemySnake) CollisionWithPlayer(xPos, yPos float64) bool {
	for i := 0; i < len(s.bodyParts); i++ {
		if xPos == s.bodyParts[i][0] && yPos == s.bodyParts[i][1] {
			return true
		}
	}
	return false
}

// TranslateHeadPos changes body position in general
func (s *EnemySnake) TranslateHeadPos(newXPos, newYPos float64) {
	newX := s.bodyParts[0][0] + newXPos
	newY := s.bodyParts[0][1] + newYPos
	s.AddParts(newX, newY)
}

// AddParts controls game's points
func (s *EnemySnake) AddParts(newX, newY float64) {
	s.bodyParts = append([][]float64{{newX, newY}}, s.bodyParts...)
	s.bodyParts = s.bodyParts[:s.numParts+1]
}
