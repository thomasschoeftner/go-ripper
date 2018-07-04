package omdb

import (
	"testing"
	"strings"
	"go-ripper/ripper"
)

var allTokens = []string {"sepp", "hat", "gelbe", "eier"}
var conf = ripper.OmdbConfig {
	MovieQuery:   "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}",
	SeriesQuery:  "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}",
	EpisodeQuery: "https://www.omdbapi.com/?apikey={omdbtoken}&i={imdbid}&Season={seasonNo}&Episode={episodeNo}"}

func TestRoundRobinTokenUsage(t *testing.T) {
	f, err := NewOmdbVideoQueryFactory(&conf, allTokens)
	if err != nil {
		t.Errorf("omdb token tFactory failed unexpectedly due to %v", err)
	}
	tokens := f.(*OmdbVideoMetaInfoSource)
	validateToken(t, allTokens[0], tokens.nextToken())
	validateToken(t, allTokens[1], tokens.nextToken())
	validateToken(t, allTokens[2], tokens.nextToken())
	validateToken(t, allTokens[3], tokens.nextToken())
	validateToken(t, allTokens[0], tokens.nextToken())
}

func validateToken(t *testing.T, expected string, got string) {
	if expected != got {
		t.Errorf("expected token \"%s\", but got \"%s\"", expected, got)
	}
}

func TestEmptyTokens(t *testing.T) {
	t.Run("empty tokens", func(t *testing.T) {
		f, err := NewOmdbVideoQueryFactory(&conf, []string{})
		if err == nil || f != nil {
			t.Errorf("expected creation of omdb meta-info query factor to fail due to missing tokens - did not happen")
		}
	})

	t.Run("nil tokens", func(t *testing.T) {
		f, err := NewOmdbVideoQueryFactory(&conf, nil)
		if err == nil || f != nil {
			t.Errorf("expected creation of omdb meta-info query factor to fail due to missing tokens - did not happen")
		}
	})
}

func TestNilConfig(t *testing.T) {
	f, err := NewOmdbVideoQueryFactory(nil, []string{"a", "b", "c"})
	if err == nil || f != nil {
		t.Errorf("expected creation of omdb meta-info query factor to fail due to missing config - did not happen")
	}
}

func TestReplaceVars(t *testing.T) {
	url := replaceUrlVars(conf.MovieQuery, map[string]string {
		urlpattern_omdbtoken : "sepp",
		urlpattern_imdbid : "hatgelbeeier"})
	if strings.Contains(url, "omdbtoken") ||
		strings.Contains(url, "imdbid") ||
		strings.ContainsAny(url, "{}") {
			t.Errorf("variable replacement failed - result is \"%s\"", url)
	}
}

func TestCreateQuery(t *testing.T) {
	t.Run("create title query", func(t *testing.T) {
		f, _ := NewOmdbVideoQueryFactory(&conf, allTokens)
		expectedId := "franzl dauni"
		q := f.NewTitleQuery(expectedId)
		oq := q.(*omdbQuery)

		toBeReplaced := []string {"omdbtoken", "imdbid"}
		if strings.Contains(oq.url, toBeReplaced[0]) || strings.Contains(oq.url, toBeReplaced[1]) {
			t.Errorf("variable replacement failed still got \"%s\" or \"%s\" in \"%s\"", toBeReplaced[0], toBeReplaced[1], oq.url)
		}
		if !strings.Contains(oq.url, expectedId) {
			t.Errorf("variable replacement failed - missing expected imdbid \"%s\" in \"%s\"", expectedId, oq.url)
		}

	})

	t.Run("create episode query", func(t *testing.T) {

	})
}