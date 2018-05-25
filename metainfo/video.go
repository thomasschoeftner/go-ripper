package metainfo

type VideoMetaInfo struct {
	Id       string   //omdb: imdbID
	Title    string   //omdb: Title
	Year     int      //omdb: Year
	Runtime  int      //omdb: Runtime
	Genres   []string //omdb: Genre
	Actors   []string //omdb: Actors
	Writers  []string //omdb: Writer
	Director string   //omdb: Director
	Poster   string   //omdb: Poster
}

func IsVideo(metaInfo MetaInfo) bool {
	return is(metaInfo, mm_type_video)
}

func (vmi *VideoMetaInfo) mediaType() multiMediaType {
	return mm_type_video
}

