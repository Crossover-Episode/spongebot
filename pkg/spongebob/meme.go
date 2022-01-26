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
	font         *truetype.Font
}

const (
	maxLength = 540
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
		font:         font,
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
	n := m.strokeSizeForLength(len(text))
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(dc.Width())/2 + float64(dx)
			y := float64(dc.Height())/2 + float64(dy)
			dc.DrawStringWrapped(text, x, y, 0.5, 0.5, float64(dc.Width())-100.0, lineSpacing, gg.AlignCenter)
		}
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringWrapped(text, float64(dc.Width())/2, float64(dc.Height())/2, 0.5, 0.5, float64(dc.Width())-100.0, lineSpacing, gg.AlignCenter)

	return dc.Image(), nil
}

const (
	maxStroke = 4
	minStroke = 2
)

//returns a number [minStroke, maxStroke] mapped as length goes from maxLength
//to 0.
func (m *MemeGenerator) strokeSizeForLength(length int) int {
	return maxStroke - int(float64(length)/(maxLength+1)*(maxStroke-minStroke+1)+minStroke) + minStroke
}

//I don't like this, but I won't put in more time to make something more elegant
func (m *MemeGenerator) fontAndLineSpacingForLength(length int) (float64, float64) {
	if length < 50 {
		return 70, 2.5
	} else if length < 100 {
		return 40, 2.25
	} else if length < 200 {
		return 30, 2
	} else if length < 400 {
		return 25, 1.75
	}
	return 20, 1.5
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
