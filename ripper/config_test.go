package ripper

import (
	"testing"
	"go-cli/test"
)

func TestValidateConfig(t *testing.T) {
	t.Run("check nil config", func(t *testing.T) {
		test.AssertOn(t).ExpectError("expected error when validating nil config, but got none")(validateConfig(nil))
	})

	t.Run("check for empty workDir", func (t* testing.T) {
		c := &AppConf{WorkDirectory: "", MetaInfoRepo: "x/y/z", OutputDirectory: "/k/l/m"}
		test.AssertOn(t).ExpectError("expected error when validating empty workDir, but got none")(validateConfig(c))
	})

	t.Run("check for empty metaInfoRepo", func (t* testing.T) {
		c := &AppConf{WorkDirectory: "a/b/c", MetaInfoRepo: "", OutputDirectory: "/k/l/m"}
		test.AssertOn(t).ExpectError("expected error when validating empty metaInfoRepo, but got none")(validateConfig(c))
	})

	t.Run("check for empty metaInfoRepo", func (t* testing.T) {
		c := &AppConf{WorkDirectory: "a/b/c", MetaInfoRepo: "x/y/z", OutputDirectory: ""}
		test.AssertOn(t).ExpectError("expected error when validating empty metaInfoRepo, but got none")(validateConfig(c))
	})

	t.Run("validate no spaces in workdir", func(t *testing.T) {
		c := &AppConf{WorkDirectory: "a/b c/d", MetaInfoRepo: "x/y/z", OutputDirectory: "/k/l/m"}
		test.AssertOn(t).ExpectError("expected error when validating space in workDir, but got none")(validateConfig(c))
	})

	t.Run("validate no spaces in repodir", func(t *testing.T) {
		c := &AppConf{WorkDirectory: "a/b/c", MetaInfoRepo: "x/y y/z", OutputDirectory: "/k/l/m"}
		test.AssertOn(t).ExpectError("expected error when validating space in metaInfoRepo, but got none")(validateConfig(c))
	})

	t.Run("validate no spaces in defaultoutputdir", func(t *testing.T) {
		c := &AppConf{WorkDirectory: "a/b/c", MetaInfoRepo: "x/y/z", OutputDirectory: "/k/l l/m"}
		test.AssertOn(t).ExpectError("expected error when validating space in defaultOutputDirectory, but got none")(validateConfig(c))
	})

	t.Run("remove leading & trailing spaces in workDir, metaInfoRepo, defaultOutputDir", func(t *testing.T) {
		assert := test.AssertOn(t)
		c := &AppConf{WorkDirectory: " a/b/c  ", MetaInfoRepo: " x/y/z  ", OutputDirectory: "  /k/l/m "}
		assert.NotError(validateConfig(c))
		assert.StringsEqual("a/b/c", c.WorkDirectory)
		assert.StringsEqual("x/y/z", c.MetaInfoRepo)
		assert.StringsEqual("/k/l/m", c.OutputDirectory)
	})

	t.Run("success", func(t *testing.T) {
		assert := test.AssertOn(t)
		c := &AppConf{WorkDirectory: "a/b/c", MetaInfoRepo: "x/y/z", OutputDirectory: "/k/l/m"}
		assert.NotError(validateConfig(c))
		assert.StringsEqual("a/b/c", c.WorkDirectory)
		assert.StringsEqual("x/y/z", c.MetaInfoRepo)
		assert.StringsEqual("/k/l/m", c.OutputDirectory)
	})
}
