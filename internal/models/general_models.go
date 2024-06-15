package models

type VersionArg struct {
	Version bool `json:"-" help:"Print version information"`
}
type Args struct {
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
}

type WallHavenAPIModel struct {
	Data []WallHavenData `json:"data"`
	Meta WallHavenMeta   `json:"meta"`
}

type WallHavenMeta struct {
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	PerPage     int    `json:"per_page"`
	Total       int    `json:"total"`
	Query       string `json:"query"`
}

type WallHavenData struct {
	ID         string              `json:"id"`
	URL        string              `json:"url"`
	ShortURL   string              `json:"short_url"`
	Views      int                 `json:"views"`
	Favorites  int                 `json:"favorites"`
	Source     string              `json:"source"`
	Purity     string              `json:"purity"`
	Category   string              `json:"category"`
	DimensionX int                 `json:"dimension_x"`
	DimensionY int                 `json:"dimension_y"`
	Resolution string              `json:"resolution"`
	Ratio      string              `json:"ratio"`
	FileSize   int                 `json:"file_size"`
	FileType   string              `json:"file_type"`
	CreatedAt  string              `json:"created_at"`
	Colors     []string            `json:"colors"`
	Path       string              `json:"path"`
	Thumbs     WallHavenDataThumbs `json:"thumbs"`
}

type WallHavenDataThumbs struct {
	Large    string `json:"large"`
	Original string `json:"original"`
	Small    string `json:"small"`
}
