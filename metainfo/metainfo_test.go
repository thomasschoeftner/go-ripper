package metainfo

import (
	"testing"
	"go-cli/test"
	"path/filepath"
)

const testMetaInfoType = "testmetainfo"
type testMetaInfo struct {
	IdInfo
	Field string
}
func (tmi *testMetaInfo) GetType() string {
	return testMetaInfoType
}

func TestTypeComparison(t *testing.T) {
	assert := test.AssertOn(t)
	assert.False("expected false when comparing meta-info type with nil meta-info")(Is(nil, "sometype"))
	assert.False("expected false when comparing meta-info type with different type")(Is(&testMetaInfo{IdInfo{"id"}, "value"}, "other-type"))
	assert.True("expected true when comparing meta-info type with correct type")(Is(&testMetaInfo{IdInfo{"id"}, "value"}, testMetaInfoType))
}

func TestReadSaveMetaInfo(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	fName := filepath.Join(dir, "meta-info.json")
	mi := testMetaInfo{IdInfo{"abcd"}, "some value"}
	assert := test.AssertOn(t)
	assert.NotError(SaveMetaInfo(fName, &mi))
	gotMi := testMetaInfo{}
	assert.NotError(ReadMetaInfo(fName, &gotMi))
	assert.StringsEqual(mi.Id, gotMi.Id)
	assert.StringsEqual(mi.Field, gotMi.Field)
}
