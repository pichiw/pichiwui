package components

import (
	"fmt"
	"image/color"
	"math"
)

func HTMLColor(c color.RGBA) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %0.2f)", c.R, c.G, c.B, float64(c.A)/float64(math.MaxUint8))
}
