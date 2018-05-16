package scan

type target struct {
	Folder string
	File string
	Id string
	Collection *string //album, series
	Group *int //cd#, season#
}

func newTarget(folder string, file string, id string, collection *string, group *int) *target {
	return &target{folder, file, id, collection, group}
}
