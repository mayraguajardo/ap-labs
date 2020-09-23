package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	pic := make([][] uint8,dy)
	for x := 0; x < dy; x++{
		pic[x] = make([]uint8, dx)
		for y := 0; y < dx; y++{
			pic[x][y] = (uint8)(x^y)
		}
	}

     return pic
}

func main() {
	pic.Show(Pic)
}
