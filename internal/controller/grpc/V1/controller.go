package v1

import (
	proto "go_clean/docs/proto/v1"
	"go_clean/internal/usecase"
)

type V1 struct {
	proto.DownloaderServer

	u usecase.Downloader
}
