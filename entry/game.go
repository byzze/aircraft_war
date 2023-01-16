package entry

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	input           *Input
	cfg             *Config
	player          *Player
	count           int
	bullets         map[*Bullet]struct{}
	blast           map[*Blast]struct{}
	aliens          sync.Map
	mode            Mode
	failCount       int // 被外星人碰撞和移出屏幕的外星人数量之和
	overMsg         string
	audioContext    *audio.Context
	audioBackground *audio.Player
	seBytes         []byte
	seCh            chan []byte
	aliensCh        chan *Alien
	viewport        []*Viewport
}

// 游戏重新开始初始化组件
func (g *Game) init() {
	g.failCount = 0
	g.overMsg = ""
	g.aliens = sync.Map{}
	g.bullets = make(map[*Bullet]struct{})
	g.aliensCh = make(chan *Alien, 1)
	g.blast = make(map[*Blast]struct{})
	g.viewport = make([]*Viewport, 0)
	g.seCh = make(chan []byte)
	g.audioContext = audio.NewContext(sampleRate)

	g.CreateFonts()
	g.CreateViewport()
	g.player = NewPlayer(g.cfg)
	go g.ProcessAliens()

}

const (
	sampleRate = 32000
)

func (g *Game) ProcessAliens() {
	for {
		if g.aliensCh != nil {
			select {
			case <-g.aliensCh:
				close(g.aliensCh)
				g.aliensCh = nil
				return
			default:
				// 变更敌人下降速度
				if g.cfg.AlienSpeedFactor < 15 {
					g.cfg.AlienSpeedFactor = g.cfg.AlienSpeedFactor + 0.000001
				}
				g.CreateAliens()
			}
		}
	}
}
func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight) // 设置窗口大小
	ebiten.SetWindowTitle(cfg.Title)                        // 设置窗口标题
	g := &Game{
		input: &Input{},
		cfg:   cfg,
	}
	g.init()
	go g.NewDefeatAudio()
	go g.ProcessBlast()

	return g
}

// 清除击杀特效
func (g *Game) ProcessBlast() {
	for {
		time.Sleep(time.Second * 2)
		if len(g.blast) != 0 {
			for k := range g.blast {
				delete(g.blast, k)
			}
		}
	}
}

func (g *Game) NewDefeatAudio() {
	// 子弹和敌人碰撞音效
	s, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
		return
	}
	b, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatal(err)
		return
	}
	g.seCh <- b
}

// 背景音乐
func (g *Game) NewBackgroundAudio() {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	var err error
	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(raudio.Ragtime_mp3))
	if err != nil {
		log.Fatal(err)
	}
	p, err := g.audioContext.NewPlayer(s)
	if err != nil {
		log.Fatal(err)
	}

	g.audioBackground = p
	g.audioBackground.Play()
}

// 画图，渲染图形化
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor)
	// 无限动画
	for _, v := range g.viewport {
		v.Draw(screen)
	}
	var titleTexts []string
	var texts []string
	// 分数
	scoreStr := fmt.Sprintf("%04d", g.player.Score)
	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"ALIEN INVASION"}
		texts = []string{"", "", "", "", "", "", "", "PRESS SPACE KEY", "", "OR LEFT MOUSE"}
	case ModeGame:

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.player.x, g.player.y)
		screen.DrawImage(g.player.image, op)
		// 子弹
		for bullet := range g.bullets {
			bullet.Draw(screen)
		}
		// 敌人
		g.aliens.Range(func(key, value interface{}) bool {
			if alien, ok := key.(*Alien); ok {
				alien.Draw(screen)
			}
			if boos, ok := key.(*Boss); ok {
				boos.Draw(screen)
			}
			if prop, ok := key.(*Prop); ok {
				prop.Draw(screen)
			}
			return true
		})
		// 击杀特效
		for blast := range g.blast {
			blast.Draw(screen)
		}

		text.Draw(screen, scoreStr, arcadeFont, g.cfg.ScreenWidth-len(scoreStr)*g.cfg.FontSize, g.cfg.FontSize, color.Black)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))

	case ModeOver:
		if g.overMsg == "You Win!" {
			texts = []string{"", "YOU WIN!", "", "YOU SCORE " + scoreStr}
		} else {
			texts = []string{"", "GAME OVER!", "", "YOU SCORE " + scoreStr}
		}
	}

	for i, l := range titleTexts {
		x := (g.cfg.ScreenWidth - len(l)*g.cfg.TitleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*g.cfg.TitleFontSize, color.White)
	}
	for i, l := range texts {
		x := (g.cfg.ScreenWidth - len(l)*g.cfg.FontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*g.cfg.FontSize, color.White)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

// 每帧刷新
func (g *Game) Update() error {
	for _, v := range g.viewport {
		v.Move(g)
		if v.outOfScreen(g.cfg) {
		}
	}
	select {
	case g.seBytes = <-g.seCh:
		close(g.seCh)
		g.seCh = nil
	default:
	}
	switch g.mode {
	case ModeTitle:
		if g.input.IsKeyPressed() {
			g.mode = ModeGame
		}
	case ModeGame:

		if g.player.Score == 20 || g.player.Score == 100 || g.player.Score == 150 {
			g.CreateBoss()
		}
		// 自动发射
		if time.Now().Sub(g.input.lastBulletTime).Milliseconds() > g.cfg.BulletInterval {
			g.CreateBullet(g.player.BulletNum)
			g.input.lastBulletTime = time.Now()
		}

		g.checkOutScreen()
		g.checkCollision()
		g.input.Update(g)

		if g.failCount >= 3 {
			g.overMsg = "Game Over!"
		} else if g.player.Score >= 500 {
			g.overMsg = "You Win!"
		}

		if len(g.overMsg) > 0 {
			g.mode = ModeOver
			g.audioBackground.Pause()
			g.aliensCh <- nil
		}

	case ModeOver:
		if g.input.IsKeyPressed() {
			g.init()
			g.mode = ModeTitle
			g.audioBackground.Play()
			if !BossLock.TryLock() {
				BossLock.Unlock()
			}
		}
	}

	return nil
}

// 检查物体是否超出布局边界
func (g *Game) checkOutScreen() {
	g.aliens.Range(func(key, value interface{}) bool {
		if alien, ok := key.(*Alien); ok {
			alien.y += alien.speedFactor
			if alien.outOfScreen(g.cfg) {
				g.aliens.Delete(key)
			}
			if CheckCollision(alien, g.player) {
				g.mode = ModeOver
				g.overMsg = "Game Over!"
				g.aliens.Delete(key)
				return false
			}
		}

		if boss, ok := key.(*Boss); ok {
			boss.y += boss.speedFactor
			if boss.outOfScreen(g.cfg) {
				g.aliens.Delete(key)
			}
			if CheckCollision(boss, g.player) {
				g.mode = ModeOver
				g.overMsg = "Game Over!"
				g.aliens.Delete(key)
				return false
			}
		}
		if prop, ok := key.(*Prop); ok {
			prop.y += prop.speedFactor
			if prop.outOfScreen(g.cfg) {
				g.aliens.Delete(key)
			}
			if CheckCollision(prop, g.player) {
				g.player.BulletNum = g.player.BulletNum + 2
				g.player.PlayerSpeedFactor = g.player.PlayerSpeedFactor + 2
				g.aliens.Delete(key)
			}
		}
		return true
	})

	for bullet := range g.bullets {
		bullet.y -= bullet.speedFactor
		if bullet.outOfScreen() {
			delete(g.bullets, bullet)
		}
	}
}

// 检查碰撞
func (g *Game) checkCollision() {
	for bullet := range g.bullets {
		g.aliens.Range(func(key, value interface{}) bool {
			if alien, ok := key.(*Alien); ok {
				r := CheckCollision(bullet, alien)
				if r {
					if g.seBytes != nil {
						sePlayer := g.audioContext.NewPlayerFromBytes(g.seBytes)
						sePlayer.Play()
						blast := NewBlast(g.cfg, alien)
						g.addBlast(blast)
					}

					g.player.Score++
					g.aliens.Delete(key)
					delete(g.bullets, bullet)
				}
			}

			if boss, ok := key.(*Boss); ok {
				r := CheckCollision(bullet, boss)
				if r {
					if g.seBytes != nil {
						sePlayer := g.audioContext.NewPlayerFromBytes(g.seBytes)
						sePlayer.Play()
					}
					boss.Score++
					delete(g.bullets, bullet)
				} else if boss.Score >= 40 {
					if !BossLock.TryLock() {
						BossLock.Unlock()
					}
					g.aliens.Delete(key)
					delete(g.bullets, bullet)
					log.Println("boss die:", boss.Score)
				}
			}
			return true
		})
	}
}

// 添加爆炸特效
func (g *Game) addBlast(blast *Blast) {
	g.blast[blast] = struct{}{}
}

// 添加子弹
func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

// 添加敌人
func (g *Game) addAlien(alien *Alien) {
	g.aliens.Store(alien, struct{}{})
}

// 添加boss
func (g *Game) addBoss(boss *Boss) {
	g.aliens.Store(boss, struct{}{})
}

// 添加道具
func (g *Game) addProp(prop *Prop) {
	g.aliens.Store(prop, struct{}{})
}

// 创建动画背景
func (g *Game) addViewport(view *Viewport) {
	g.viewport = append(g.viewport, view)
}

// 创建动画
func (g *Game) CreateViewport() {
	view := NewViewport(g.cfg)
	g.addViewport(view)
}

// 创建特效
func (g *Game) CreateProp() {
	prop := NewProp(g.cfg)
	g.addProp(prop)
}

// 创建子弹
func (g *Game) CreateBullet(num int) {
	alien := NewBullet(g.cfg, g.player)
	var xList = make([]int, num)
	tmpVal := 0
	x := int(alien.x)
	// 捡到道具后，子弹数量变化,生成子弹坐标
	for i := 0; i < num; i++ {
		if i%2 == 0 {
			xList[i] = x + tmpVal
			tmpVal = 10 + tmpVal
			continue
		}
		xList[i] = x - tmpVal
	}
	for _, v := range xList {
		alien := NewBullet(g.cfg, g.player)
		alien.x = float64(v)
		g.addBullet(alien)
	}
}

// 创建敌人
func (g *Game) CreateAliens() {
	rand.Seed(time.Now().Unix())
	alien := NewAlien(g.cfg)
	var flag bool
	// 创建时校验新的敌人是否和之前的敌人处于叠加，碰撞状态
	g.aliens.Range(func(key, value interface{}) bool {
		if a, ok := key.(*Alien); ok {
			if CheckCollision(alien, a) {
				flag = true
				return false
			}
		}
		if b, ok := key.(*Boss); ok {
			if CheckCollision(alien, b) {
				flag = true
				return false
			}
		}
		return true
	})

	if !flag {
		g.addAlien(alien)
	}
}

var BossLock sync.Mutex

// 创建boos
func (g *Game) CreateBoss() {
	if !BossLock.TryLock() {
		log.Println("trylock")
		return
	}
	boss := NewBoss(g.cfg)
	prop := NewProp(g.cfg)
	var flag bool
	g.aliens.Range(func(key, value interface{}) bool {
		if alien, ok := key.(*Alien); ok {
			if CheckCollision(alien, boss) {
				g.aliens.Delete(key)
				return false
			}
		}
		return true
	})
	if !flag {
		g.addBoss(boss)
		g.addProp(prop)
	}
}

func (g *Game) CreateFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.TitleFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.FontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.SmallFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

// 通用对象
type GameObject struct {
	width  int
	height int
	x      float64
	y      float64
}

func (gameObj *GameObject) Width() int {
	return gameObj.width
}

func (gameObj *GameObject) Height() int {
	return gameObj.height
}

func (gameObj *GameObject) X() float64 {
	return gameObj.x
}

func (gameObj *GameObject) Y() float64 {
	return gameObj.y
}

type Entity interface {
	Width() int
	Height() int
	X() float64
	Y() float64
}
