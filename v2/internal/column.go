package internal

type ColType string

var (
	ColName ColType = "Name"
	ColVer  ColType = "Version"
	ColDate ColType = "Installed Date"
	ColDesc ColType = "Description"
)

var COLUMNS = []ColType{ColName, ColVer, ColDate, ColDesc}
