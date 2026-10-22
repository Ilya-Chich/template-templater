package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rendering/internal/adapters/config"
	httpadapter "rendering/internal/adapters/http"
	"rendering/internal/adapters/httpserver"
	"rendering/internal/adapters/templates"
	"rendering/internal/app"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var staticDirs = []string{"css", "js", "img"}

var staticFiles = []string{
	"robots.txt", "sitemap.xml", "site.webmanifest", "favicon.ico",
	"favicon-16x16.png", "favicon-32x32.png", "apple-touch-icon.png",
	"android-chrome-192x192.png", "android-chrome-512x512.png",
}

func loadEnv() {
	envFile := ".env.dev"
	if os.Getenv("DEBUG") == "false" {
		envFile = ".env.prod"
	}
	err := godotenv.Load(envFile)
	if err == nil {
		return
	}
	if errors.Is(err, os.ErrNotExist) && os.Getenv("HTTP_PORT") != "" {
		return
	}
	log.Fatalf("error loading %s file", envFile)
}

func main() {
	loadEnv()

	cfg, err := config.Parse(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := templates.Load("templates")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger(), httpadapter.SecurityHeaders())
	for _, dir := range staticDirs {
		router.Static("/"+dir, dir)
	}
	for _, file := range staticFiles {
		router.StaticFile("/"+file, file)
	}
	httpadapter.Register(router, httpadapter.NewHandler(app.NewPages(), renderer))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := httpserver.Run(ctx, srv, cfg); err != nil {
		log.Fatal(err)
	}
}
