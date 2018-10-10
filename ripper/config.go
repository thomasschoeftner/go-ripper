package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	IgnorePrefix  string
	WorkDirectory string
	MetaInfoRepo  string
	Processing    *task.ProcessingConfig
	Output        *OutputConfig
	Scan          *ScanConfigGroup
	Resolve       *ResolveConfig
	Tag           *TagConfig
}

type OutputConfig struct {
	Video string
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

type TagConfig struct {
	Video *VideoTagConfig
}

type VideoTagConfig struct {
	Tagger string
	AtomicParsley *struct {
		Path string
	}
}

type HandbrakeConf struct {
	//TODO
}

type VlcConf struct {
	//TODO
}
