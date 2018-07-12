package htmlhelp

import (
	"fmt"
	"image/color"
	"math"
	"sync"
)

// ColorWheel generates visually distinct colors in a deterministic sequence
// Based on: https://stackoverflow.com/questions/309149/generate-distinctly-different-rgb-colors-in-graphs
type ColorWheel struct {
	currentIndex     uint8
	currentIntensity uint8
	m                sync.Mutex
}

// NextColor gets the next color in the pattern. If we've reached the last available color
// it will loop back to the beginning.
func (c *ColorWheel) NextColor() color.RGBA {
	c.m.Lock()
	defer c.m.Unlock()

	if c.currentIntensity == 0 {
		c.currentIntensity = 0xFF
	}

	oldIndex := c.currentIndex

	c.currentIndex++
	if c.currentIndex > 6 {
		c.currentIndex = 0
		c.currentIntensity /= 2
	}

	switch oldIndex {
	case 0:
		return color.RGBA{R: 0, G: 0, B: c.currentIntensity, A: 0xFF}
	case 1:
		return color.RGBA{R: 0, G: c.currentIntensity, B: 0, A: 0xFF}
	case 2:
		return color.RGBA{R: c.currentIntensity, G: 0, B: 0, A: 0xFF}
	case 3:
		return color.RGBA{R: 0, G: c.currentIntensity, B: c.currentIntensity, A: 0xFF}
	case 4:
		return color.RGBA{R: c.currentIntensity, G: 0, B: c.currentIntensity, A: 0xFF}
	case 5:
		return color.RGBA{R: c.currentIntensity, G: c.currentIntensity, B: 0, A: 0xFF}
	case 6:
		return color.RGBA{R: c.currentIntensity, G: c.currentIntensity, B: c.currentIntensity, A: 0xFF}
	default:
		panic("nope!")
	}
}

func HTMLColor(c color.RGBA) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %0.2f)", c.R, c.G, c.B, float64(c.A)/float64(math.MaxUint8))
}
