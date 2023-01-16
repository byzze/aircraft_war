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
	width, height := bgImage.Size()
	g := GameObject{
		width:  width,
		height: height,
		x:      0,
		y:      0,
	}
	ent := &Viewport{
		image:      bgImage,
		GameObject: g,
	}
	return ent
}

func (viewport *Viewport) Draw(screen *ebiten.Image) {
	x16, y16 := viewport.Position()
	offsetX, offsetY := float64(-x16), float64(-y16)/16
	w, h := bgImage.Size()
	const repeat = 3
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			x := i * w
			y := h * j
			op.GeoM.Translate(float64(x), float64(y))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(bgImage, op)
		}
	}
}

func (viewport *Viewport) outOfScreen(cfg *Config) bool {
	return -viewport.y/16 > cfg.ScreenHeight
}

func (viewport *Viewport) Move(g *Game) {
	_, h := bgImage.Size()
	// maxX16 := w * 16
	maxY16 := h * 16

	// viewport.x += w / 32
	viewport.y += h / 32
	// viewport.x %= maxX16
	viewport.y %= maxY16
	if viewport.y >= 3000 {
		viewport.y = 0
	}
	logrus.Info(viewport.y)
}

func (p *Viewport) Position() (int, int) {
	return p.x, p.y
}
