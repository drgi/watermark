package png_label

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	stepPercent = 0.05
)

type PngLabel struct {
	watermarkTextHorizontalOffset int
	watermarkTextVerticalOffset   int
	watermarkTextColor            *image.Uniform
	watermarkBackgroundColor      *image.Uniform
	rowSpacing                    float64
	dpi                           float64
	maxFontSize                   float64

	f *truetype.Font
}

func New(f *truetype.Font, textHorizontalOffset, textVerticalOffset int, textColor, bgColor *image.Uniform, spacing, dpi, maxFontSize float64) *PngLabel {
	return &PngLabel{
		f: f,

		watermarkTextHorizontalOffset: textHorizontalOffset,
		watermarkTextVerticalOffset:   textVerticalOffset,
		watermarkTextColor:            textColor,
		watermarkBackgroundColor:      bgColor,
		rowSpacing:                    spacing,
		dpi:                           dpi,
		maxFontSize:                   maxFontSize,
	}
}

func (p *PngLabel) LabelFromText(text []string, watermarkW, watermarkH int) ([]byte, error) {

	// base layer
	rgba := image.NewRGBA(image.Rect(0, 0, watermarkW, watermarkH))
	draw.Draw(rgba, rgba.Bounds(), p.watermarkBackgroundColor, image.ZP, draw.Src)
	var fontSizePt float64

	// fontSize count
	rows := rowCount(text)
	fontSizePtChras := p.maxFontSizePtByChar(text, p.f, watermarkW-(p.watermarkTextVerticalOffset*2), 4)
	fontSizePtRows := p.fontSizePxToPt(p.maxFontSizePxByRow(watermarkH-p.watermarkTextHorizontalOffset, p.rowSpacing, rows))
	if fontSizePtRows < fontSizePtChras {
		fontSizePt = fontSizePtRows
	} else {
		fontSizePt = fontSizePtChras
	}
	// Draw the text.
	h := font.HintingNone
	fontFace := truetype.NewFace(p.f, &truetype.Options{
		Size:    fontSizePt,
		DPI:     p.dpi,
		Hinting: h,
	})

	d := &font.Drawer{
		Dst:  rgba,
		Src:  p.watermarkTextColor,
		Face: fontFace,
	}
	fontSizePx := fontSizePt / 72 * p.dpi
	dy := int(math.Ceil(fontSizePx * p.rowSpacing))

	y := 0
	y += dy
	for _, s := range text {
		d.Dot = fixed.P(p.watermarkTextHorizontalOffset, y)
		d.DrawString(s)
		y += dy
	}

	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, rgba)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bytes := buffer.Bytes()
	return bytes, nil
}

// count target text rows
func rowCount(t []string) int {
	return len(t)
}

// count max chars in target text
func (p *PngLabel) maxFontSizePtByChar(t []string, f *truetype.Font, maxW int, fontSizePt float64) float64 {
	var max int
	step := fontSizePt * stepPercent
	h := font.HintingNone
	face := truetype.NewFace(f, &truetype.Options{
		Size:    fontSizePt,
		DPI:     p.dpi,
		Hinting: h,
	})
	for _, c := range t {
		runes := []rune(c)
		pxCount := 0
		for _, r := range runes {

			advance, ok := face.GlyphAdvance(r)
			if ok {
				pxCount += advance.Round()
			}
		}
		if max < pxCount {
			max = pxCount
		}
	}
	if max > maxW {
		return fontSizePt - step
	}
	if fontSizePt > p.maxFontSize {
		return p.maxFontSize
	}

	return p.maxFontSizePtByChar(t, f, maxW, fontSizePt+step)

}

// count max fontSize in px vertical
func (p *PngLabel) maxFontSizePxByRow(watermarkHight int, spacing float64, rowCount int) float64 {
	return (float64(watermarkHight) / float64(rowCount)) / spacing
}

// convert font size to pt
func (p *PngLabel) fontSizePxToPt(size float64) (pt float64) {
	pt = (size * 72) / p.dpi
	return
}
