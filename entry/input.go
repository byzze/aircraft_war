package entry

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	msg            string
	lastBulletTime time.Time
}

func (i *Input) IsKeyPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (i *Input) Update(g *Game) {

	// 向左移动
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.x -= g.player.PlayerSpeedFactor
		if g.player.x < 1 {
			g.player.x = 1
		}
	}

	// 向右移动
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.x += g.player.PlayerSpeedFactor
		if g.player.x > float64(g.cfg.ScreenWidth)-float64(g.player.width) {
			g.player.x = float64(g.cfg.ScreenWidth) - float64(g.player.width)
		}
	}

	// 向上移动
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.y -= g.player.PlayerSpeedFactor
		if g.player.y < 1 {
			g.player.y = 1
		}
	}

	// 向下移动
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.y += g.player.PlayerSpeedFactor
		if g.player.y > float64(g.cfg.ScreenHeight)-float64(g.player.height) {
			g.player.y = float64(g.cfg.ScreenHeight) - float64(g.player.height)
		}
	}

	// 按键发射子弹
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(g.bullets) < g.cfg.MaxBulletNum && time.Now().Sub(i.lastBulletTime).Milliseconds() > g.cfg.BulletInterval {
			g.CreateBullet(g.player.BulletNum)
			i.lastBulletTime = time.Now()
		}
	}
}
