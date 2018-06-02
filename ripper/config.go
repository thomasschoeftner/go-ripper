package ripper

import (
	"go-cli/task"
	"strings"
	"fmt"
)

type AppConf struct {
	IgnoreFolderPrefix  string
	TempDirectoryName   string
	OutputDirectoryName string
	Processing          *task.ProcessingConfig
	Scan                *ScanConfigGroup
	Resolve             *ResolveConfig
	Tool                *ToolConfig
}

func (conf *AppConf) AppendIgnorePrefix() {
	// hardcode ignore-prefix on temp and output dirs to avoid configuration issues
	conf.TempDirectoryName = conf.appendIgnorePrefix(conf.TempDirectoryName)
	conf.OutputDirectoryName = conf.appendIgnorePrefix(conf.OutputDirectoryName)
}

func (conf *AppConf) appendIgnorePrefix(s string) string {
	if strings.HasPrefix(s, conf.IgnoreFolderPrefix) {
		return s
	} else {
		return fmt.Sprintf("%s%s", conf.IgnoreFolderPrefix, s)
	}
}


type ScanConfigGroup struct {
	Video *ScanConfig
}

type ScanConfig struct {
	IdPattern  string
	CollectionPattern string
	ItemNoPattern string
	Patterns []string
}

type ResolveConfig struct {
	Video *VideoResolveConfig
}

type VideoResolveConfig struct {
	Omdb *OmdbConfig
}

type OmdbConfig struct {
	TitleQuery string
	SeasonQuery string
	EpisodeQuery string
}

type ToolConfig struct {
	Handbrake *HandbrakeConf
	Vlc *VlcConf
}

type HandbrakeConf struct {
	//TODO
}

type VlcConf struct {
	//TODO
}
