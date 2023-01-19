package watermark

import (
	"bytes"

	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"

	"github.com/disintegration/imaging"
	png_label "github.com/drgi/watermark/png"
	"github.com/golang/freetype/truetype"
)

type Watermark struct {
	watermarkHPercent             float64
	watermarkWPercent             float64
	watermarkHorizontalOffset     int
	watermarkVerticalOffset       int
	watermarkTextHorizontalOffset int
	watermarkTextVerticalOffset   int
	watermarkTextColor            *image.Uniform
	watermarkBackgroundColor      *image.Uniform
	watermarkOpacity              uint8

	jpegOptions     *jpeg.Options
	autoOrientation bool

	font        *truetype.Font
	rowSpacing  float64
	dpi         float64
	maxFontSize float64

	pngLabelGenerator *png_label.PngLabel
}

func New(fontPath string, opt *Options) (*Watermark, error) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	// parse font
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	w := &Watermark{
		watermarkHPercent:             opt.watermarkHPercent(),
		watermarkWPercent:             opt.watermarkWPercent(),
		watermarkHorizontalOffset:     opt.watermarkHorizontalOffset(),
		watermarkVerticalOffset:       opt.watermarkVerticalOffset(),
		watermarkTextHorizontalOffset: opt.watermarkTextHorizontalOffset(),
		watermarkTextVerticalOffset:   opt.watermarkTextVerticalOffset(),
		watermarkTextColor:            opt.watermarkTextColor(),
		watermarkBackgroundColor:      opt.watermarkBackgroundColor(),
		watermarkOpacity:              opt.watermarkOpacity(),

		jpegOptions:     opt.jpegOptions(),
		autoOrientation: opt.autoOrientation(),

		font:        f,
		rowSpacing:  opt.rowSpacing(),
		dpi:         opt.dpi(),
		maxFontSize: opt.maxFontSize(),
	}

	pngGen := png_label.New(
		f,
		w.watermarkTextHorizontalOffset,
		w.watermarkTextVerticalOffset,
		w.watermarkTextColor,
		w.watermarkBackgroundColor,
		w.rowSpacing,
		w.dpi,
		w.maxFontSize,
	)

	w.pngLabelGenerator = pngGen

	return w, nil
}

func (wm *Watermark) AddFromText(inputImage []byte, text []string) ([]byte, error) {

	// Open a test image.
	targetImage, err := imaging.Decode(bytes.NewReader(inputImage), imaging.AutoOrientation(wm.autoOrientation))
	if err != nil {
		return nil, err
	}
	// create png watermark from text
	targetImageW, targetImageH := wm.targetImageSize(targetImage)
	watermarkW, watermarkH := wm.watermarkSize(targetImageW, targetImageH)
	watermarkByte, err := wm.pngLabelGenerator.LabelFromText(text, watermarkW, watermarkH)
	if err != nil {
		return nil, err
	}

	watermark, err := imaging.Decode(bytes.NewReader(watermarkByte), imaging.AutoOrientation(wm.autoOrientation))
	if err != nil {
		return nil, err
	}

	mask := image.NewUniform(color.Alpha{wm.watermarkOpacity})
	offset := image.Pt(wm.watermarkHorizontalOffset, targetImageH-(watermarkH+wm.watermarkVerticalOffset))
	b := targetImage.Bounds()

	resultImage := image.NewRGBA(b)

	draw.Draw(resultImage, b, targetImage, image.ZP, draw.Src)
	draw.DrawMask(resultImage, resultImage.Bounds().Add(offset), watermark, image.ZP, mask, image.ZP, draw.Over)

	buffer := new(bytes.Buffer)

	err = jpeg.Encode(buffer, resultImage, wm.jpegOptions)
	if err != nil {
		return nil, err
	}
	bytes := buffer.Bytes()

	return bytes, nil
}

func (wm *Watermark) targetImageSize(targetImage image.Image) (int, int) {
	rect := targetImage.Bounds()
	w := rect.Size().X
	h := rect.Size().Y
	return w, h
}

func (wm *Watermark) watermarkSize(imageW, imageH int) (int, int) {
	wW := float64(imageW) * float64(wm.watermarkWPercent)
	wH := float64(imageH) * float64(wm.watermarkHPercent)
	return int(wW), int(wH)
}
