package snapshot

import (
	"fmt"
	"image"
)

type DiffFiles struct {
	compositePath string
	actualPath    string
}

func (e DiffFiles) CompositePath() string {
	return e.compositePath
}

func (e DiffFiles) ActualPath() string {
	return e.actualPath
}

type ErrBoundsMismatch struct {
	DiffFiles
	expected image.Rectangle
	actual   image.Rectangle
}

func (e ErrBoundsMismatch) Error() string {
	return fmt.Sprintf("snapshot and image bounds differ: %v != %v", e.expected.Bounds(), e.actual.Bounds())
}

type ErrPixelsDiffer struct {
	DiffFiles
	count int
}

func (e ErrPixelsDiffer) Error() string {
	return fmt.Sprintf("snapshot and image differ by %d pixels", e.count)
}
