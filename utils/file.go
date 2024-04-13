package utils

import (
	"github.com/df-mc/dragonfly/server/player/skin"
	"image/png"
	"os"
)

// Read ...
func Read(path string) skin.Cape {
	f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	defer f.Close()

	i, err := png.Decode(f)

	if err != nil {
		panic(err)
	}

	c := skin.NewCape(i.Bounds().Max.X, i.Bounds().Max.Y)

	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			color := i.At(x, y)
			r, g, b, a := color.RGBA()
			i := x*4 + i.Bounds().Max.X*y*4
			c.Pix[i], c.Pix[i+1], c.Pix[i+2], c.Pix[i+3] = uint8(r), uint8(g), uint8(b), uint8(a)
		}
	}
	return c
}
