package services

import (
	"encoding/json"
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
	args                  models.Args
	imageCache            *caches.ImageMetaCache
	metaUpdateTicker      *time.Ticker
	wallpaperUpdateTicker *time.Ticker
	imageUpdate           chan<- string
	Quit                  chan bool
}

func WallHavenServiceService(args models.Args, imageUpdate chan string, imageCache *caches.ImageMetaCache) *WallHavenService {
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
	w.wallpaperUpdateTicker = time.NewTicker(time.Duration(w.args.Refresh) * time.Minute)
	w.getMetaData()
	w.wallpaperUpdater()

	go w.ListenMetaTicker()
	go w.ListeWallpaperTicker()

	<-w.Quit
	w.metaUpdateTicker.Stop()
	w.wallpaperUpdateTicker.Stop()
	close(w.Quit)
}

func (w *WallHavenService) ListenMetaTicker() {
	for {
		<-w.metaUpdateTicker.C
		w.getMetaData()
	}
}

func (w *WallHavenService) ListeWallpaperTicker() {
	for {
		<-w.wallpaperUpdateTicker.C
		w.wallpaperUpdater()
	}
}

func (w *WallHavenService) wallpaperUpdater() {
	if w.args.MaxImages > w.imageCache.GetImagesLenght() || w.args.Rotate {
		image := w.getRandomImage()
		if !w.imageCache.IsImageExist(image.ID) {
			imagePath := utils.ConstructImagePath(image.ID, w.args.ImageNamePrefix, w.args.OutputDir, image.FileType)
			fsEngine := engines.NewFSEngine(imagePath, w.args)
			fsEngine.WriteImageToAFile(w.fetchImage(image))
			w.imageCache.SetImage(image.ID, imagePath)
			w.imageUpdate <- image.ID
		}
		if w.args.MaxImages < w.imageCache.GetImagesLenght() && w.args.Rotate {
			w.imageCache.RotateOne()
		}
		return
	}

	w.imageUpdate <- w.imageCache.GetRandomImage()
}

func (w *WallHavenService) getMetaData() {
	response := w.requestEngine.DoRequest(http.MethodGet, utils.CreateUrlQuery(w.args))
	var wallHaveApiResponse models.WallHavenAPIModel
	err := json.Unmarshal(response, &wallHaveApiResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal wallhaven response: %v", err)
	}

	w.metaData = wallHaveApiResponse.Meta
}

func (w WallHavenService) getRandomImage() models.WallHavenData {
	log.Println("Downloading image")
	var newArgs models.Args = w.args
	newArgs.Page = rand.IntN(w.metaData.LastPage-1) + 1
	response := w.requestEngine.DoRequest(http.MethodGet, utils.CreateUrlQuery(newArgs))
	var wallHaveApiResponse models.WallHavenAPIModel
	err := json.Unmarshal(response, &wallHaveApiResponse)
	if err != nil {
		log.Fatalf("Error fetching image list: %v", err)
	}

	selectedImage := rand.IntN(wallHaveApiResponse.Meta.PerPage)
	return wallHaveApiResponse.Data[selectedImage]
}

func (w WallHavenService) fetchImage(selectedImage models.WallHavenData) []byte {
	log.Println("Fetching image")
	requestEngine := engines.NewHttpRequestEngine(selectedImage.Path, 0)
	response := requestEngine.DoRequest(http.MethodGet, "")
	return response
}
