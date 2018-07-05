package omdb

import (
	"testing"
	"strings"
	"go-ripper/ripper"
)

var allTokens = []string {"sepp", "hat", "gelbe", "eier"}

func conf(tokens []string) *ripper.VideoResolveConfig {
	return &ripper.VideoResolveConfig{
		Resolver: CONF_OMDB_RESOLVER,
		Omdb: &ripper.OmdbConfig {
			OmdbTokens:   tokens,
			MovieQuery:   "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}",
			SeriesQuery:  "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}",
			EpisodeQuery: "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}&Season={seasonNo}&Episode={episodeNo}"},
	}
}

func TestRoundRobinTokenUsage(t *testing.T) {
	f, err := NewOmdbVideoMetaInfoSource(conf(allTokens))
	if err != nil {
		t.Errorf("omdb token tFactory failed unexpectedly due to %v", err)
	}
	tokens := f.(*omdbVideoMetaInfoSource)
	for i:=0; i<len(allTokens); i++ {
		validateToken(t, allTokens[i], tokens.nextToken())
	}
	validateToken(t, allTokens[0], tokens.nextToken()) //validate round-robin
}

func validateToken(t *testing.T, expected string, got string) {
	if expected != got {
		t.Errorf("expected token \"%s\", but got \"%s\"", expected, got)
	}
}

func TestEmptyTokens(t *testing.T) {
	t.Run("empty tokens", func(t *testing.T) {
		f, err := NewOmdbVideoMetaInfoSource(conf([]string{}))
		if err == nil || f != nil {
			t.Errorf("expected creation of omdb meta-info query factor to fail due to missing tokens - did not happen")
		}
	})

	t.Run("nil tokens", func(t *testing.T) {
		f, err := NewOmdbVideoMetaInfoSource(conf(nil))
		if err == nil || f != nil {
			t.Errorf("expected creation of omdb meta-info query factor to fail due to missing tokens - did not happen")
		}
	})
}

func TestNilConfig(t *testing.T) {
	f, err := NewOmdbVideoMetaInfoSource(nil)
	if err == nil || f != nil {
		t.Errorf("expected creation of omdb meta-info query factor to fail due to missing config - did not happen")
	}
}

func TestReplaceVars(t *testing.T) {
	url := replaceUrlVars(conf(nil).Omdb.MovieQuery, map[string]string {
		urlpattern_omdbtoken : "sepp",
		urlpattern_imdbid : "hatgelbeeier"})
	if strings.Contains(url, "omdbtoken") ||
		strings.Contains(url, "imdbid") ||
		strings.ContainsAny(url, "{}") {
			t.Errorf("variable replacement failed - result is \"%s\"", url)
	}
}
