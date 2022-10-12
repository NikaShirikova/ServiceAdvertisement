package main

import (
	"advertisement/internal/handler"
	"advertisement/internal/service"
	"advertisement/server"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	port := flag.String("p", "8081", "port input")
	portsParther := flag.String("d", "127.0.0.1:8085", "ports parthers input")
	flag.Parse()
	adrParthers := strings.Split(*portsParther, ",")
	fmt.Println(adrParthers)

	logger, err := zap.NewProduction()
	if err != nil {
		return
	}
	defer logger.Sync()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	validate := validator.New()
	services := service.NewService()
	handlers := handler.NewHandler(services, logger, validate, adrParthers)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(*port, handlers.InitRoutes(), logger); err != nil {
			logger.Error("Start server to fatal", zap.Error(err))
			return
		}
	}()

	<-done
	logger.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed ", zap.String("error ", err.Error()))
		return
	}
}
