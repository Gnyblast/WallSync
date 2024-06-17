package main

import (
	"context"
	"log"
	"os/signal"
	"path"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/gnyblast/WallSync/internal/caches"
	"github.com/gnyblast/WallSync/internal/engines"
	"github.com/gnyblast/WallSync/internal/interfaces"
	"github.com/gnyblast/WallSync/internal/models"
	"github.com/gnyblast/WallSync/internal/services"
	"github.com/gnyblast/WallSync/internal/utils"
)

var serviceArgs models.ServiceArgs
var args interfaces.IArgs
var versionArg models.VersionArg

func main() {
	// f, err := os.Create("cpuprofile")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// pprof.StartCPUProfile(f)
	arg.Parse(&versionArg)
	if versionArg.Version {
		utils.PrintVersionAndDie()
	}

	arg.Parse(&serviceArgs)
	if serviceArgs.Version {
		utils.PrintVersionAndDie()
	}

	args = getArgsForService(serviceArgs)
	if versionArg.Version || args.ShouldPrintVersion() {
		utils.PrintVersionAndDie()
	}

	args.SetOutputDir(path.Join(args.GetOutputDir(), "WallSync"))
	var imageUpdateChannel chan string = make(chan string)
	defer close(imageUpdateChannel)

	utils.CheckDependencies(args.GetCommand())

	fsEngine := engines.NewFSEngine(args.GetOutputDir(), args)
	fsEngine.CreateDirectory()

	imageCache := caches.NewImageMetaCache(args)
	externalService := initService(serviceArgs, args, imageUpdateChannel, imageCache)
	backgroundImagecService := services.NewBackgroundImageService(args, imageUpdateChannel, imageCache)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	go backgroundImagecService.Listen()
	go ServiceRunner(externalService)

	<-ctx.Done()

	log.Print("Stoping the application...")

	externalService.Stop()
	<-externalService.IsStoped()

	backgroundImagecService.Quit <- true
	<-backgroundImagecService.Quit
	// pprof.StopCPUProfile()
}

func getArgsForService(serviceArgs models.ServiceArgs) interfaces.IArgs {
	switch serviceArgs.ServiceName {
	case "wallhaven":
		var args models.WallhavenArgs
		arg.MustParse(&args)
		return &args
	default:
		arg.MustParse(&serviceArgs)
		return &models.WallhavenArgs{}
	}
}

func initService(serviceArgs models.ServiceArgs, args interfaces.IArgs, imageUpdate chan string, imageCache *caches.ImageMetaCache) interfaces.IDownloadService {
	switch serviceArgs.ServiceName {
	case "wallhaven":
		return services.NewWallHavenService(args.(*models.WallhavenArgs), imageUpdate, imageCache)
	default:
		return services.NewWallHavenService(args.(*models.WallhavenArgs), imageUpdate, imageCache)
	}
}

func ServiceRunner(service interfaces.IDownloadService) {
	service.Start()
}
