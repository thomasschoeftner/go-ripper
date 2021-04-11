package omdb

import (
	"testing"
	"strconv"
	"github.com/thomasschoeftner/go-cli/test"
)

const jsonPattern = `
{
  "Title" : "{title}",
  "Year" :  "{year}",
  "Poster" : "{poster}",
  "imdbID" : "{id}",
  "Type" : "{type}",
  "Genre" : "sci-fi",
  "Season" : "{season}",
  "Episode" : "{episode}",
  "totalSeasons" : "{totalseasons}"
}`
var replaceVars = replaceUrlVars

func TestMovieMapping(t *testing.T) {
	vals := map[string]string {
		"title" : "another story",
		"year" : "2012",
		"poster" : "http://a.bcd/image.jpeg",
		"id" : "blah",
		"type" : "movie",
	}

	t.Run("valid movie", func(t *testing.T) {
		raw := []byte(replaceVars(jsonPattern, vals))
		assert := test.AssertOn(t)
		got, err := toMovieMetaInfo(raw)

		assert.NotError(err)
		if nil == got {
			t.Errorf("expected not nil, but got %v", got)
		}
		assert.StringsEqual(vals["title"], got.Title)
		assert.StringsEqual(vals["year"], got.Year)
		assert.StringsEqual(vals["poster"], got.Poster)
		assert.StringsEqual(vals["id"], got.Id)
	})

	t.Run("invalid type", func(t * testing.T) {
		assert := test.AssertOn(t)
		raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"type" : "not a movie"}), vals))
		got, err := toMovieMetaInfo(raw)
		assert.ExpectError("did not catch expected error when mapping movie data without type = \"movie\"")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})


	////obsolete test - all fields are strings -> no mor incorrect field type possible
	//t.Run("incorrect field type", func(t *testing.T) {
	//	raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"year" : "twenty eighteen"}), vals))
	//	got, err := toMovieMetaInfo(raw)
	//	assert := test.AssertOn(t)
	//	assert.ExpectError("did not catch expected error when mapping movie data with invalid content")(err)
	//	if nil != got {
	//		t.Errorf("expected nil, but got %v", got)
	//	}
	//})

	t.Run("missing fields", func(t *testing.T) {
		illformed := []byte(`
			"Title" : "title",
			"Year" :  "2012",
			"imdbID" : "tt23456",
			"Type" : "movie",
		`)
		got, err := toMovieMetaInfo(illformed)
		assert := test.AssertOn(t)
		assert.ExpectError("did not catch expected error when mapping movie data without field \"poster\"")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})
}

func TestSeriesMapping(t *testing.T) {
	vals := map[string]string {
		"title" : "another story",
		"year" : "2012",
		"poster" : "http://a.bcd/image.jpeg",
		"id" : "blah",
		"type" : "series",
		"totalseasons": "17",
	}

	t.Run("valid series", func(t *testing.T) {
		raw := []byte(replaceVars(jsonPattern, vals))
		assert := test.AssertOn(t)
		got, err := toSeriesMetaInfo(raw)

		assert.NotError(err)
		if nil == got {
			t.Errorf("expected not nil, but got %v", got)
		}
		assert.StringsEqual(vals["title"], got.Title)
		assert.StringsEqual(vals["year"], got.Year)
		assert.StringsEqual(vals["poster"], got.Poster)
		assert.StringsEqual(vals["id"], got.Id)
		assert.StringsEqual(vals["totalseasons"], strconv.Itoa(got.Seasons))
	})

	t.Run("invalid type", func(t * testing.T) {
		assert := test.AssertOn(t)
		raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"type" : "not a series"}), vals))
		got, err := toSeriesMetaInfo(raw)
		assert.ExpectError("did not catch expected error when mapping series data without type = \"series\"")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})

	t.Run("incorrect field type", func(t *testing.T) {
		raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"totalseasons" : "three"}), vals))
		got, err := toSeriesMetaInfo(raw)
		assert := test.AssertOn(t)
		assert.ExpectError("did not catch expected error when mapping series data with invalid content")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		illformed := []byte(`
			"Title" : "title",
			"Year" :  "2012",
			"imdbID" : "tt23456",
			"Type" : "series",
		`)
		got, err := toMovieMetaInfo(illformed)
		assert := test.AssertOn(t)
		assert.ExpectError("did not catch expected error when mapping series data without number of total seasons")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})
}

func TestEpisodeMapping(t *testing.T) {
	vals := map[string]string {
		"title" : "another story",
		"year" : "2012",
		"poster" : "http://a.bcd/image.jpeg",
		"id" : "blah",
		"type" : "episode",
		"season": "4",
		"episode" : "9",
	}

	t.Run("valid episode", func(t *testing.T) {
		raw := []byte(replaceVars(jsonPattern, vals))
		assert := test.AssertOn(t)
		got, err := toEpisodeMetaInfo(raw)
		assert.NotError(err)
		if nil == got {
			t.Errorf("expected not nil, but got %v", got)
		}
		assert.StringsEqual(vals["title"], got.Title)
		assert.StringsEqual(vals["year"], got.Year)
		assert.StringsEqual(vals["id"], got.Id)
		assert.StringsEqual(vals["season"], strconv.Itoa(got.Season))
		assert.StringsEqual(vals["episode"], strconv.Itoa(got.Episode))
	})

	t.Run("invalid type", func(t * testing.T) {
		assert := test.AssertOn(t)
		raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"type" : "not an episode"}), vals))
		got, err := toEpisodeMetaInfo(raw)
		assert.ExpectError("did not catch expected error when mapping episode data without type = \"episode\"")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})

	t.Run("incorrect field type", func(t *testing.T) {
		raw := []byte(replaceVars(replaceVars(jsonPattern, map[string]string {"year" : "twenty eighteen"}), vals))
		got, err := toSeriesMetaInfo(raw)
		assert := test.AssertOn(t)
		assert.ExpectError("did not catch expected error when mapping episode data with invalid content")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		illformed := []byte(`
			"Title" : "title",
			"Year" :  "2012",
			"imdbID" : "tt23456",
			"Type" : "episode",
			"season": "4",
		`)
		got, err := toMovieMetaInfo(illformed)
		assert := test.AssertOn(t)
		assert.ExpectError("did not catch expected error when mapping series data without number of total seasons")(err)
		if nil != got {
			t.Errorf("expected nil, but got %v", got)
		}
	})
}
