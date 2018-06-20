package imgresize

import (
	"image"
)

// Fit makes a best fit of current into desired with constant
// aspect ratio.
func Fit(desired, current image.Rectangle) image.Rectangle {
	if current.Dx() == 0 || current.Dy() == 0 {
		return image.Rect(0, 0, 0, 0)
	}
	if desired.Dx() == 0 || desired.Dy() == 0 {
		return image.Rect(0, 0, 0, 0)
	}
	y := (desired.Dx() * current.Dy()) / current.Dx()
	if y == 0 {
		y = 1
	}
	xfit := image.Rect(0, 0, desired.Dx(), y)
	if xfit.In(desired) {
		return xfit
	}

	x := (desired.Dy() * current.Dx()) / current.Dy()
	if x == 0 {
		x = 1
	}
	yfit := image.Rect(0, 0, x, desired.Dy())
	if yfit.In(desired) {
		return yfit
	}

	panic("unreachable")
	return image.Rectangle{}
}
