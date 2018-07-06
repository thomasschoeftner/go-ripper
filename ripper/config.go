package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	IgnorePrefix  string
	WorkDirectory string
	MetaInfoRepo  string
	Processing    *task.ProcessingConfig
	Scan          *ScanConfigGroup
	Resolve       *ResolveConfig
	Tool          *ToolConfig
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
	Resolver string
	Omdb *OmdbConfig
}

type OmdbConfig struct {
	Timeout int
	Retries int
	MovieQuery string
	SeriesQuery string
	EpisodeQuery string
	OmdbTokens []string
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
