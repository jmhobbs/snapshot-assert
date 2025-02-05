package snapshot

import (
	"image"
	"image/color"
	"os"
	"testing"
)

type Snapshots struct {
	root      string
	tmp       string
	diffColor color.Color
	tempFiles []string
}

// Remove all of the temporary files created by this instance.
func (s *Snapshots) Cleanup() {
	for _, path := range s.tempFiles {
		_ = os.Remove(path)
	}
}

type Option func(*Snapshots)

func WithStorageRoot(root string) Option {
	return func(s *Snapshots) {
		s.root = root
	}
}

func WithTempDir(tmp string) Option {
	return func(s *Snapshots) {
		s.tmp = tmp
	}
}

func WithDiffColor(clr color.Color) Option {
	return func(s *Snapshots) {
		s.diffColor = clr
	}
}

func New(opts ...Option) *Snapshots {
	s := &Snapshots{
		root:      ".snapshots",
		tmp:       "",
		diffColor: color.RGBA{0, 255, 0, 255},
		tempFiles: []string{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

/*
Test compares the actual image to the snapshot image. If the snapshot image does not exist, it will be created.
If the images differ, an error is returned, either ErrBoundsMismatch or ErrPixelsDiffer.
The error can be inspected to get the paths to the actual and composite diff images.
A non-test failure will cause `testing.TB.Fatal()` to be called.
*/
func (s *Snapshots) Test(t testing.TB, actual image.Image) error {
	t.Helper()
	return s.TestWithName(t, t.Name(), actual)
}

// Run `Test` with a custom name.
func (s *Snapshots) TestWithName(t testing.TB, name string, actual image.Image) error {
	t.Helper()

	if ensureRootExists(s.root) {
		t.Logf("created snapshot directory %q", s.root)
	}

	snapshot, created, err := loadOrCreateSnapshot(s.root, name, actual)
	if err != nil {
		return err
	}
	if created {
		t.Log("created new snapshot")
		return nil
	}

	diff, diffPixels := generateDiff(snapshot, actual, s.diffColor)

	var diffFiles DiffFiles

	if diffPixels > 0 {
		compositePath, err := writeComposite(s.tmp, name, diff, snapshot, actual)
		if compositePath != "" {
			s.tempFiles = append(s.tempFiles, compositePath)
		}
		if err != nil {
			return err
		}

		actualPath, err := writeActual(s.tmp, name, actual)
		if actualPath != "" {
			s.tempFiles = append(s.tempFiles, actualPath)
		}
		if err != nil {
			return err
		}

		diffFiles = DiffFiles{
			compositePath: compositePath,
			actualPath:    actualPath,
		}
		t.Logf("new image: %q", actualPath)
		t.Logf("diff image: %q", compositePath)
	}

	if snapshot.Bounds() != actual.Bounds() {
		return ErrBoundsMismatch{
			DiffFiles: diffFiles,
			expected:  snapshot.Bounds(),
			actual:    actual.Bounds(),
		}
	}

	if diffPixels != 0 {
		return ErrPixelsDiffer{
			DiffFiles: diffFiles,
			count:     diffPixels,
		}
	}

	return nil
}

// Assert that the provided image matches the snapshot. Cause a test error if not.
func (s *Snapshots) Assert(t testing.TB, actual image.Image) {
	t.Helper()
	s.AssertWithName(t, t.Name(), actual)
}

// Assert that the provided image matches the named snapshot. Cause a test error if not.
func (s *Snapshots) AssertWithName(t testing.TB, name string, actual image.Image) {
	t.Helper()

	if err := s.TestWithName(t, name, actual); err != nil {
		t.Error(err)
	}
}

// Assert that the provided image matches the snapshot. Fail the test if not.
func (s *Snapshots) Fail(t testing.TB, actual image.Image) {
	t.Helper()
	s.FailWithName(t, t.Name(), actual)
}

// Assert that the provided image matches the named snapshot. Fail the test if not.
func (s *Snapshots) FailWithName(t testing.TB, name string, actual image.Image) {
	t.Helper()

	if err := s.Test(t, actual); err != nil {
		t.Fatal(err)
	}
}
