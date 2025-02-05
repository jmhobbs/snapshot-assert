# Snapshot Testing for Go

This is intended to be a simple drop in for adding image backed assertions to Go tests.

In short, it stores PNG files in a sub-directory and then compares the test inputs to the existing image.

## Usage

```go
package whatever_test

import (
    "image"
    "image/color"
    "image/draw"
    "testing"

    "github.com/jmhobbs/snapshot-assert"
)

func Test_Snapshot(t *testing.T) {
    // make a blue image, 50x50
    img := image.NewRGBA(image.Rect(0, 0, 50, 50))
    draw.Draw(
        img,
        img.Bounds(),
        image.NewUniform(color.RGBA{0, 0, 255, 255}),
        image.Point{},
        draw.Over,
    )

    // ensure it looks like the previous blue image
    snapshot.Assert(t, img)
}
```
