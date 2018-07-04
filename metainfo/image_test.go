package metainfo

import (
	"testing"
	"fmt"
	"go-cli/test"
	"path/filepath"
	"bytes"
)

func TestImageFileName(t *testing.T) {
	repoPath := "a/b/c"
	id := "tt12345"
	ext := "png"
	expectedFileName := fmt.Sprintf("%s/%s/%s.%s", repoPath, SUBDIR_IMAGES, id, ext)

	fName := ImageFileName(repoPath, id, ext)
	test.AssertOn(t).StringsEqual(expectedFileName, fName)
}

func TestReadSaveImage(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	image := []byte{1,13,24,12,15,16,199}
	fName := filepath.Join(dir, "image.jpg")

	assert := test.AssertOn(t)
	assert.NotError(SaveImage(fName, image))
	gotImage, err := ReadImage(fName)
	assert.NotError(err)
	if !bytes.Equal(image, gotImage) {
		t.Errorf("saved and read images are not equal - saved %v, but read %v", image, gotImage)
	}
}
