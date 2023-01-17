package entry

import (
	"aircraft_war/resources"
	"bytes"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Prop struct {
	GameObject
	image       *ebiten.Image
	speedFactor float64
}

var (
	propImg *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.GopherPng))
	if err != nil {
		log.Fatal(err)
	}
	propImg = ebiten.NewImageFromImage(img)
}

func NewProp(cfg *Config) *Prop {
	width, height := propImg.Size()
	x := float64(rand.Int63n(int64(cfg.ScreenWidth - width)))
	g := GameObject{
		width:  width,
		height: height,
		x:      x,
		y:      0,
	}
	ent := &Prop{
		image:       propImg,
		GameObject:  g,
		speedFactor: cfg.PropSpeedFactor,
	}
	return ent
}

func (prop *Prop) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(prop.x, prop.y)
	screen.DrawImage(prop.image, op)
}

func (prop *Prop) outOfScreen(cfg *Config) bool {
	return prop.y > float64(cfg.ScreenHeight)
}
