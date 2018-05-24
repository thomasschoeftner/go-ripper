package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	TempDirectoryName   string
	OutputDirectoryName string
	Processing          *task.ProcessingConfig
	Scan                *ScanConfigGroup
	Resolve             *ResolveConfig
	Tool                *ToolConfig
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
