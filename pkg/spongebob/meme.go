package spongebob

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "embed"

	"github.com/fogleman/gg"
)

//go:embed imgs/meme.jpg
var spongebobJPG []byte

func GenerateMeme(text string) (image.Image, error) {
	img, err := jpeg.Decode(bytes.NewReader(spongebobJPG))
	if err != nil {
		return nil, err
	}
	
	meme, err := addTextToImage(img, text)
	if err != nil {
		return nil, err
	}

	return meme, nil
}

func addTextToImage(img image.Image, text string) (image.Image, error) {
	dc := gg.NewContextForImage(img)
	
	if err := dc.LoadFontFace("/Library/Fonts/Impact.ttf", 96); err != nil {
		return nil, err
	}
	dc.SetRGB(0, 0, 0)
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(dc.Width())/2 + float64(dx)
			y := float64(dc.Height())/2 + float64(dy)
			dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
		}
	}

	return dc.Image(), nil
}