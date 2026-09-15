package main

import (
	"context"
	"log"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
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

	// 2. Load Redis config
	redisCfg, err := pkgredis.NewConfig()
	if err != nil {
		panic(err)
	}

	// 3. Create Redis client
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("error closing redis client: %v", err)
		}
	}()

	// 4. Ping Redis
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("cannot connect to Redis: " + err.Error())
	}

	// 5. Create repository
	linkRepo := repository.NewLinkRepository(redisClient)

	// 6. Create shorten service
	shortenSvc := service.NewShortenLink(linkRepo)

	// 7. Start HTTP server
	app := api.NewEngine(cfg, shortenSvc)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
