package omdb

import (
	"mmlib/config"
)

type Omdb struct {
	omdbToken string
	titleQueryPattern string
	seasonQueryPattern string
}

func Init(omdbToken string, conf *config.Config) Omdb {
	return Omdb {
		omdbToken: omdbToken,
		titleQueryPattern:  conf.Omdb.TitleQuery.ReplaceVariable(config.VariableOmdbToken, omdbToken).Value(),
		seasonQueryPattern: conf.Omdb.SeasonQuery.ReplaceVariable(config.VariableOmdbToken, omdbToken).Value(),
	}
}



func (Omdb) GetTitleMetaData(imdbId string) {
	//TODO
}

func (Omdb) GetSeasonMetaData(imdbId string, season int) {
	//TODO
}