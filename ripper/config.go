package ripper

import (
	"go-cli/task"
)

type AppConf struct {
	WorkDirectory string
	OutputDirectory string
	Processing *task.ProcessingConf
	Omdb *OmdbConf
	Tool *ToolConfig
}

type ToolConfig struct {
	Handbrake *HandbrakeConf
	Vlc *VlcConf
}

type OmdbConf struct {
	TitleQuery string
	SeasonQuery string
}

type HandbrakeConf struct {
	//TODO
}

type VlcConf struct {
	//TODO
}
