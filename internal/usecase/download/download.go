package download

import (
	"go_clean/internal/repo/persistent"
	"go_clean/internal/repo/temporary"
	"go_clean/internal/usecase"
)

type downlaoder struct {
	dbRepo    *persistent.SQLRepo
	cacheRepo *temporary.Cache
}

func NewDownload(dbRepo *persistent.SQLRepo, cache *temporary.Cache) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
	}
}
