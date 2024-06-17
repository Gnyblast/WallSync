package interfaces

type IArgs interface {
	GetRefreshRate() int
	GetOutputDir() string
	SetOutputDir(string)
	GetCommand() string
	GetArgumentsTemplate() string
	GetMaxImages() int
	GetImageNamePrefix() string
	ShouldRotate() bool
	ShouldPrintVersion() bool
}
