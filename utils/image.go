package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func CreateCaptchaImage(captchaCode string) (*bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, 240, 80))
	buffer := &bytes.Buffer{}
	addLabel(img, 100, 45, captchaCode)

	if err := png.Encode(buffer, img); err != nil {
		return nil, err
	}

	return buffer, nil
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
