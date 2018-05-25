package metainfo

import (
	"io/ioutil"
	"io"
)

type multiMediaType string
const (
	mm_type_video multiMediaType = "video"
)


type MetaInfo interface {
	mediaType() multiMediaType
}

func is(metaInfo MetaInfo, kind multiMediaType) bool {
	return kind == metaInfo.mediaType()
}


type metaInfoQuery interface {
	invoke() (io.ReadCloser, error)
	convert(raw []byte) (MetaInfo, error)
}

func Get(miq metaInfoQuery) (MetaInfo, error) {
	readCloser, err := miq.invoke()
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	raw, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}

	return miq.convert(raw)
}
