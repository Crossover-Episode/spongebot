package spongebob

import (
	"bytes"
	_ "embed"
	"image"
	"image/jpeg"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type MemeGenerator struct {
	spongebobImg image.Image
	font *truetype.Font
}

const (
	maxLength = 930
)

//go:embed imgs/meme.jpg
var spongebobJPG []byte

//go:embed fonts/impact.ttf
var impactFontBytes []byte

func NewGenerator() (*MemeGenerator, error) {
	img, err := jpeg.Decode(bytes.NewReader(spongebobJPG))
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(impactFontBytes)
	if err != nil {
		return nil, err
	}
	
	generator := &MemeGenerator{
		font: font,
		spongebobImg: img,

	}

	return generator, nil
}



func (m *MemeGenerator) GenerateMeme(text string) (*bytes.Buffer, error) {	
	spongebobText := ToText(text, false)
	meme, err := m.addTextToImage(m.spongebobImg, spongebobText)
	if err != nil {
		return nil, err
	}

	return prepMeme(meme)
}

func (m *MemeGenerator) addTextToImage(img image.Image, text string) (image.Image, error) {
	if len(text) > maxLength {
		text = text[:maxLength-3] + "..."
	}

	dc := gg.NewContextForImage(img)

	fontSize, lineSpacing := m.fontAndLineSpacingForLength(len(text))
	dc.SetFontFace(m.LoadFontFace(fontSize))

	dc.SetRGB(0, 0, 0)
	n := 4 // "stroke" size - increase this if you wanna watch the bot struggle
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(dc.Width())/2 + float64(dx)
			y := float64(dc.Height())/2 + float64(dy)
			dc.DrawStringWrapped(text, x, y, 0.5, 0.5, float64(dc.Width())-100.0, lineSpacing, gg.AlignCenter) //doing this at O(n^2) = horribly bad performance
		}
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringWrapped(text, float64(dc.Width())/2, float64(dc.Height())/2, 0.5, 0.5, float64(dc.Width())-100.0, lineSpacing, gg.AlignCenter)
	return dc.Image(), nil
}

//I don't like this, but I won't put in more time to make something more elegant
func (m *MemeGenerator) fontAndLineSpacingForLength(length int) (float64, float64) {
	if length < 100 {
		return 150, 3
	} else if length < 400 {
		return 125, 3
	} else if length < 700 {
		return 100, 2
	}
	return 75, 1
}

func (m *MemeGenerator) LoadFontFace(points float64) font.Face {
	face := truetype.NewFace(m.font, &truetype.Options{
		Size: points,
	})
	return face
}

func prepMeme(meme image.Image) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	jpegOpts := &jpeg.Options{
		Quality: 80,
	}
	err := jpeg.Encode(&buff, meme, jpegOpts)

	return &buff, err
}


