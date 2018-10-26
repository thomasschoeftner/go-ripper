package main

import (
	"testing"
	"go-cli/test"
	"go-ripper/ripper"
)

func TestValidateConfig(t *testing.T) {
	t.Run("check nil config", func(t *testing.T) {
		test.AssertOn(t).ExpectError("expected error when validating nil config, but got none")(validateConfig(nil))
	})

	t.Run("check for empty workDir", func (t* testing.T) {
		c := &ripper.AppConf{
			WorkDirectory: "",
			MetaInfoRepo: "x/y/z"}
		test.AssertOn(t).ExpectError("expected error when validating empty workDir, but got none")(validateConfig(c))
	})

	t.Run("check for empty metaInfoRepo", func (t* testing.T) {
		c := &ripper.AppConf{
			WorkDirectory: "a/b/c",
			MetaInfoRepo: ""}
		test.AssertOn(t).ExpectError("expected error when validating empty metaInfoRepo, but got none")(validateConfig(c))
	})

	t.Run("validate no spaces in workdir", func(t *testing.T) {
		c := &ripper.AppConf{
			WorkDirectory: "a/b c/d",
			MetaInfoRepo: "x/y/z"}
		test.AssertOn(t).ExpectError("expected error when validating space in workDir, but got none")(validateConfig(c))
	})

	t.Run("validate no spaces in repodir", func(t *testing.T) {
		c := &ripper.AppConf{
			WorkDirectory: "a/b/c",
			MetaInfoRepo: "x/y y/z"}
		test.AssertOn(t).ExpectError("expected error when validating space in metaInfoRepo, but got none")(validateConfig(c))
	})

	t.Run("remove leading & trailing spaces in workDir & metaInfoRepo", func(t *testing.T) {
		assert := test.AssertOn(t)
		c := &ripper.AppConf{
			WorkDirectory: " a/b/c  ",
			MetaInfoRepo: " x/y/z  "}
		assert.NotError(validateConfig(c))
		assert.StringsEqual("a/b/c", c.WorkDirectory)
		assert.StringsEqual("x/y/z", c.MetaInfoRepo)
	})

	t.Run("success", func(t *testing.T) {
		assert := test.AssertOn(t)
		c := &ripper.AppConf{
			WorkDirectory: "a/b/c",
			MetaInfoRepo: "x/y/z"}
		assert.NotError(validateConfig(c))
		assert.StringsEqual("a/b/c", c.WorkDirectory)
		assert.StringsEqual("x/y/z", c.MetaInfoRepo)
	})

}
