package apiv1

import (
	proto "go_clean/docs/proto/v1"

	"go_clean/internal/usecase"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases usecase.Downloader) {
	r := &V1{
		u: usecases,
	}

	proto.RegisterDownloaderServer(app, r)
}
