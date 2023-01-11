package main

import (
	"first/entry"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:generate file2byteslice -input ./images/blast.png -output resources/blast.go -package resources -var BlastPng
//go:generate file2byteslice -input ./images/alien.png -output resources/alien.go -package resources -var AlienPng
//go:generate file2byteslice -input ./images/boss.png -output resources/boss.go -package resources -var BossPng
//go:generate file2byteslice -input ./images/player.png -output resources/player.go -package resources -var PlayerPng
//go:generate file2byteslice -input ./images/gopher.png -output resources/gopher.go -package resources -var GopherPng
//go:generate file2byteslice -input config.json -output resources/config.go -package resources -var ConfigJson
func main() {
	game := entry.NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
