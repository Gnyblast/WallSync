package services

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gnyblast/WallSync/internal/caches"
	"github.com/gnyblast/WallSync/internal/models"
)

type BackgroundImageService struct {
	imageCache  *caches.ImageMetaCache
	imageUpdate <-chan string
	Quit        chan bool
	args        models.Args
}

func NewBackgroundImageService(args models.Args, imageUpdate chan string, imageCache *caches.ImageMetaCache) *BackgroundImageService {
	return &BackgroundImageService{
		args:        args,
		imageCache:  imageCache,
		imageUpdate: imageUpdate,
		Quit:        make(chan bool),
	}
}

func (b BackgroundImageService) Listen() {
	go b.ListenImageUpdates()
	<-b.Quit
	close(b.Quit)
}

func (b BackgroundImageService) ListenImageUpdates() {
	for {
		randImg := <-b.imageUpdate
		b.setWallPaper(randImg)
	}
}

func (b *BackgroundImageService) setWallPaper(imageId string) {
	log.Println("Setting wallpaper")
	o, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("command -v %s", b.args.Command)).Output()
	if err != nil {
		log.Fatalf("Feh is not installed or cannot be found: %v", err)
	}

	log.Println(b.imageCache.GetImagePath(imageId))
	arguments := fmt.Sprintf(b.args.ArgumentsTemplate, b.imageCache.GetImagePath(imageId))
	_, err = exec.Command(strings.TrimSuffix(string(o), "\n"), strings.Split(arguments, " ")[0:]...).Output()
	if err != nil {
		log.Fatalf("Could not set background: %v", err)
	}
}
