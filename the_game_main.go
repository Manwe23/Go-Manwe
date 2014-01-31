// The Game project main.go
package main

import (
	"TheGameMap"
	"fmt"
)

func main() {
	fmt.Println("Witaj swiat")
	var mapa TheGameMap.Mapa
	mapa.Render()
	mapa.Get(1, 2)
	mapa.Set(3, 4, "Wysokosc", "Góry")
	mapa.Load()

}
