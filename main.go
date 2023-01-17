package main

import (
	"aircraft_war/entry"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)




func main() {
	game := entry.NewGame()
 /// adsfxzcvxzxcvzxcvasd
  // Sdsadsadasd
  // 
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
