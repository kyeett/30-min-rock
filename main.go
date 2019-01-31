package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
	"github.com/peterhellberg/gfx"
)

const (
	screenWidth  = 400
	screenHeight = 400
	coinSize     = 32
)

type Player struct {
	gfx.Rect
	img *ebiten.Image
}

type Coin struct {
	gfx.Rect
	value int
}

type Game struct {
	player Player
	coins  []Coin
	score  int
}

func (g *Game) update(screen *ebiten.Image) error {
	// Movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Rect = g.player.Moved(gfx.V(-2, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Rect = g.player.Moved(gfx.V(2, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Rect = g.player.Moved(gfx.V(0, -2))
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Rect = g.player.Moved(gfx.V(0, 2))
	}

	// Check collision
	for i, coin := range g.coins {
		if coin.value == 0 {
			continue
		}
		if coin.Overlaps(g.player.Rect) {
			g.score += coin.value

			// "Remove" coin
			g.coins[i].value = 0

			g.coins = append(g.coins, newRandomCoin())
		}
	}

	// Draw coints
	for _, coin := range g.coins {
		if coin.value == 0 {
			continue
		}
		ebitenutil.DrawRect(screen, coin.Min.X, coin.Min.Y, coin.W(), coin.H(), colornames.Yellow)
	}

	// Draw player
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(g.player.Min.X, g.player.Min.Y)
	screen.DrawImage(g.player.img, options)

	// Draw score
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Money $%d", g.score), 5, 5)
	return nil
}

func newRandomCoin() Coin {
	x := rand.Float64() * screenWidth
	y := rand.Float64() * screenHeight
	return Coin{
		Rect:  gfx.R(x, y, x+coinSize, y+coinSize),
		value: rand.Intn(20) + 1,
	}
}

func main() {

	b, err := Asset("player.png")
	if err != nil {
		log.Fatal(b)
	}
	tmpImg, err := gfx.DecodePNG(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	img, err := ebiten.NewImageFromImage(tmpImg, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	g := Game{
		player: Player{
			Rect: gfx.R(20, 20, 20+64, 20+64),
			img:  img,
		},
		coins: []Coin{newRandomCoin(), newRandomCoin()},
	}

	ebiten.SetFullscreen(true)
	if err := ebiten.Run(g.update, screenWidth, screenHeight, 2, "game tutorial"); err != nil {
		log.Fatal(err)
	}
}
