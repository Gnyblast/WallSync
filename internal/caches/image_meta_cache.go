package caches

import (
	"math/rand/v2"
	"path"
	"path/filepath"
	"strings"

	"github.com/gnyblast/WallSync/internal/engines"
	"github.com/gnyblast/WallSync/internal/models"
)

type ImageMetaCache struct {
	args       models.Args
	images     map[string]string
	rotateList []string
}

func NewImageMetaCache(args models.Args) *ImageMetaCache {
	fsEngine := engines.NewFSEngine(args.OutputDir, args)
	imageCache := &ImageMetaCache{
		args:       args,
		images:     make(map[string]string),
		rotateList: make([]string, 0),
	}

	files := fsEngine.CollectExistingFiles()

	for _, v := range files {
		fileBaseName := strings.TrimSuffix(v, filepath.Ext(v))
		imageCache.images[fileBaseName] = path.Join(args.OutputDir, v)
	}

	if args.MaxImages < len(imageCache.images) {
		var deletionCount int = len(imageCache.images) - args.MaxImages
		for i := 0; i < deletionCount; i++ {
			imageCache.RotateOne()
		}
	}

	return imageCache
}

func (i *ImageMetaCache) GetImagePath(id string) string {
	return i.images[id]
}

func (i *ImageMetaCache) GetRandomImage() string {
	l := rand.IntN(len(i.images))
	it := 0
	randImageId := ""
	for k := range i.images {
		if it == l {
			randImageId = k
		}
		it++
	}

	return randImageId
}

func (i *ImageMetaCache) SetImage(id string, path string) {
	i.images[id] = path
}

func (i *ImageMetaCache) RemoveImage(id string) {
	fsEngine := engines.NewFSEngine(i.images[id], i.args)
	fsEngine.RemoveFile()
	delete(i.images, id)
	if len(i.rotateList) > 0 {
		var index int = 0
		for i, v := range i.rotateList {
			if v == id {
				index = i
				break
			}
		}

		i.rotateList = append(i.rotateList[:index], i.rotateList[index+1:]...)
	}

}

func (i *ImageMetaCache) IsImageExist(id string) bool {
	_, ok := i.images[id]
	return ok
}

func (i *ImageMetaCache) GetImagesLenght() int {
	return len(i.images)
}

func (i *ImageMetaCache) RotateOne() {
	if len(i.rotateList) < 1 {

		i.rotateList = make([]string, 0, len(i.images))

		for id := range i.images {
			i.rotateList = append(i.rotateList, id)
		}

	}

	var index int = rand.IntN(len(i.rotateList))
	var id string = i.rotateList[index]
	i.RemoveImage(id)
}
