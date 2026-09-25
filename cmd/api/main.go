package main

import (
	"context"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
	"github.com/nguyenthienan91/bookmark-manager/internal/logger"
	"github.com/nguyenthienan91/bookmark-manager/internal/repository"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
	pkgredis "github.com/nguyenthienan91/bookmark-manager/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
)

// @title           Bookmark Manager API
// @version         1.0.0
// @description     API for Bookmark Manager application.
// @BasePath  		/

func main() {
	// 1. Load API config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	// 2. Initialize structured logger
	logger.Init(logger.Config{
		ServiceName: cfg.ServiceName,
		InstanceID:  cfg.InstanceID,
	})

	// 3. Load Redis config
	redisCfg, err := pkgredis.NewConfig()
	if err != nil {
		panic(err)
	}

	// 4. Create Redis client
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.L().Error().Err(err).Msg("error closing redis client")
		}
	}()

	// 5. Ping Redis
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("cannot connect to Redis: " + err.Error())
	}

	// 6. Create repository
	linkRepo := repository.NewLinkRepository(redisClient)
	healthRepo := repository.NewHealthCheckRepository(redisClient)

	// 7. Create shorten service
	shortenSvc := service.NewShortenLink(linkRepo)
	healthCheckSvc := service.NewHealthCheck(cfg.ServiceName, cfg.InstanceID, healthRepo)

	// 8. Start HTTP server
	app := api.NewEngine(cfg, healthCheckSvc, shortenSvc)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
