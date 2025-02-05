package snapshot

import (
	"image"
	"image/color"
	"image/draw"
)

func generateDiff(left, right image.Image, diffColor color.Color) (image.Image, int) {
	width := max(left.Bounds().Dx(), right.Bounds().Dx())
	height := max(left.Bounds().Dy(), right.Bounds().Dy())

	minWidth := min(left.Bounds().Dx(), right.Bounds().Dx())
	minHeight := min(left.Bounds().Dy(), right.Bounds().Dy())

	diff := image.NewRGBA(image.Rect(0, 0, width, height))
	diffPixels := 0
	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			if left.At(x, y) != right.At(x, y) {
				diff.Set(x, y, diffColor)
				diffPixels = diffPixels + 1
			}
		}
	}

	if minHeight < height {
		rows := height - minHeight
		diffPixels = diffPixels + rows*width
	}

	if minWidth < width {
		columns := width - minWidth
		// Calculate this using minHeight so we do not double count pixels from above
		diffPixels = diffPixels + columns*minHeight
	}

	draw.Draw(diff, image.Rect(minWidth, 0, width, height), image.NewUniform(diffColor), image.Point{}, draw.Over)
	draw.Draw(diff, image.Rect(0, minHeight, width, height), image.NewUniform(diffColor), image.Point{}, draw.Over)

	return diff, diffPixels
}

// Compose the expected, diff and actual into a single image, in that order.
func generateCompositeImage(diff, expected, actual image.Image) image.Image {
	width := expected.Bounds().Dx() + diff.Bounds().Dx() + actual.Bounds().Dx()
	height := diff.Bounds().Dy()

	composite := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(
		composite,
		expected.Bounds(),
		expected,
		image.Point{},
		draw.Over,
	)

	draw.Draw(
		composite,
		image.Rect(
			expected.Bounds().Dx(),
			0,
			expected.Bounds().Dx()+diff.Bounds().Dx(),
			diff.Bounds().Dy(),
		),
		diff,
		image.Point{},
		draw.Over,
	)

	draw.Draw(
		composite,
		image.Rect(
			expected.Bounds().Dx()+diff.Bounds().Dx(),
			0,
			width,
			actual.Bounds().Dy(),
		),
		actual,
		image.Point{},
		draw.Over,
	)

	return composite
}
