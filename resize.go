package imgresize

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/anthonynsimon/bild/transform"
)

// Resize resizes JPEG image src to fit width and height best and writes
// the resulting JPEG image to dst.
func ResizeJPEG(dst io.Writer, src io.Reader, width, height int) error {
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	desired := image.Rect(0, 0, width, height)
	fit := Fit(desired, img.Bounds())
	resizedImg := transform.Resize(img, fit.Dx(), fit.Dy(), transform.Linear)
	err = jpeg.Encode(dst, resizedImg, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	return nil
}
