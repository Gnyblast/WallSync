package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/gnyblast/WallSync/internal/caches"
	"github.com/gnyblast/WallSync/internal/constants"
	"github.com/gnyblast/WallSync/internal/engines"
	"github.com/gnyblast/WallSync/internal/models"
	"github.com/gnyblast/WallSync/internal/services"
	"github.com/gnyblast/WallSync/internal/utils"
)

var args models.Args
var versionArg models.VersionArg

func main() {
	// f, err := os.Create("cpuprofile")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// pprof.StartCPUProfile(f)
	arg.Parse(&versionArg)
	if versionArg.Version {
		PrintVersionAndDie()
	}

	arg.MustParse(&args)
	if versionArg.Version || args.Version {
		PrintVersionAndDie()
	}

	args.Page = 1
	args.OutputDir = path.Join(args.OutputDir, "WallSync")
	var imageUpdateChannel chan string = make(chan string)
	defer close(imageUpdateChannel)

	utils.CheckDependencies(args.Command)

	fsEngine := engines.NewFSEngine(args.OutputDir, args)
	fsEngine.CreateDirectory()

	imageCache := caches.NewImageMetaCache(args)
	wallhavenService := services.NewWallHavenServiceService(args, imageUpdateChannel, imageCache)
	backgroundImagecService := services.NewBackgroundImageService(args, imageUpdateChannel, imageCache)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	go backgroundImagecService.Listen()
	go wallhavenService.Start()

	<-ctx.Done()

	log.Print("Stoping the application...")

	wallhavenService.Quit <- true
	<-wallhavenService.Quit
	backgroundImagecService.Quit <- true
	<-backgroundImagecService.Quit
	// pprof.StopCPUProfile()
}

func PrintVersionAndDie() {
	fmt.Println(constants.VERSION)
	os.Exit(0)
}
