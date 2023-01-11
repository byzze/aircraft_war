package entry

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
)

type Config struct {
	ScreenWidth       int        `json:"screenWidth"`  // 屏幕宽度
	ScreenHeight      int        `json:"screenHeight"` // 屏幕高度
	Title             string     `json:"title"`        // 标题
	BgColor           color.RGBA `json:"bgColor"`
	PlayerSpeedFactor float64    `json:"playerSpeedFactor"`
	BulletWidth       int        `json:"bulletWidth"`
	BulletHeight      int        `json:"bulletHeight"`
	BulletColor       color.RGBA `json:"bulletColor"`
	BulletSpeedFactor float64    `json:"bulletSpeedFactor"`
	MaxBulletNum      int        `json:"maxBulletNum"`
	BulletInterval    int64      `json:"bulletInterval"`
	AlienSpeedFactor  float64    `json:"alienSpeedFactor"`
	TitleFontSize     int        `json:"titleFontSize"`
	FontSize          int        `json:"fontSize"`
	SmallFontSize     int        `json:"smallFontSize"`
	BossSpeedFactor   float64    `json:"bossSpeedFactor"`
	PropSpeedFactor   float64    `json:"propSpeedFactor"`
}

func loadConfig() *Config {
	f, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("os.Open failed: %v\n", err)
	}

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}
