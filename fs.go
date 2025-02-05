package snapshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// Create our snapshot root if it does not exist
// Returns true if created.
func ensureRootExists(root string) bool {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		_ = os.MkdirAll(root, 0755)
		return true
	}
	return false
}

// Load the snapshot by name, or create it if it does not exist.
// Returns the snapshot, if it was created, and an error if anything failed.
func loadOrCreateSnapshot(root, name string, actual image.Image) (image.Image, bool, error) {
	path := filepath.Join(root, name+".png")

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			if err != nil {
				return nil, false, fmt.Errorf("unable to create snapshot file %q: %w", path, err)
			}
			defer f.Close()
			if err = png.Encode(f, actual); err != nil {
				return nil, false, fmt.Errorf("unable to create snapshot file %q: %w", path, err)
			}
			return nil, true, nil
		} else {
			return nil, false, fmt.Errorf("unable to open snapshot file %q: %w", path, err)
		}
	}
	defer f.Close()

	snapshot, err := png.Decode(f)
	if err != nil {
		return nil, false, fmt.Errorf("unable to open snapshot file %q: %w", path, err)
	}

	return snapshot, false, nil
}

// Write a composite file (called "diff") which is left, diff, right.
// Returns the file name or an error on failure.
func writeComposite(tmpDir, name string, diff, left, right image.Image) (string, error) {
	return writeTempFile(
		tmpDir,
		name+"-diff-*.png",
		generateCompositeImage(diff, left, right),
	)
}

// Write the "actual" image to disk, i.e. not the snapshot we have, the new one
// Returns the file name or an error on failure.
func writeActual(tmpDir, name string, actual image.Image) (string, error) {
	return writeTempFile(tmpDir, name+"-actual-*.png", actual)
}

// Writes a file to the temp dir.
// Returns the file name or an error on failure.
func writeTempFile(tmpDir, template string, img image.Image) (string, error) {
	f, err := os.CreateTemp(tmpDir, template)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		return f.Name(), fmt.Errorf("unable to create file: %w", err)
	}

	return f.Name(), nil
}
