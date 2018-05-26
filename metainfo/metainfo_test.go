package metainfo

import (
	"io"
	"testing"
	"strings"
)

type TestMetaInfo struct {
	response string
}
const mm_type_test = "test"

func (tm *TestMetaInfo) mediaType() multiMediaType {
	return mm_type_test
}

func TestMetaInfoIs(t *testing.T) {
	tmi := &TestMetaInfo{}
	if is(tmi, mm_type_video) {
		t.Error("test meta info must not match type video")
	}
	if !is(tmi, mm_type_test) {
		t.Error("test meta info must not match type video")
	}
}


type TestMetaInfoQuery struct { //is also an io.ReadCloser
	message string
	toAppend string
	stillOpen bool
	reader io.Reader
}

func (tq *TestMetaInfoQuery) Invoke() (io.ReadCloser, error) {
	tq.stillOpen = true
	tq.reader = strings.NewReader(tq.message+tq.toAppend)
	return tq, nil
}

func (tq *TestMetaInfoQuery) Convert(raw []byte) (MetaInfo, error) {
	return &TestMetaInfo{string(raw)}, nil
}

func (tq *TestMetaInfoQuery) Read(p []byte) (n int, err error) {
	return tq.reader.Read(p)
}

func (tq *TestMetaInfoQuery) Close() error {
	tq.stillOpen = false
	return nil
}


func TestGet(t *testing.T) {
	tmiq := &TestMetaInfoQuery{message: "just some", toAppend: " random text"}
	expected := tmiq.message + tmiq.toAppend
	mi, err := Get(tmiq)
	if err != nil {
		t.Error("querying test meta info raised unexpected error: ", err)
	}
	if tmiq.stillOpen {
		t.Error("MetaInfoQuery ReadWriter was not closed")
	}
	if !is(mi, mm_type_test) {
		t.Errorf("TestMetaInfoQuery result was expected to be of type TestMetaInfo, but was %s", mi.mediaType())
	}
	tmi := mi.(*TestMetaInfo)
	if tmi.response != expected {
		t.Errorf("TestMetaInfoQuery yielded unexpected result - expected \"%s\", but got \"%s\"", expected, tmi.response)
	}
}
