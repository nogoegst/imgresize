package imgresize

import (
	"image"
	"log"
	"testing"

	"github.com/matryer/is"
)

func TestFit(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		current  image.Rectangle
		desired  image.Rectangle
		expected image.Rectangle
	}{
		{
			current:  image.Rect(0, 0, 400, 799),
			desired:  image.Rect(0, 0, 300, 500),
			expected: image.Rect(0, 0, 250, 500),
		},
		{
			current:  image.Rect(0, 0, 30, 70),
			desired:  image.Rect(0, 0, 300, 500),
			expected: image.Rect(0, 0, 214, 500),
		},
		{
			current:  image.Rect(0, 0, 30, 70),
			desired:  image.Rect(0, 0, 0, 0),
			expected: image.Rect(0, 0, 0, 0),
		},
		{
			current:  image.Rect(0, 0, 0, 0),
			desired:  image.Rect(0, 0, 12, 343),
			expected: image.Rect(0, 0, 0, 0),
		},
	}
	for _, c := range cases {
		fit := Fit(c.desired, c.current)
		log.Printf("%v", fit)
		is.Equal(fit, c.expected)
	}
}
