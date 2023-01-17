package entry

import (
	"aircraft_war/resources"
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Boss struct {
	GameObject
	image       *ebiten.Image
	speedFactor float64
	Score       int
}

var (
	bossImg *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.BossPng))
	if err != nil {
		log.Fatal(err)
	}
	bossImg = ebiten.NewImageFromImage(img)
}
func NewBoss(cfg *Config) *Boss {
	width, height := bossImg.Size()
	g := GameObject{
		width:  width,
		height: height,
		x:      float64(cfg.ScreenWidth-width) / 2,
		y:      float64(-height),
	}
	ent := &Boss{
		image:       bossImg,
		speedFactor: 1,
		GameObject:  g,
	}
	return ent
}

func (boss *Boss) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(boss.x, boss.y)
	screen.DrawImage(boss.image, op)
}

func (boss *Boss) outOfScreen(cfg *Config) bool {
	return boss.y > float64(cfg.ScreenHeight)
}
