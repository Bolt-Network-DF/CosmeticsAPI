package utils

import (
	"github.com/anthonynsimon/bild/imgio"
	"github.com/df-mc/dragonfly/server/player/skin"
	"image"
	"image/color"
	"os"
)

var (
	bounds = map[int]struct {
		width  int
		height int
	}{
		64 * 32 * 4: {
			width:  64,
			height: 32,
		},
		64 * 64 * 4: {
			width:  64,
			height: 64,
		},
		128 * 128 * 4: {
			width:  128,
			height: 128,
		},
	}
)

// ReadCosmeticData ...
func ReadCosmeticData(path string) (image.Image, []byte) {
	texturePath := path + ".png"
	geometryPath := path + ".json"

	texture, err := imgio.Open(texturePath)
	if err != nil {
		panic(err)
	}

	geometry, err := os.ReadFile(geometryPath)
	if err != nil {
		panic(err)
	}

	return texture, geometry
}

// ReadCapeData ...
func ReadCapeData(path string) skin.Cape {
	i, err := imgio.Open(path + ".png")
	if err != nil {
		panic(err)
	}

	c := skin.NewCape(i.Bounds().Max.X, i.Bounds().Max.Y)

	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			col := i.At(x, y)
			r, g, b, a := col.RGBA()
			i := x*4 + i.Bounds().Max.X*y*4
			c.Pix[i], c.Pix[i+1], c.Pix[i+2], c.Pix[i+3] = uint8(r), uint8(g), uint8(b), uint8(a)
		}
	}
	return c
}

// GetImageFromSkin ...
func GetImageFromSkin(s skin.Skin) image.Image {
	skinData := string(s.Pix)
	bound, ok := bounds[len(skinData)]
	if !ok {
		panic("invalid skin size")
	}

	height := bound.height
	width := bound.width

	img := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}})

	pos := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := skinData[pos]
			pos++
			g := skinData[pos]
			pos++
			b := skinData[pos]
			pos++
			a := skinData[pos]
			pos++

			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return img
}
