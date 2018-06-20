package imgresize

import (
	"fmt"
	"hash"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"testing"

	"github.com/matryer/is"
	"golang.org/x/crypto/blake2b"
)

func NewHash() hash.Hash {
	h, err := blake2b.New(16, nil)
	if err != nil {
		panic(err)
	}
	return h
}

func Sum(h interface{}) string {
	return fmt.Sprintf("%x", h.(hash.Hash).Sum(nil))
}

func RandRGBA(r io.Reader) color.RGBA {
	b := make([]byte, 4)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	return color.RGBA{b[0], b[1], b[2], b[3]}
}

func NewImage(bounds image.Rectangle, seed int64) *image.RGBA {
	r := rand.New(rand.NewSource(seed))
	img := image.NewRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			img.SetRGBA(x, y, RandRGBA(r))
		}
	}
	return img
}

func EncodeJPEG(img *image.RGBA) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		if err := jpeg.Encode(pw, img, &jpeg.Options{Quality: 100}); err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr
}

func TestResizeJPEG(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		seed         int64
		srcBounds    image.Rectangle
		dstBounds    image.Rectangle
		resultBounds image.Rectangle
		hash         string
	}{
		{
			seed:         1,
			srcBounds:    image.Rect(0, 0, 32, 32),
			dstBounds:    image.Rect(0, 0, 100, 100),
			resultBounds: image.Rect(0, 0, 100, 100),
			hash:         "72a7a00d94ad547701ad0f195a0fb8e9",
		},
		{
			seed:         2,
			srcBounds:    image.Rect(0, 0, 100, 100),
			dstBounds:    image.Rect(0, 0, 32, 32),
			resultBounds: image.Rect(0, 0, 32, 32),
			hash:         "ec06674a06acc0d501209ca0e275cd7b",
		},
		{
			seed:         3,
			srcBounds:    image.Rect(0, 0, 1, 1),
			dstBounds:    image.Rect(0, 0, 32, 10),
			resultBounds: image.Rect(0, 0, 10, 10),
			hash:         "b9c2d7899a04c5ea84a60aa3a9ccdc08",
		},
		{
			seed:         4,
			srcBounds:    image.Rect(0, 0, 100, 50),
			dstBounds:    image.Rect(0, 0, 25, 50),
			resultBounds: image.Rect(0, 0, 25, 12),
			hash:         "0a05f2f92aabf6c6dbc2db7a03263309",
		},
		{
			seed:         5,
			srcBounds:    image.Rect(0, 0, 15, 150),
			dstBounds:    image.Rect(0, 0, 400, 40),
			resultBounds: image.Rect(0, 0, 4, 40),
			hash:         "eb9ad92fd03504147356ffbcfaa9f5bc",
		},
		{
			seed:         6,
			srcBounds:    image.Rect(0, 0, 0, 0),
			dstBounds:    image.Rect(0, 0, 32, 32),
			resultBounds: image.Rect(0, 0, 0, 0),
			hash:         "2514fe0a0678efeeafb7467d03a51b2a",
		},
		{
			seed:         7,
			srcBounds:    image.Rect(0, 0, 32, 32),
			dstBounds:    image.Rect(0, 0, 0, 0),
			resultBounds: image.Rect(0, 0, 0, 0),
			hash:         "2514fe0a0678efeeafb7467d03a51b2a",
		},
		{
			seed:         8,
			srcBounds:    image.Rect(0, 0, 0, 100),
			dstBounds:    image.Rect(0, 0, 32, 23),
			resultBounds: image.Rect(0, 0, 0, 0),
			hash:         "2514fe0a0678efeeafb7467d03a51b2a",
		},
		{
			seed:         9,
			srcBounds:    image.Rect(0, 0, 100, 0),
			dstBounds:    image.Rect(0, 0, 32, 23),
			resultBounds: image.Rect(0, 0, 0, 0),
			hash:         "2514fe0a0678efeeafb7467d03a51b2a",
		},
	}
	for _, c := range cases {
		h := NewHash()
		jpegBounds := make(chan image.Rectangle)
		jpegR, jpegW := io.Pipe()
		go func() {
			image, _, err := image.Decode(jpegR)
			if err != nil {
				jpegR.CloseWithError(err)
			}
			jpegBounds <- image.Bounds()
		}()
		dst := io.MultiWriter(h, jpegW)
		src := EncodeJPEG(NewImage(c.srcBounds, c.seed))
		err := ResizeJPEG(dst, src, c.dstBounds.Dx(), c.dstBounds.Dy())
		log.Printf("%v hash: %s", c.seed, Sum(h))
		is.NoErr(err)
		is.Equal(Sum(h), c.hash)
		bounds := <-jpegBounds
		log.Printf("-> %v", bounds)
	}

}
