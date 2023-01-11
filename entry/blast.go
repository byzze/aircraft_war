package entry

import (
	"bytes"
	"first/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type Blast struct {
	GameObject
	image       *ebiten.Image
	speedFactor float64
}

var (
	blastImg *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.BlastPng))
	if err != nil {
		log.Fatal(err)
	}
	blastImg = ebiten.NewImageFromImage(img)
}

func NewBlast(cfg *Config, alien *Alien) *Blast {
	width, height := blastImg.Size()
	g := GameObject{
		width:  width,
		height: height,
		x:      alien.x,
		y:      alien.y,
	}
	ent := &Blast{
		image:      blastImg,
		GameObject: g,
	}
	return ent
}

func (blast *Blast) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(blast.x, blast.y)
	screen.DrawImage(blast.image, op)
}

func (blast *Blast) outOfScreen(cfg *Config) bool {
	return blast.y > float64(cfg.ScreenHeight)
}
