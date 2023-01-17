package entry

import (
	"aircraft_war/resources"
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	image *ebiten.Image
	GameObject
	Score             int
	BulletNum         int
	PlayerSpeedFactor float64
}

var (
	playerImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.PlayerPng))
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
}

func NewPlayer(cfg *Config) *Player {
	width, height := playerImage.Size()
	g := GameObject{
		width:  width,
		height: height,
		x:      float64(cfg.ScreenWidth-width) / 2,
		y:      float64(cfg.ScreenHeight - height),
	}

	player := &Player{
		image:             playerImage,
		GameObject:        g,
		BulletNum:         1,
		PlayerSpeedFactor: cfg.PlayerSpeedFactor,
	}

	return player
}
