package engines

import (
	"log"
	"os"

	"github.com/gnyblast/WallSync/internal/models"
)

type FSEngine struct {
	path string
	args models.Args
}

func NewFSEngine(path string, args models.Args) *FSEngine {
	return &FSEngine{
		path: path,
		args: args,
	}
}

func (f FSEngine) WriteImageToAFile(imageContent []byte) {
	log.Println("Writing image")

	file, err := os.Create(f.path)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	file.Write(imageContent)
}

func (f FSEngine) CollectExistingFiles() []string {
	var files []string

	d, err := os.ReadDir(f.path)
	if err != nil {
		log.Fatalf("Could not open directory: %v", err)
	}

	for _, v := range d {
		if !v.IsDir() {
			files = append(files, v.Name())
		}
	}

	return files
}

func (f FSEngine) CreateDirectory() {
	err := os.MkdirAll(f.path, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
}

func (f FSEngine) RemoveFile() {
	err := os.Remove(f.path)
	if err != nil {
		log.Fatalf("Error removing file: %v", err)
	}
}
