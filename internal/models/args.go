package models

type ServiceArgs struct {
	ServiceName string `json:"-" arg:"--service-name,required" help:"Name of the service to connect. I.e: --service-name wallhaven"`
	Version     bool   `json:"-" help:"Print version information"`
}
type VersionArg struct {
	Version bool `json:"-" help:"Print version information"`
}
type WallhavenArgs struct {
	Search            string `json:"q" arg:"required" help:"WallHaven: search query"`
	Categories        int    `json:"categories" default:"100" help:"WallHaven: categories"`
	Purity            int    `json:"purity" default:"100" help:"WallHaven: purity"`
	Ratios            string `json:"ratios" default:"16X9" help:"WallHaven: ratios"`
	Sorting           string `json:"sorting" default:"relevance" help:"WallHaven: sorting"`
	Order             string `json:"order" default:"desc" help:"WallHaven: order"`
	AIArtFilter       int    `json:"ai_art_filter" arg:"--ai-art-filter" default:"1" help:"WallHaven: AI Art Filter"`
	Refresh           int    `json:"-" default:"10" help:"Refresh background every x Minutes. Default is 10"`
	Page              int    `json:"page" arg:"-"`
	OutputDir         string `json:"-" arg:"--output-dir,required" help:"Output directory to save image files: there will be a subdirectory created named 'WallSync'"`
	Command           string `json:"-" arg:"required" help:"Binary name to be used as background manager. i.e: feh"`
	ArgumentsTemplate string `json:"-" arg:"--arguments-template,required" help:"Arguments for the command with %s used for image path substitution: i.e: --bg-fill %s"`
	MaxImages         int    `json:"-" arg:"--max-images,required" help:"Max number of images that will be stored in disk"`
	ImageNamePrefix   string `json:"-" arg:"--image-name-prefix" help:"Prefix for the image names that will be written to disk"`
	Rotate            bool   `json:"-" help:"If set, the program will be rotating the images from remote when reached MaxImages"`
	Version           bool   `json:"-" help:"Print version information"`
	ServiceName       string `json:"-" arg:"--service-name,required" help:"Name of the service to connect. I.e: --service-name wallhaven"`
}

func (w WallhavenArgs) GetRefreshRate() int {
	return w.Refresh
}
func (w WallhavenArgs) GetOutputDir() string {
	return w.OutputDir
}

func (w *WallhavenArgs) SetOutputDir(outputDir string) {
	w.OutputDir = outputDir
}
func (w WallhavenArgs) GetCommand() string {
	return w.Command
}
func (w WallhavenArgs) GetArgumentsTemplate() string {
	return w.ArgumentsTemplate
}
func (w WallhavenArgs) GetMaxImages() int {
	return w.MaxImages
}
func (w WallhavenArgs) GetImageNamePrefix() string {
	return w.ImageNamePrefix
}
func (w WallhavenArgs) ShouldRotate() bool {
	return w.Rotate
}

func (w WallhavenArgs) ShouldPrintVersion() bool {
	return w.Version
}
