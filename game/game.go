package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player           *Player
	lasers           []*Laser
	meteors          []*Meteor
	meteorSpawnTimer *Timer
	score            int
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer: NewTimer(20),
	}
	player := NewPlayer(g)
	g.player = player
	return g
}

// atualiza a lógica do jogo
func (g *Game) Update() error {
	g.player.Update()

	for _, l := range g.lasers {
		l.Update()
	}

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		m := NewMeteor()
		g.meteors = append(g.meteors, m)
	}

	for _, m := range g.meteors {
		m.Update()
	}

	playerHit := false
	for _, m := range g.meteors {
		if m.Collider().Intersects(g.player.Collider()) {
			playerHit = true
			break
		}
	}

	if playerHit {
		fmt.Println("Você perdeu")
		g.Reset()
	}

	for i := len(g.meteors) - 1; i >= 0; i-- {
		m := g.meteors[i]
		for j := len(g.lasers) - 1; j >= 0; j-- {
			l := g.lasers[j]
			if m.Collider().Intersects(l.Collider()) {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
				g.lasers = append(g.lasers[:j], g.lasers[j+1:]...)
				g.score += 1
				break
			}
		}
	}

	return nil
}

// desenha objetos na tela
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)
	for _, l := range g.lasers {
		l.Draw(screen)
	}
	for _, m := range g.meteors {
		m.Draw(screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddLasers(laser *Laser) {
	g.lasers = append(g.lasers, laser)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.lasers = nil
	g.meteors = nil
	g.meteorSpawnTimer.Reset()
	g.score = 0
}
