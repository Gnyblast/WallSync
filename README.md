# WallSync

## What is WallSync

`WallSync` is an app written in Golang, to perform search queries to <https://wallhaven.cc> through the API and save images to your local machine by making them also background image by using your favorite background manager.

`WallSync` will collect predefined about of images and then stop fetching new images but instead it will keep refreshing the background image from the ones on disk unless you pass `--rotate`.

## Usage

1. Run `go install github.com/gnyblast/WallSync/cmd/wallSync@latest`
    1. Make sure installation complete without error
    2. Make sure the `GOPATH` is in your `PATH`
    3. Construct your command and add it to your start-up script or create a `systemd` file and enable it.
2. Download the `wallSync` binary
    1. Download the `wallSync` binary from the releases.
    2. Move the `wallSync` binary somewhere in your `PATH`
    3. Construct your command and add it to your start-up script or create a `systemd` file and enable it.
3. Compile from source code
    1. Clone the source code `git clone https://github.com/gnyblast/wallsync`
    2. `cd` into the directory top level and run `go build ./cmd/wallSync`
    3. Move the compiled `wallSync` binary somewhere in your `PATH`
    4. Construct your command and add it to your start-up script or create a `systemd` file and enable it.

### Examples

Keep 10 files and upon reaching 10 don't download new:\
`wallSync --search "world of warcraft" --ratios 16X9 --output-dir /tmp --command feh --arguments-template="--bg-fill %s" --refresh 1 --max-images 10`

Keep 10 files and upon reaching 10 start replacing existing files from remote but always keep 10 files:\
`wallSync --search "world of warcraft" --ratios 16X9 --output-dir /tmp --command feh --arguments-template="--bg-fill %s" --refresh 1 --max-images 10 --rotate`

P.S: `--arguments-template` always has to be used with `=` sign and arguments wrapped with quotes("). Check examples above.

## Arguments

`--search`: (required) WallHaven: search query - required\
`--categories`: WallHaven: categories [default: 100]\
`--purity` WallHaven: purity [default: 100]\
`--ratios`: WallHaven: ratios [default: 16X9]\
`--sorting`: WallHaven: sorting [default: relevance]\
`--order`: WallHaven: order [default: desc]\
`--ai-art-filter`: WallHaven: AI Art Filter [default: 1]\
`--refresh`: Refresh background every x Minutes. Default is 10 [default: 10]\
`--output-dir`: (required) Output directory to save image files: there will be a subdirectory created named 'WallSync'\
`--command`: (required) Binary name to be used as background manager. i.e: feh\
`--arguments-template`: (required) Arguments for the command with %s used for image path substitution: i.e: --bg-fill %s\
`--max-images`: (required) Max number of images that will be stored in disk\
`--image-name-prefix`: Prefix for the image names that will be written to disk\
`--rotate`: If set, the program will be rotating the images from remote when reached MaxImages\
`--version`: Print version information