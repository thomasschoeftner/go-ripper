package omdb

import (
	"encoding/json"
	"strings"
	"fmt"
	"strconv"
	"go-ripper/metainfo/video"
)

const (
	omdb_id      = "imdbid"
	omdb_title   = "title"
	omdb_year    = "year"
	omdb_poster  = "poster"
	omdb_seasons = "totalseasons"
	omdb_season  = "season"
	omdb_episode = "episode"
	omdb_type    = "type"
)

const (
	omdb_type_movie   = "movie"
	omdb_type_series  = "series"
	omdb_type_episode = "episode"
)

func toMap(raw []byte) (map[string]string, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		return nil, err
	}

	results := map[string]string {}
	for k, v := range parsed {
		results[strings.ToLower(k)] = fmt.Sprintf("%v", v)
	}
	return results, nil
}


func toMovieMetaInfo(raw []byte) (*video.MovieMetaInfo, error) {
	values, err := toMap(raw)
	if err != nil {
		return nil, err
	}

	kind := values[omdb_type]
	if omdb_type_movie  != kind {
		return nil, fmt.Errorf("mapping omdb-response to movie meta-info failed: expected type %s, but got type %s", omdb_type_movie, kind)
	}

	movie := &video.MovieMetaInfo{}
	if err := assignString(&movie.Id, values, omdb_id); err != nil {
		return nil, err
	}
	if err := assignString(&movie.Poster, values, omdb_poster); err != nil {
		return nil, err
	}
	if err := assignString(&movie.Title, values, omdb_title); err != nil {
		return nil, err
	}
	if err := assignString(&movie.Year, values, omdb_year); err != nil {
		return nil, err
	}
	return movie, nil
}

func toSeriesMetaInfo(raw []byte) (*video.SeriesMetaInfo, error) {
	values, err := toMap(raw)
	if err != nil {
		return nil, err
	}

	kind := values[omdb_type]
	if omdb_type_series != kind {
		return nil, fmt.Errorf("mapping omdb-response to series meta-info failed: expected type %s, but got type %s", omdb_type_series, kind)
	}

	series := &video.SeriesMetaInfo{}
	if err := assignString(&series.Id, values, omdb_id); err != nil {
		return nil, err
	}
	if err := assignString(&series.Poster, values, omdb_poster); err != nil {
		return nil, err
	}
	if err := assignString(&series.Title, values, omdb_title); err != nil {
		return nil, err
	}
	if err := assignString(&series.Year, values, omdb_year); err != nil {
		return nil, err
	}
	if err := assignInt(&series.Seasons, values, omdb_seasons); err != nil {
		return nil, err
	}
	return series, nil
}

func toEpisodeMetaInfo(raw []byte) (*video.EpisodeMetaInfo, error) {
	values, err := toMap(raw)
	if err != nil {
		return nil, err
	}

	kind := values[omdb_type]
	if omdb_type_episode != kind {
		return nil, fmt.Errorf("mapping omdb-response to episode meta-info failed: expected type %s, but got type %s", omdb_type_episode, kind)
	}

	episode := &video.EpisodeMetaInfo{}
	if err := assignString(&episode.Id, values, omdb_id); err != nil {
		return nil, err
	}
	if err := assignString(&episode.Title, values, omdb_title); err != nil {
		return nil, err
	}
	if err := assignString(&episode.Year, values, omdb_year); err != nil {
		return nil, err
	}
	if err := assignInt(&episode.Season, values, omdb_season); err != nil {
		return nil, err
	}
	if err := assignInt(&episode.Episode, values, omdb_episode); err != nil {
		return nil, err
	}
	return episode, nil
}

func assignString(target *string, values map[string]string, key string) error {
	if val, defined := values[key]; defined {
		*target = val
		return nil
	} else {
		return fmt.Errorf("mapping omdb-response to meta-info failed: omdb field \"%s\" is missing", key)
	}
}

func assignInt(target *int, values map[string]string, key string) error {
	if val, defined := values[key]; defined {
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*target = i
		return nil
	} else {
		return fmt.Errorf("mapping omdb-response to meta-info failed: omdb field \"%s\" is missing", key)
	}
}
