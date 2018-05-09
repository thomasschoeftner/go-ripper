package main

type OmdbConfig struct {
	Omdbtoken *string
	TitleQuery  *string
	SeasonQuery *string
}

const (
	VariableOmdbToken = "omdb.token"
	VariableImdbTitleId = "imdb.id"
	VariableImdbSeasonNo = "imdb.seasonNo"
)

