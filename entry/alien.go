package entry

import (
	"bytes"
	"first/resources"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Alien struct {
	GameObject
	image       *ebiten.Image
	speedFactor float64
}

var (
	alienImg *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.AlienPng))
	if err != nil {
		log.Fatal(err)
	}
	alienImg = ebiten.NewImageFromImage(img)
}

func NewAlien(cfg *Config) *Alien {
	width, height := alienImg.Size()
	x := float64(rand.Int63n(int64(cfg.ScreenWidth - width)))

	g := GameObject{
		width:  width,
		height: height,
		x:      x,
		y:      float64(-height),
	}
	ent := &Alien{
		image:       alienImg,
		speedFactor: cfg.AlienSpeedFactor,
		GameObject:  g,
	}
	return ent
}

func (alien *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(alien.x, alien.y)
	screen.DrawImage(alien.image, op)
}

func (alien *Alien) outOfScreen(cfg *Config) bool {
	return alien.y > float64(cfg.ScreenHeight)
}
