package snapshot

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateDiff(t *testing.T) {
	diffColor := color.RGBA{0, 255, 0, 255}

	t.Run("expected smaller than actual", func(t *testing.T) {
		expected := image.NewRGBA(image.Rect(0, 0, 2, 2))
		actual := image.NewRGBA(image.Rect(0, 0, 3, 3))

		expectedDiff := image.NewRGBA(image.Rect(0, 0, 3, 3))
		expectedDiff.Set(2, 0, diffColor)
		expectedDiff.Set(2, 1, diffColor)
		expectedDiff.Set(0, 2, diffColor)
		expectedDiff.Set(1, 2, diffColor)
		expectedDiff.Set(2, 2, diffColor)

		// EEA
		// EEA
		// AAA

		diff, diffPixels := generateDiff(expected, actual, diffColor)

		assert.Equal(t, 5, diffPixels)

		for y := 0; y < diff.Bounds().Dy(); y++ {
			for x := 0; x < diff.Bounds().Dx(); x++ {
				assert.Equal(t, expectedDiff.At(x, y), diff.At(x, y))
			}
		}
	})

	t.Run("expected larger than actual", func(t *testing.T) {
		expected := image.NewRGBA(image.Rect(0, 0, 4, 4))
		actual := image.NewRGBA(image.Rect(0, 0, 2, 2))

		expectedDiff := image.NewRGBA(image.Rect(0, 0, 4, 4))
		expectedDiff.Set(2, 0, diffColor)
		expectedDiff.Set(3, 0, diffColor)
		expectedDiff.Set(2, 1, diffColor)
		expectedDiff.Set(3, 1, diffColor)
		expectedDiff.Set(0, 2, diffColor)
		expectedDiff.Set(1, 2, diffColor)
		expectedDiff.Set(2, 2, diffColor)
		expectedDiff.Set(3, 2, diffColor)
		expectedDiff.Set(0, 3, diffColor)
		expectedDiff.Set(1, 3, diffColor)
		expectedDiff.Set(2, 3, diffColor)
		expectedDiff.Set(3, 3, diffColor)

		// AAEE
		// AAEE
		// EEEE
		// EEEE

		diff, diffPixels := generateDiff(expected, actual, diffColor)

		assert.Equal(t, 12, diffPixels)

		for y := 0; y < diff.Bounds().Dy(); y++ {
			for x := 0; x < diff.Bounds().Dx(); x++ {
				assert.Equal(t, expectedDiff.At(x, y), diff.At(x, y))
			}
		}
	})

	t.Run("no diff", func(t *testing.T) {
		expected := image.NewRGBA(image.Rect(0, 0, 2, 2))
		actual := image.NewRGBA(image.Rect(0, 0, 2, 2))

		expectedDiff := image.NewRGBA(image.Rect(0, 0, 2, 2))

		diff, diffPixels := generateDiff(expected, actual, diffColor)

		assert.Equal(t, 0, diffPixels)

		for y := 0; y < diff.Bounds().Dy(); y++ {
			for x := 0; x < diff.Bounds().Dx(); x++ {
				assert.Equal(t, expectedDiff.At(x, y), diff.At(x, y))
			}
		}
	})

	t.Run("partial diff", func(t *testing.T) {
		expected := image.NewRGBA(image.Rect(0, 0, 10, 10))
		actual := image.NewRGBA(image.Rect(0, 0, 10, 10))

		expectedDiff := image.NewRGBA(image.Rect(0, 0, 10, 10))
		expectedDiff.Set(3, 3, diffColor)
		expectedDiff.Set(4, 3, diffColor)
		expectedDiff.Set(3, 4, diffColor)
		expectedDiff.Set(4, 4, diffColor)

		// -----------
		// -----------
		// -----------
		// ---RR------
		// ---RR------
		// -----------
		// -----------
		// -----------
		// -----------
		// -----------

		// change 4px on actual to red
		draw.Draw(
			actual,
			image.Rect(3, 3, 5, 5),
			image.NewUniform(color.RGBA{255, 0, 0, 255}),
			image.Point{},
			draw.Over,
		)

		diff, diffPixels := generateDiff(expected, actual, diffColor)

		assert.Equal(t, 4, diffPixels)

		for y := 0; y < diff.Bounds().Dy(); y++ {
			for x := 0; x < diff.Bounds().Dx(); x++ {
				assert.Equal(t, expectedDiff.At(x, y), diff.At(x, y))
			}
		}
	})
}
