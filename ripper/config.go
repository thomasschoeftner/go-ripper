package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	TempDirectoryName   string
	OutputDirectoryName string
	Processing          *task.ProcessingConfig
	Scan                *ScanConfigGroup
	Omdb                *OmdbConfig
	Tool                *ToolConfig

}

type ScanConfigGroup struct {
	Video *ScanConfig
	//TODO add: Audio *ScanConfig?
}

type ScanConfig struct {
	IdPattern  string
	CollectionPattern string
	ItemNoPattern string
	Patterns []string
}

type ToolConfig struct {
	Handbrake *HandbrakeConf
	Vlc *VlcConf
}

type OmdbConfig struct {
	OmdbTokens []string
	TitleQuery string
	SeasonQuery string
}

type HandbrakeConf struct {
	//TODO
}

type VlcConf struct {
	//TODO
}
