package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/url"
	"os"
	"os/exec"
	"path"

	"github.com/gnyblast/WallSync/internal/constants"
)

func CheckDependencies(command string) {
	_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("command -v %s", command)).Output()
	if err != nil {
		log.Fatalf("Feh is not installed or cannot be found: %v", err)
	}
}

func CreateUrlQuery(args any) string {
	argsJson, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("Error while parsing url query: %v", err)
	}

	var queryMap map[string]interface{} = make(map[string]interface{})

	err = json.Unmarshal(argsJson, &queryMap)
	if err != nil {
		log.Fatalf("Error while parsing url query: %v", err)
	}

	params := url.Values{}
	for k, v := range queryMap {
		params.Add(k, fmt.Sprintf("%v", v))
	}

	return params.Encode()
}

func RandToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func ConstructImagePath(imageId string, imageNamePrefix string, outputDir string, fileType string) string {
	ext, err := mime.ExtensionsByType(fileType)
	if err != nil || len(ext) == 0 {
		log.Fatalf("Could not find file extension: %v", err)
	}

	var imageName string = imageId + ext[0]
	if len(imageNamePrefix) > 0 {
		imageName = imageNamePrefix + "_" + imageId + ext[0]
	}
	imagePath := path.Join(outputDir, imageName)
	return imagePath
}

func PrintVersionAndDie() {
	fmt.Println(constants.VERSION)
	os.Exit(0)
}
