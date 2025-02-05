package snapshot

import (
	"image"
	"testing"
)

var Default *Snapshots = New()

func Assert(t *testing.T, actual image.Image) {
	t.Helper()
	Default.Assert(t, actual)
}

func Fail(t *testing.T, actual image.Image) {
	t.Helper()
	Default.Fail(t, actual)
}

func Test(t *testing.T, actual image.Image) error {
	t.Helper()
	return Default.Test(t, actual)
}
