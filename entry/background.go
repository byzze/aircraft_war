package entry

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var (
	bgImage *ebiten.Image
)

type Viewport struct {
	x     int
	y     int
	image *ebiten.Image
	GameObject
}

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tile_png))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(img)
}

func NewViewport(cfg *Config) *Viewport {
	width, height := alienImg.Size()
	g := GameObject{
		width:  width,
		height: height,
		x:      0,
		y:      -160,
	}
	ent := &Viewport{
		image:      alienImg,
		GameObject: g,
	}
	return ent
}

func (viewport *Viewport) Draw(screen *ebiten.Image) {
	x16, y16 := viewport.Position()
	offsetX, offsetY := float64(-x16)/16, float64(-y16)/16
	// Draw bgImage on the screen repeatedly.
	const repeat = 3
	w, h := bgImage.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			op.GeoM.Translate(float64(offsetX), float64(offsetY))
			screen.DrawImage(bgImage, op)
		}
	}

}
func (viewport *Viewport) outOfScreen(cfg *Config) bool {
	logrus.Println(viewport.y / 16)
	return -viewport.y/16 > cfg.ScreenHeight
}
func (p *Viewport) Move() {
	p.y = p.y - 16
	// maxY16 := p.height * 16
	p.y += p.height / 32
	// p.y %= maxY16
}

func (p *Viewport) Position() (int, int) {
	return p.x, p.y
}
