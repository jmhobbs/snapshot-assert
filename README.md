[![Go Reference](https://pkg.go.dev/badge/github.com/jmhobbs/snapshot-assert.svg)](https://pkg.go.dev/github.com/jmhobbs/snapshot-assert)
[![golangci-lint](https://github.com/jmhobbs/snapshot-assert/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/jmhobbs/snapshot-assert/actions/workflows/golangci-lint.yml)
[![Test and coverage](https://github.com/jmhobbs/snapshot-assert/actions/workflows/test.yml/badge.svg)](https://github.com/jmhobbs/snapshot-assert/actions/workflows/test.yml)

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
