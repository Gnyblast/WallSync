package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gnyblast/WallSync/internal/caches"
	"github.com/gnyblast/WallSync/internal/engines"
	"github.com/gnyblast/WallSync/internal/models"
	"github.com/gnyblast/WallSync/internal/utils"
)

type WallHavenService struct {
	requestEngine         *engines.HttpRequestEngine
	metaData              models.WallHavenMeta
	args                  *models.WallhavenArgs
	imageCache            *caches.ImageMetaCache
	metaUpdateTicker      *time.Ticker
	wallpaperUpdateTicker *time.Ticker
	imageUpdate           chan<- string
	Quit                  chan bool
}

func NewWallHavenService(args *models.WallhavenArgs, imageUpdate chan string, imageCache *caches.ImageMetaCache) *WallHavenService {
	args.Page = 1
	return &WallHavenService{
		requestEngine: engines.NewHttpRequestEngine("https://wallhaven.cc/api/v1/search", 0),
		args:          args,
		imageUpdate:   imageUpdate,
		imageCache:    imageCache,
		Quit:          make(chan bool),
	}
}

func (w *WallHavenService) Start() {
	w.metaUpdateTicker = time.NewTicker(2 * time.Minute)
	w.wallpaperUpdateTicker = time.NewTicker(time.Duration(w.args.GetRefreshRate()) * time.Minute)
	w.getMetaData()
	w.wallpaperUpdater()

loop:
	for {
		select {
		case <-w.metaUpdateTicker.C:
			w.getMetaData()
		case <-w.wallpaperUpdateTicker.C:
			w.wallpaperUpdater()
		case <-w.Quit:
			w.metaUpdateTicker.Stop()
			w.wallpaperUpdateTicker.Stop()
			break loop
		}
	}
	close(w.Quit)
}

func (w *WallHavenService) Stop() {
	w.Quit <- true
}

func (w *WallHavenService) IsStoped() <-chan bool {
	return w.Quit
}

func (w *WallHavenService) wallpaperUpdater() {
	if (w.args.GetMaxImages() > w.imageCache.GetImagesLenght() || w.args.ShouldRotate()) && w.metaData.LastPage > 0 {
		image, err := w.getRandomImage()
		if err != nil {
			if w.imageCache.GetImagesLenght() > 0 {
				w.imageUpdate <- w.imageCache.GetRandomImage()
			}
			return
		}

		if !w.imageCache.IsImageExist(image.ID) {
			imagePath := utils.ConstructImagePath(image.ID, w.args.GetImageNamePrefix(), w.args.GetOutputDir(), image.FileType)
			fsEngine := engines.NewFSEngine(imagePath, w.args)
			imageContent, err := w.fetchImage(image)
			if err != nil {
				w.imageUpdate <- w.imageCache.GetRandomImage()
				return
			}

			fsEngine.WriteImageToAFile(imageContent)
			w.imageCache.SetImage(image.ID, imagePath)
		}

		w.imageUpdate <- image.ID

		if w.args.GetMaxImages() < w.imageCache.GetImagesLenght() && w.args.ShouldRotate() {
			w.imageCache.RotateOne()
		}

		return
	}

	w.imageUpdate <- w.imageCache.GetRandomImage()
}

func (w *WallHavenService) getMetaData() error {
	response, err := w.requestEngine.DoRequest(http.MethodGet, utils.CreateUrlQuery(w.args))
	if err != nil {
		return fmt.Errorf("failed to get meta data: %v", err)
	}

	var wallHaveApiResponse models.WallHavenAPIModel
	err = json.Unmarshal(response, &wallHaveApiResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal wallhaven response: %v", err)
	}

	w.metaData = wallHaveApiResponse.Meta
	return nil
}

func (w WallHavenService) getRandomImage() (models.WallHavenData, error) {
	log.Println("Downloading image")
	var newArgs *models.WallhavenArgs = w.args
	newArgs.Page = rand.IntN(w.metaData.LastPage)
	response, err := w.requestEngine.DoRequest(http.MethodGet, utils.CreateUrlQuery(newArgs))
	if err != nil {
		return models.WallHavenData{}, err
	}

	var wallHaveApiResponse models.WallHavenAPIModel
	err = json.Unmarshal(response, &wallHaveApiResponse)
	if err != nil {
		return models.WallHavenData{}, err
	}

	selectedImage := rand.IntN(wallHaveApiResponse.Meta.PerPage)
	return wallHaveApiResponse.Data[selectedImage], nil
}

func (w WallHavenService) fetchImage(selectedImage models.WallHavenData) ([]byte, error) {
	log.Println("Fetching image")
	requestEngine := engines.NewHttpRequestEngine(selectedImage.Path, 0)
	return requestEngine.DoRequest(http.MethodGet, "")
}
