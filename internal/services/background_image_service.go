package services

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gnyblast/WallSync/internal/caches"
	"github.com/gnyblast/WallSync/internal/interfaces"
)

type BackgroundImageService struct {
	imageCache  *caches.ImageMetaCache
	imageUpdate <-chan string
	Quit        chan bool
	args        interfaces.IArgs
}

func NewBackgroundImageService(args interfaces.IArgs, imageUpdate chan string, imageCache *caches.ImageMetaCache) *BackgroundImageService {
	return &BackgroundImageService{
		args:        args,
		imageCache:  imageCache,
		imageUpdate: imageUpdate,
		Quit:        make(chan bool),
	}
}

func (b BackgroundImageService) Listen() {
loop:
	for {
		select {
		case randImg := <-b.imageUpdate:
			b.setWallPaper(randImg)
		case <-b.Quit:
			break loop
		}
	}
	close(b.Quit)
}

func (b *BackgroundImageService) setWallPaper(imageId string) {
	log.Println("Setting wallpaper")
	o, err := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", b.args.GetCommand())).Output()
	if err != nil {
		log.Fatalf("Feh is not installed or cannot be found: %v", err)
	}

	imagePath, err := b.imageCache.GetImagePath(imageId)
	if err != nil {
		log.Printf("Could not find image: %v", err)
	}

	log.Println(imagePath)
	arguments := fmt.Sprintf(b.args.GetArgumentsTemplate(), imagePath)
	res, err := exec.Command(strings.TrimSuffix(string(o), "\n"), strings.Split(arguments, " ")[0:]...).Output()
	if err != nil {
		log.Fatalf("Could not set background: %v, %s", err, string(res))
	}
}
