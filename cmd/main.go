// Package main  это исполняемый файл приложения.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/victor8titov/rest-api-notes/internal/adaptor"
	"github.com/victor8titov/rest-api-notes/internal/service/http"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	diContainer, err := adaptor.NewDIContainer()
	if err != nil {
		log.Fatal("server", err)
	}
	defer diContainer.Close()

	httpService := http.NewService(diContainer)
	log.Println("started server")
	err = httpService.ListenAndServe(3000)
	if err != nil {
		log.Fatal("server", zap.Error(err))
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	log.Println("stopping service")

	const stopTimeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, stopTimeout)
	go func() {
		<-signals
		log.Println("force stopping service")
		cancel()
	}()

	log.Println("stopped")
}
