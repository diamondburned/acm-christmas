package vision

import (
	"image"
	"image/color"

	"github.com/pierrre/imageutil"
)

// BoundaryImage is an image that marks some boundary by filling the boundary
// with a color and leaving the rest of the image in a different color.
type BoundaryImage struct {
	img image.Image
	at  imageutil.AtFunc
	bc  [4]uint32
}

// NewBoundaryImage creates a new BoundaryImage from the given image and
// boundary color.
func NewBoundaryImage(img image.Image, bc color.Color) *BoundaryImage {
	cr, cg, cb, ca := bc.RGBA()
	return &BoundaryImage{
		img: img,
		at:  imageutil.NewAtFunc(img),
		bc:  [4]uint32{cr, cg, cb, ca},
	}
}

// PtIn returns true if the given point is in the boundary.
func (bi *BoundaryImage) PtIn(pt image.Point) bool {
	ar, ag, ab, aa := bi.at(pt.X, pt.Y)
	return ar == bi.bc[0] && ag == bi.bc[1] && ab == bi.bc[2] && aa == bi.bc[3]
}

func (bi *BoundaryImage) EachPt(f func(pt image.Point) (stop bool)) {
	y0 := bi.img.Bounds().Min.Y
	y1 := bi.img.Bounds().Max.Y
	x0 := bi.img.Bounds().Min.X
	x1 := bi.img.Bounds().Max.X

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if bi.PtIn(image.Point{X: x, Y: y}) {
				if f(image.Point{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

// Count returns the number of points in the boundary.
func (bi *BoundaryImage) Count() int {
	var count int
	bi.EachPt(func(image.Point) bool {
		count++
		return false
	})
	return count
}

func (bi *BoundaryImage) PtAt(i int) (image.Point, bool) {
	var pt image.Point
	bi.EachPt(func(p image.Point) bool {
		if i == 0 {
			pt = p
			return true
		}
		i--
		return false
	})
	return pt, i == 0
}