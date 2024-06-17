package interfaces

type IDownloadService interface {
	Start()
	Stop()
	IsStoped() <-chan bool
}
