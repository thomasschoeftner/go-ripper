package ripper

import (
	"go-cli/task"
	"go-cli/require"
	"go-cli/config"
	"fmt"
	"strings"
)

func GetConfig(configFile string) *AppConf {
	conf := AppConf{}
	require.NotFailed(config.FromFile(&conf, configFile, map[string]string {}))
	require.NotFailed(validateConfig(&conf))
	return &conf
}

func validateConfig(c *AppConf) error {
	if c == nil {
		return fmt.Errorf("config is nil - no config available")
	}

	validatePath:= func(path string, fieldName string) error {
		if 0 == len(path) {
			return fmt.Errorf("[config error] \"%s\" is empty", fieldName)
		}
		if strings.ContainsRune(path, ' ') {
			return fmt.Errorf("[config error] \"%s\" must not contain spaces", fieldName)
		}
		return nil
	}

	c.WorkDirectory = strings.Trim(c.WorkDirectory, " ")
	if err := validatePath(c.WorkDirectory, "workDirectory"); err != nil {
		return err
	}

	//validate metainforepo
	c.MetaInfoRepo = strings.Trim(c.MetaInfoRepo, " ")
	if err := validatePath(c.MetaInfoRepo, "metaInfoRepo"); err != nil {
		return err
	}

	//validate outputfolder
	c.DefaultOutputDirectory = strings.Trim(c.DefaultOutputDirectory, " ")
	if err := validatePath(c.DefaultOutputDirectory, "defaultOutputDirectory"); err != nil {
		return err
	}

	return nil
}


type AppConf struct {
	IgnorePrefix           string
	WorkDirectory          string
	MetaInfoRepo           string
	DefaultOutputDirectory string
	Processing             *task.ProcessingConfig
	Output                 *OutputConfig
	Scan                   *ScanConfigGroup
	Resolve                *ResolveConfig
	Tag                    *TagConfig
}

type OutputConfig struct {
	Video string
}

type ScanConfigGroup struct {
	Video *ScanConfig
}

type ScanConfig struct {
	IdPattern         string
	CollectionPattern string
	ItemNoPattern     string
	Patterns          []string
	AllowSpaces       bool
	AllowedExtensions []string
}

type ResolveConfig struct {
	Video *VideoResolveConfig
}

type VideoResolveConfig struct {
	Resolver string
	Omdb     *OmdbConfig
}

type OmdbConfig struct {
	Timeout      int
	Retries      int
	MovieQuery   string
	SeriesQuery  string
	EpisodeQuery string
	OmdbTokens   []string
}

type TagConfig struct {
	Video *VideoTagConfig
}

type VideoTagConfig struct {
	Tagger        string
	AtomicParsley *AtomicParsleyConfig
}

type AtomicParsleyConfig struct {
	Path    string
	Timeout string
}

type HandbrakeConfig struct {
	//TODO
}
