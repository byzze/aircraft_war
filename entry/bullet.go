package entry

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	image *ebiten.Image
	GameObject
	speedFactor float64
}

func NewBullet(cfg *Config, player *Player) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImageWithOptions(rect, nil)
	img.Fill(cfg.BulletColor)
	g := GameObject{
		width:  cfg.BulletWidth,
		height: cfg.BulletHeight,
		x:      player.x + float64(player.width-cfg.BulletWidth)/2,
		y:      float64(player.Y()),
	}

	return &Bullet{
		image:       img,
		GameObject:  g,
		speedFactor: cfg.BulletSpeedFactor,
	}
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bullet.x, bullet.y)
	screen.DrawImage(bullet.image, op)
}

func (bullet *Bullet) outOfScreen() bool {
	return bullet.y < -float64(bullet.height)
}
