package processor

import (
	"testing"
	"go-ripper/targetinfo"
	"go-cli/test"
	"go-cli/commons"
)


var ti targetinfo.TargetInfo = targetinfo.NewMovie("c.foo", "a/b", "id")

func TestDefaultCheckLazy(t *testing.T) {
	t.Run("recommend not lazy if lazy is off", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(false, "foo")
		test.AssertOn(t).False("expected checklazy to be false")(checkLazy(ti))
	})

	t.Run("recommend not lazy if lazy is on but extension does not match", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(true, "bar")
		test.AssertOn(t).False("expected checklazy to be false")(checkLazy(ti))
	})

	t.Run("recommend lazy if lazy is on and extension matches", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(true, "foo")
		test.AssertOn(t).True("expected checklazy to be true")(checkLazy(ti))
	})
}

func TestNeverLazy(t *testing.T) {
	t.Run("recommend not lazy if lazy is off", func(t *testing.T) {
		checkLazy := NeverLazy(false, "procName", commons.Printf)
		test.AssertOn(t).False("expected checklazy to return false when never-lazy")(checkLazy(ti))
	})

	t.Run("recommend not lazy if lazy is on", func(t *testing.T) {
		checkLazy := NeverLazy(true, "procName", commons.Printf)
		test.AssertOn(t).False("expected checklazy to return false when never-lazy")(checkLazy(ti))
	})
}