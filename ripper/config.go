package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	TempDirectoryName   string
	OutputDirectoryName string
	Processing          *task.ProcessingConf
	Omdb                *OmdbConf
	Tool                *ToolConfig
}

type ToolConfig struct {
	Handbrake *HandbrakeConf
	Vlc *VlcConf
}

type OmdbConf struct {
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
