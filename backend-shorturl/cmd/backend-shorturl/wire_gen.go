// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend-shorturl/internal/biz"
	"backend-shorturl/internal/conf"
	"backend-shorturl/internal/data"
	"backend-shorturl/internal/server"
	"backend-shorturl/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confService *conf.Service, confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	client, err := data.NewEntClient(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	redisClient, err := data.NewRedisClient(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(logger, client, redisClient)
	if err != nil {
		return nil, nil, err
	}
	shortUrlRepo := data.NewShortUrlRepo(dataData, logger)
	shortUrlUsecase := biz.NewShortUrlUsecase(shortUrlRepo, logger)
	publicService := service.NewPublicService(shortUrlUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, publicService, logger)
	httpServer := server.NewHTTPServer(confServer, publicService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
