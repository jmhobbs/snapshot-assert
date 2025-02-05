package snapshot_test

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmhobbs/snapshot-assert"
	"github.com/stretchr/testify/assert"
)

func solidTestImage(width, height int, clr color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(clr), image.Point{}, draw.Over)
	return img
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return png.Decode(f)
}

var (
	red   color.Color = color.RGBA{255, 0, 0, 255}
	green color.Color = color.RGBA{0, 255, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
)

// Expects .snapshots/Test_Assert_NoChanges.png to exist
func Test_NoChanges(t *testing.T) {
	snapshot.Assert(t, solidTestImage(10, 10, blue))
}

// Expects .snapshots-custom/Test_Assert_CustomRoot.png NOT to exist
func Test_CustomRoot(t *testing.T) {
	defer os.RemoveAll(".snapshots-custom")

	i := snapshot.New(snapshot.WithStorageRoot(".snapshots-custom"))
	err := i.Test(t, solidTestImage(10, 10, blue))
	assert.NoError(t, err)

	_, err = os.Stat(".snapshots-custom/Test_CustomRoot.png")
	assert.NoError(t, err)

	err = i.Test(t, solidTestImage(10, 10, red))
	assert.Error(t, err)
}

// Expects .snapshots/Test_CustomDiffColor.png to exist but be BLUE
func Test_CustomDiffColor(t *testing.T) {
	// run snapshot, we expect an error
	i := snapshot.New(snapshot.WithDiffColor(red))
	err := i.Test(t, solidTestImage(10, 10, green))
	assert.Error(t, err)

	var errPixels snapshot.ErrPixelsDiffer
	if !errors.As(err, &errPixels) {
		t.Fatalf("Expected ErrPixelsDiffer error, got %T", err)
	}

	// Load the actual diff image generated
	actual, err := loadPNG(errPixels.CompositePath())
	if err != nil {
		t.Fatal(err)
	}

	// and what we expect
	expected, err := loadPNG("./fixtures/CustomDiffColor-diff.png")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

// Expects .snapshots/Test_CustomTempDir.png to exist, but not be blue
func Test_CustomTempDir(t *testing.T) {
	defer os.RemoveAll(".snapshots-temp")
	assert.NoError(t, os.MkdirAll(".snapshots-temp", 0755))

	i := snapshot.New(snapshot.WithTempDir(".snapshots-temp"))
	err := i.Test(t, solidTestImage(10, 10, blue))
	assert.Error(t, err)

	var errPixels snapshot.ErrPixelsDiffer
	if !errors.As(err, &errPixels) {
		t.Fatalf("Expected ErrPixelsDiffer error, got %T", err)
	}

	_, tempDir := filepath.Split(filepath.Dir(errPixels.CompositePath()))
	assert.Equal(t, ".snapshots-temp", tempDir)

	_, tempDir = filepath.Split(filepath.Dir(errPixels.ActualPath()))
	assert.Equal(t, ".snapshots-temp", tempDir)
}

// Expects .snapshots/Test_BoundsChange.png to exist, but be 10x10
func Test_BoundsChange(t *testing.T) {
	err := snapshot.Test(t, solidTestImage(20, 20, blue))
	assert.Error(t, err)

	var errBounds snapshot.ErrBoundsMismatch
	if !errors.As(err, &errBounds) {
		t.Fatalf("Expected ErrBoundsMismatch error, got %T", err)
	}
}

// Expects .snapshots/Test_Cleanup.png to exist, but not be red
func Test_Cleanup(t *testing.T) {
	defer os.RemoveAll(".snapshots-cleanup-temp")
	assert.NoError(t, os.MkdirAll(".snapshots-cleanup-temp", 0755))

	i := snapshot.New(snapshot.WithTempDir(".snapshots-cleanup-temp"))
	err := i.Test(t, solidTestImage(10, 10, red))
	assert.Error(t, err)

	files, err := os.ReadDir(".snapshots-cleanup-temp")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))

	i.Cleanup()

	files, err = os.ReadDir(".snapshots-cleanup-temp")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(files))
}
