package watermark

import (
	"image"
	"image/jpeg"
)

var (
	defaultWatermarkTextColor       = image.Black
	defaultWatermarkBackgroundColor = image.White
	defaultJpegOptions              = &jpeg.Options{Quality: 95}
	defaultWatermarkOpacity         = 128
)

const ()

type Options struct {
	WatermarkHPercent             float64
	WatermarkWPercent             float64
	WatermarkHorizontalOffset     int
	WatermarkVerticalOffset       int
	WatermarkTextHorizontalOffset int
	WatermarkTextVerticalOffset   int
	WatermarkTextColor            *image.Uniform
	WatermarkBackgroundColor      *image.Uniform
	WatermarkOpacity              uint8

	JpegOptions     *jpeg.Options
	AutoOrientation bool

	RowSpacing  float64
	DPI         float64
	MaxFontSize float64
}

func (o *Options) watermarkHPercent() float64 {
	return o.WatermarkHPercent
}
func (o *Options) watermarkWPercent() float64 {
	return o.WatermarkWPercent
}
func (o *Options) watermarkHorizontalOffset() int {
	return o.WatermarkHorizontalOffset
}

func (o *Options) watermarkVerticalOffset() int {
	return o.WatermarkVerticalOffset
}
func (o *Options) watermarkTextHorizontalOffset() int {
	return o.WatermarkTextHorizontalOffset
}
func (o *Options) watermarkTextVerticalOffset() int {
	return o.WatermarkTextVerticalOffset
}
func (o *Options) watermarkTextColor() *image.Uniform {
	if o.WatermarkTextColor == nil {
		return defaultWatermarkTextColor
	}
	return o.WatermarkTextColor
}
func (o *Options) watermarkBackgroundColor() *image.Uniform {
	if o.WatermarkBackgroundColor == nil {
		return defaultWatermarkBackgroundColor
	}
	return o.WatermarkBackgroundColor
}

func (o *Options) watermarkOpacity() uint8 {
	if o.WatermarkOpacity == 0 {
		return uint8(defaultWatermarkOpacity)
	}
	return o.WatermarkOpacity
}
func (o *Options) jpegOptions() *jpeg.Options {
	if o.JpegOptions == nil {
		return defaultJpegOptions
	}
	return o.JpegOptions
}

func (o *Options) autoOrientation() bool {
	return o.AutoOrientation
}

func (o *Options) rowSpacing() float64 {
	return o.RowSpacing
}
func (o *Options) dpi() float64 {
	return o.DPI
}

func (o *Options) maxFontSize() float64 {
	if o.MaxFontSize == 0 {
		return 200
	}
	return o.MaxFontSize
}
