package app

import (
	"fmt"
	"go_clean/internal/controller/grpc"
	"go_clean/internal/controller/http"
	"go_clean/internal/repo/persistent"
	"go_clean/internal/repo/temporary"
	"go_clean/internal/usecase/download"
	"go_clean/pkg/cache"
	"go_clean/pkg/grpcserver"
	"go_clean/pkg/httpserver"
	"go_clean/pkg/sqldb"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AlekSi/pointer"
)

func NewApp(configPath string) error {
	lastSlash := 0
	for i, v := range configPath {
		if v == '/' {
			lastSlash = i
		}
	}

	conf, err := ReadConfig(configPath[:lastSlash], configPath[lastSlash+1:])
	if err != nil {
		return fmt.Errorf("ReadConfig: %w", err)
	}

	go func() {
		httpserver.NewHttpServer(&httpserver.Server{
			Host:   conf.HTTP.Host,
			Port:   conf.HTTP.Port,
			Enable: conf.HTTP.Enable,
			Mux:    http.NewRoute(),
		})
	}()

	cc := cache.NewCache(conf.Cache.Host, conf.Cache.Port)
	downloadUsecases := download.NewDownload(
		pointer.To(persistent.NewSQLRepo(sqldb.NewDB(conf.PathToDB, conf.NameDB))),
		temporary.NewMemCache(cc),
	)

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases)
	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	cc.Close()
	err = grpcServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
